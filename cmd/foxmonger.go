package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FoxFurry/foxmonger"
	"github.com/spf13/viper"
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
		fmt.Fprintf(os.Stderr, "config is mandatory for execution\n")
		os.Exit(1)
	}

	viper.SetConfigFile(*config)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read config: %v\n", err)
		os.Exit(1)
	}

	conf := foxmonger.Config{}

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal config: %v\n", err)
		os.Exit(1)
	}

	monger := foxmonger.NewMonger(conf)

	if err := monger.PopulateDatabase(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to populate db: %v\n", err)
		os.Exit(1)
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
