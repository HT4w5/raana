package model

import (
	"encoding/json"
	"errors"
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

}

func (b *Blob) loadFromLocal() error {
	// Parse url
	u, err := url.Parse(b.cfg.URL)
	if err != nil {
		return err
	}

	if u.Scheme != "file" {
		return errors.New("Incorrect scheme for blob type local")
	}

	dataBytes, err := os.ReadFile(u.Path)
	if err != nil {
		return err
	}

	// Determine blob format from extension and unmarshal
	ext := filepath.Ext(u.Path)
	switch ext {
	case ".json":
		json.Unmarshal(dataBytes, b.data)
	case ".yaml":
		fallthrough
	case ".yml":
		yaml.Unmarshal(dataBytes, b.data)
	default:
		return errors.New("blob format not supported")
	}

	return nil
}
