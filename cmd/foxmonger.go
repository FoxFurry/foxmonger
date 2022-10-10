package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/FoxFurry/foxmonger"
	"github.com/spf13/viper"
)

const (
	mysqlType      = "mysql"
	postgresqlType = "postgresql"
)

var (
	flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	config = flags.String("config", "", "File containing SQL configuration for monger. Find more about it on official github page")
	help   = flags.Bool("help", false, "Test")
)

func main() {
	flags.Parse(os.Args[1:])

	if *help {
		usage()
		os.Exit(0)
	}

	if *config == "" {
		osError("config is mandatory for execution\n")
	}

	viper.SetConfigFile(*config)
	if err := viper.ReadInConfig(); err != nil {
		osError("failed to read config: %v\n", err)
	}

	conf := new(foxmonger.Config)

	if err := viper.Unmarshal(conf); err != nil {
		osError("failed to unmarshal config: %v\n", err)
	}

	db, err := openConnection(conf)
	if err != nil {
		osError("failed to open db connection: %v\n", err)
	}

	monger := foxmonger.NewMonger(conf, db)

	if err := monger.PopulateDatabase(); err != nil {
		osError("failed to populate db: %v\n", err)
	}

	return
}

func usage() {
	fmt.Printf(`
			Usage: %s [flags]
			TBD: General description will go here	
			Flags available:
		`,
		os.Args[0],
	)

	flags.PrintDefaults()
}

func osError(format string, opts ...any) {
	fmt.Fprintf(os.Stderr, format, opts...)
	os.Exit(1)
}

func openConnection(conf *foxmonger.Config) (*sql.DB, error) {
	switch conf.DBType {
	case mysqlType:
		return sql.Open(mysqlType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName))
	case postgresqlType:
		return nil, fmt.Errorf("%s is not supported yet", conf.DBType)
	default:
		return nil, fmt.Errorf("unknown db type: %s", conf.DBType)
	}
}
