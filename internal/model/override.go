package model

import (
	"fmt"

	"github.com/HT4w5/raana/internal/config"
)

type Override struct {
	cfg      *config.Override
	prepends []*Blob
	appends  []*Blob
}

func NewOverride(cfg *config.Override, bp *BlobPool) (*Override, error) {
	o := &Override{
		cfg:      cfg,
		prepends: make([]*Blob, len(cfg.Prepends)),
		appends:  make([]*Blob, len(cfg.Appends)),
	}

	for _, v := range cfg.Prepends {
		b := bp.GetBlob(v)
		if b == nil {
			return nil, fmt.Errorf("override %s has missing blob %s", cfg.Tag, v)
		}
		o.prepends = append(o.prepends, b)
	}

	for _, v := range cfg.Appends {
		b := bp.GetBlob(v)
		if b == nil {
			return nil, fmt.Errorf("override %s has missing blob %s", cfg.Tag, v)
		}
		o.prepends = append(o.prepends, b)
	}

	return o, nil
}
