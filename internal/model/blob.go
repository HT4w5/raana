package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/HT4w5/raana/internal/config"
	"github.com/goccy/go-yaml"
)

type BlobType int

const (
	TypeHTTP BlobType = iota
	TypeLocal
)

var blobTypeName = map[BlobType]string{
	TypeHTTP:  "http",
	TypeLocal: "local",
}

var blobNameType = map[string]BlobType{
	"http":  TypeHTTP,
	"local": TypeLocal,
}

type Blob struct {
	cfg      *config.Blob
	blobType BlobType
	data     map[string]interface{}
}

func NewBlob(cfg *config.Blob) (*Blob, error) {
	b := &Blob{
		cfg: cfg,
	}

	// Determine type
	if t, ok := blobNameType[cfg.Type]; !ok {
		return nil, errors.New("invalid blob type")
	} else {
		b.blobType = t
	}

	if err := b.load(); err != nil {
		return nil, err
	}

	return b, nil
}

// Load blob according to its type
func (b *Blob) load() error {
	// Parse url
	u, err := url.Parse(b.cfg.URL)
	if err != nil {
		return err
	}

	// Fetch data according to type
	var dataBytes []byte
	errScheme := errors.New("incorrect scheme for blob type local")
	switch b.blobType {
	case TypeLocal:
		if u.Scheme != "file" {
			return errScheme
		}
		dataBytes, err = os.ReadFile(u.Path)
		if err != nil {
			return err
		}
	case TypeHTTP:
		if u.Scheme != "http" && u.Scheme != "https" {
			return errScheme
		}
		resp, err := http.Get(b.cfg.URL)
		if err != nil {
			return fmt.Errorf("error making HTTP request: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("bad status code: %s", resp.Status)
		}
		dataBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported blob type")
	}

	// Determine blob format from extension and unmarshal
	ext := filepath.Ext(u.Path)
	switch ext {
	case ".json":
		if err := json.Unmarshal(dataBytes, &b.data); err != nil {
			return err
		}
	case ".yaml":
		fallthrough
	case ".yml":
		if err := yaml.Unmarshal(dataBytes, &b.data); err != nil {
			return err
		}
	default:
		return errors.New("blob format not supported")
	}

	return nil
}
