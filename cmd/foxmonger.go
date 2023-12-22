package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/FoxFurry/foxmonger"
	"github.com/FoxFurry/foxmonger/internal/database"
)

const (
	mysqlType      = "mysql"
	postgresqlType = "postgresql"
)

var (
	flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	help = flags.Bool("help", false, "TBD")

	dbType = flags.String("type", "", "TBD")
	dbName = flags.String("db", "", "TBD")
	dbUser = flags.String("user", "", "TBD")
	dbPass = flags.String("pass", "", "TBD")
	dbHost = flags.String("host", "", "TBD")
	dbPort = flags.String("port", "", "TBD")
)

func main() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		osError("failed reading flags: %v\n", err)
	}

	if *help {
		usage()
		os.Exit(0)
	}

	if err := checkFlags(dbType, dbName, dbUser, dbPass, dbHost, dbPort); err != nil {
		osError("error: %v\n", err)
	}

	targetDB, err := openConnection(*dbType, *dbUser, *dbPass, *dbHost, *dbPort, *dbName)
	if err != nil {
		osError("failed to open db connection: %v\n", err)
	}

	db := database.NewDatabase(targetDB)

	monger := foxmonger.NewMonger(db)
	if err := monger.PopulateDatabase(context.Background()); err != nil {
		osError("failed to populate db: %v\n", err)
	}

	return
}

func checkFlags(dbType, dbName, dbUser, dbPass, dbHost, dbPort *string) error {
	switch "" {
	case *dbType:
		return isMissingError("type")
	case *dbName:
		return isMissingError("name")
	case *dbUser:
		return isMissingError("user")
	case *dbPass:
		return isMissingError("pass")
	case *dbHost:
		return isMissingError("host")
	case *dbPort:
		return isMissingError("port")
	}

	return nil
}

func isMissingError(missing string) error {
	return fmt.Errorf("flag '%s' is not set", missing)
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

func openConnection(dbType, dbUser, dbPass, dbHost, dbPort, dbName string) (*sql.DB, error) {
	return sql.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName))
}
