package config

type Override struct {
	Tag       string   `json:"tag"`
	Prepends  []string `json:"prepends"`
	Appends   []string `json:"appends"`
	Overrides []string `json:"overrides"`
}
