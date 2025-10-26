package config

import (
	"flag"
)

type Config struct {
	SqliteDsn    string
	ListenAddres string
}

func NewEmptyConfig() *Config {
	return &Config{
		SqliteDsn:    "",
		ListenAddres: "",
	}
}

func (c *Config) PopulateFromArgs() {
	flag.StringVar(&c.SqliteDsn, "sqlite", ":memory:", "where the sqlite database should be stored")
	flag.StringVar(&c.ListenAddres, "listen", "127.0.0.1:5000", "where this program should listen for api requests")

	flag.Parse()
}
