package config

type Config struct {
	Profiles []*Profile `json:"profiles"`
}

func (c *Config) GetAllBlobConfigs() []*Blob {
	return nil
}
