package config

import (
	"fmt"
	"os"
	"strconv"
)

type config struct {
	PostgresConfig postgresConfig
}

type postgresConfig struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}

func (conf *postgresConfig) ToString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.host, conf.port, conf.user, conf.password, conf.dbname,
	)
}

const (
	pgHostEnvVar  = "POSTGRES_HOST"
	defaultPgHost = "localhost"

	pgPortEnvVar  = "POSTGRES_PORT"
	defaultPgPort = 5432

	pgUserEnvVar  = "POSTGRES_USER"
	defaultPgUser = "ais_user"

	pgPwdEnvVar       = "POSTGRES_PASSWORD"
	defaultPgPassword = "12345678"

	pgDbEnvVar      = "POSTGRES_DB"
	defaultPgDbname = "ais_sug"
)

func getPgConfig() postgresConfig {
	var PGConf postgresConfig
	var found bool

	if PGConf.host, found = os.LookupEnv(pgHostEnvVar); !found {
		PGConf.host = defaultPgHost
	}

	var PortString string
	if PortString, found = os.LookupEnv(pgPortEnvVar); !found {
		PGConf.port = defaultPgPort
	} else {
		if portInt, err := strconv.Atoi(PortString); err != nil {
			panic("Port should be a number if POSTGRES_PORT is specified")
		} else {
			PGConf.port = portInt
		}
	}

	if PGConf.user, found = os.LookupEnv(pgUserEnvVar); !found {
		PGConf.user = defaultPgUser
	}

	if PGConf.password, found = os.LookupEnv(pgPwdEnvVar); !found {
		PGConf.password = defaultPgPassword
	}

	if PGConf.dbname, found = os.LookupEnv(pgDbEnvVar); !found {
		PGConf.dbname = defaultPgDbname
	}

	return PGConf
}

func Config() config {
	return config{
		PostgresConfig: getPgConfig(),
	}
}
