package config

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	SqliteDsn    string
	ListenAddres string
}

const (
	defaultSqliteDsn     = ":memory:"
	defaultListenAddress = "127.0.0.1:5000"
)

func NewDefaultConfig() *Config {
	return &Config{
		SqliteDsn:    defaultSqliteDsn,
		ListenAddres: defaultListenAddress,
	}
}

func (c *Config) PopulateFromArgs() {
	var sqliteDsn = flag.String("sqlite", "", "where the sqlite database should be stored (default 127.0.0.1:5000)")
	var listenAddress = flag.String("listen", "", "where this program should listen for api requests (default :memory:)")

	flag.Parse()

	if *sqliteDsn != "" {
		c.SqliteDsn = *sqliteDsn
	}

	if *listenAddress != "" {
		c.ListenAddres = *listenAddress
	}
}

func (c *Config) PopulateFromEnv() {
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatalf("unable to read environment variables: %v", err)
	}

	if sqliteDsn, ok := myEnv["SQLITE_DSN"]; ok {
		c.SqliteDsn = sqliteDsn
	}

	if listenAddress, ok := myEnv["LISTEN_ADDRESS"]; ok {
		c.ListenAddres = listenAddress
	}
}
