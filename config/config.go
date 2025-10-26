package config

import (
	"flag"
)

type Config struct {
	SqliteDsn    string
	ListenAddres string
}

const (
	defaultSqliteDsn     = ":memory"
	defaultListenAddress = "127.0.0.1:5000"
)

func NewDefaultConfig() *Config {
	return &Config{
		SqliteDsn:    defaultSqliteDsn,
		ListenAddres: defaultListenAddress,
	}
}

func (c *Config) PopulateFromArgs() {
	var sqliteDsn = flag.String("sqlite", "", "where the sqlite database should be stored")
	var listenAddress = flag.String("listen", "", "where this program should listen for api requests")

	flag.Parse()

	if *sqliteDsn != "" {
		c.SqliteDsn = *sqliteDsn
	}

	if *listenAddress != "" {
		c.ListenAddres = *listenAddress
	}
}
