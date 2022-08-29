package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

const (
	tagSplitter = ";"
)

type table struct {
	Name           string            `mapstructure:"name"`
	BaseMultiplier int               `mapstructure:"base_multiplier"`
	Data           map[string]string `mapstructure:"data"`
}

type queryParams struct {
	RowName      string
	RowGenerator *Generator
}

type mongerConfig struct {
	BaseCount int     `mapstructure:"base_count"`
	DBType    string  `mapstructure:"db_type"`
	DBName    string  `mapstructure:"db_name"`
	DBHost    string  `mapstructure:"db_host"`
	DBUser    string  `mapstructure:"db_user"`
	DBPass    string  `mapstructure:"db_pass"`
	DBPort    string  `mapstructure:"db_port"`
	Tables    []table `mapstructure:"tables"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	conf := mongerConfig{}
	tagsManager := tagManager{}

	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}

	//db, err := openConnection(&conf)
	//if err != nil {
	//	panic(err)
	//}

	for _, table := range conf.Tables {

		var tableQueryParams []queryParams

		//fmt.Printf("Working on table %s with multiplier %d\n\n", table.Name, table.BaseMultiplier)

		for rowName, rowTag := range table.Data {
			tagGenerator, err := tagsManager.resolveTag(rowTag)
			if err != nil {
				//fmt.Printf("Tag \"%s\" could not be generated: %v\n", rowName, err)
			} else if tagGenerator == nil {
				//fmt.Printf("Tag \"%s\" is ommitted\n", rowName)
			} else {
				//fmt.Printf("Tag \"%s\" generated: %s\n", rowName, value)
				tableQueryParams = append(tableQueryParams, queryParams{rowName, tagGenerator})
			}
		}

		for i := 0; i < conf.BaseCount*table.BaseMultiplier; i++ {
			query, err := generateQuery(table.Name, tableQueryParams)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(query)
			}
		}
	}
}

func openConnection(conf *mongerConfig) (*sql.DB, error) {
	switch conf.DBType {
	case "mysql":
		return sql.Open(conf.DBType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName))
	default: // TODO: Add more DBs support
		return nil, fmt.Errorf("%s is not supported yet", conf.DBType)
	}
}

func generateQuery(tableName string, tableParams []queryParams) (string, error) {
	if tableParams == nil {
		return "", fmt.Errorf("bruh")
	}

	output := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		paramsToRowsString(tableParams),
		paramsToValueString(tableParams))

	return output, nil
}

func paramsToRowsString(tableParams []queryParams) string {
	if tableParams == nil {
		return ""
	}

	rowsString := tableParams[0].RowName

	for idx := 1; idx < len(tableParams); idx++ {
		rowsString += fmt.Sprintf(", %s", tableParams[idx].RowName)
	}

	return rowsString
}

func paramsToValueString(tableParams []queryParams) string {
	if tableParams == nil {
		return ""
	}

	rowsString := tableParams[0].RowGenerator.do()

	for idx := 1; idx < len(tableParams); idx++ {
		rowsString += fmt.Sprintf(", %s", tableParams[idx].RowGenerator.do())
	}

	return rowsString
}
