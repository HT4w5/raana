package model

import (
	"fmt"

	"github.com/HT4w5/raana/internal/config"
)

type Override struct {
	cfg       *config.Override
	prepends  []*Blob
	appends   []*Blob
	overrides []*Blob
}

func NewOverride(cfg *config.Override, bp *BlobPool) (*Override, error) {
	o := &Override{
		cfg:       cfg,
		prepends:  make([]*Blob, len(cfg.Prepends)),
		appends:   make([]*Blob, len(cfg.Appends)),
		overrides: make([]*Blob, len(cfg.Overrides)),
	}

	groupTagMap := map[*[]string]*[]*Blob{
		&cfg.Prepends:  &o.prepends,
		&cfg.Appends:   &o.appends,
		&cfg.Overrides: &o.overrides,
	}

	for strs, bs := range groupTagMap {
		for _, v := range *strs {
			b := bp.GetBlob(v)
			if b == nil {
				return nil, fmt.Errorf("override %s has missing blob %s", cfg.Tag, v)
			}
			*bs = append((*bs), b)
		}
	}

	return o, nil
}
