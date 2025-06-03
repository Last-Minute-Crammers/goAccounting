package config

type Config struct {
	System System `yaml:"system"`
}

type System struct {
	Addr   int    `yaml:"addr"`
	DbType string `yaml:"db-type"`
	Mode   string `yaml:"mode"`
}
