package foxmonger

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/go-sql-driver/mysql"

	"github.com/FoxFurry/foxmonger/internal/tag"
	"github.com/FoxFurry/foxmonger/internal/util"
	"github.com/jaswdr/faker"
)

var (
	enumPat    = regexp.MustCompile("enum:.+")
	limitPat   = regexp.MustCompile("limit:\\d+")
	foreignPat = regexp.MustCompile("foreign:.+")
)

type FoxMonger interface {
	PopulateDatabase() error
}

type monger struct {
	db            *sql.DB
	fakerInstance faker.Faker
	conf          Config
}

func NewMonger(conf Config) FoxMonger {
	database, err := openConnection(conf)
	if err != nil {
		log.Fatalf("failed initialize: %v", err)
	}

	return &monger{
		fakerInstance: faker.New(),
		conf:          conf,
		db:            database,
	}
}

func (m *monger) PopulateDatabase() error {
	var generators []tag.Generator

	for _, table := range m.conf.Tables {
		fmt.Printf("Working on table: %s\n", table.Name)

		for row, tagString := range table.Data {
			generator, err := m.tagsToGenerator(tagString, row)
			if err != nil {
				return fmt.Errorf("failed to create row \"%s\" generator: %w", row, err)
			}

			generators = append(generators, *generator)
		}

		queryTemplate := fmt.Sprintf("INSERT INTO %s (%s) VALUES", table.Name, paramsToRowsString(generators))

		var transaction string
		for idx := 0; idx < table.BaseMultiplier*m.conf.BaseCount; idx++ {
			transaction += fmt.Sprintf("%s ( %s );\n", queryTemplate, paramsToValueString(generators))
		}

		fmt.Println("Transaction generated, applying")
		if err := m.executeQuery(transaction); err != nil {
			return err
		}
	}

	return nil
}

func (m *monger) tagsToGenerator(tagsString, rowName string) (*tag.Generator, error) {
	tagsValues, err := util.SplitTags(tagsString)
	if err != nil {
		return nil, fmt.Errorf("could not split tags: %w", err)
	}

	tagGenerator := tag.Generator{
		RowName: rowName,
	}

	for _, tagValue := range tagsValues {
		resolvedTag, err := m.resolveTag(tagValue)
		if err != nil {
			return nil, err
		}

		if producerTag, ok := resolvedTag.(tag.Producer); ok {
			if err := tagGenerator.SetProducer(producerTag); err != nil {
				return nil, fmt.Errorf("could not set producer %s: %w", tagValue, err)
			}

		} else if modifierTag, ok := resolvedTag.(tag.Modifier); ok {
			tagGenerator.AddModifier(modifierTag)
		}
	}

	return &tagGenerator, nil
}

func (m *monger) resolveTag(tagValue string) (any, error) {
	switch {
	// Exact tags
	case tagValue == "fullname":
		return tag.NewFullNameProducer(&m.fakerInstance), nil

	case tagValue == "email":
		return tag.NewEmailProducer(&m.fakerInstance), nil

	// Pattern tags
	case enumPat.MatchString(tagValue):
		return tag.NewEnumProducer(tagValue)

	case limitPat.MatchString(tagValue):
		return tag.NewLimitModifier(tagValue)

	case foreignPat.MatchString(tagValue):
		foreignTarget, err := util.ResolveForeignKey(tagValue)
		if err != nil {
			return nil, fmt.Errorf("could not create foreign key producer: %w", err)
		}

		foreignRows, err := m.getRows(foreignTarget.TargetTable, foreignTarget.TargetRow)
		if err != nil {
			return nil, fmt.Errorf("could not get foreign target rows: %w", err)
		}

		return tag.NewForeignProducer(foreignRows), nil
	default:
		return nil, fmt.Errorf("unknown tag: %s", tagValue)
	}
}

func (m *monger) getRows(tableName, rowName string) ([]string, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", rowName, tableName)

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		rowsResult []string
		rowBuffer  any
	)

	for rows.Next() {
		if err := rows.Scan(&rowBuffer); err != nil {
			panic("this should not have happened: " + err.Error())
		}

		rowsResult = append(rowsResult, fmt.Sprintf("%v", rowBuffer))
	}

	return rowsResult, nil
}

func (m *monger) executeQuery(query string) error {
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func paramsToRowsString(tableParams []tag.Generator) string {
	if tableParams == nil {
		return ""
	}

	rowsString := tableParams[0].RowName

	for idx := 1; idx < len(tableParams); idx++ {
		rowsString += fmt.Sprintf(", %s", tableParams[idx].RowName)
	}

	return rowsString
}

func paramsToValueString(tableParams []tag.Generator) string {
	if tableParams == nil {
		return ""
	}

	rowsString := fmt.Sprintf("'%s'", tableParams[0].Do())

	for idx := 1; idx < len(tableParams); idx++ {
		rowsString += fmt.Sprintf(", '%s'", tableParams[idx].Do())
	}
	return rowsString
}

func openConnection(conf Config) (*sql.DB, error) {
	switch conf.DBType {
	case MySQLType:
		return sql.Open(MySQLType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName))
	case PostgreSQL:
		return nil, fmt.Errorf("%s is not supported yet", conf.DBType)
	default:
		return nil, fmt.Errorf("unknown db type: %s", conf.DBType)
	}
}
