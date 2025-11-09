package config

type Profile struct {
	Tag       string      `json:"tag"`
	Primary   *Blob       `json:"primary"`
	Overrides []*Override `json:"overrides"`
}
