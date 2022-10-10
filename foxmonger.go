package foxmonger

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/FoxFurry/foxmonger/internal/tag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jaswdr/faker"
)

var (
	enumPat    = regexp.MustCompile("enum:.+")
	limitPat   = regexp.MustCompile("limit:\\d+")
	foreignPat = regexp.MustCompile("foreign:.+")

	TagSplitter = ";"
)

type foreignKey struct {
	TargetTable,
	TargetRow string
}

type FoxMonger struct {
	db            *sql.DB
	fakerInstance faker.Faker
	conf          Config
}

func NewMonger(conf Config) *FoxMonger {
	database, err := openConnection(conf)
	if err != nil {
		log.Fatalf("failed initialize: %v", err)
	}

	return &FoxMonger{
		fakerInstance: faker.New(),
		conf:          conf,
		db:            database,
	}
}

func (m *FoxMonger) PopulateDatabase() error {
	var (
		generators      []tag.Generator
		generatorBuffer *tag.Generator
		queryTemplate   string
		err             error
	)

	for i := range m.conf.Tables {
		table := &m.conf.Tables[i]

		fmt.Printf("Working on table: %s\n", table.Name)

		for row, tagString := range table.Data {
			generatorBuffer, err = m.tagsToGenerator(tagString, row)
			if err != nil {
				return fmt.Errorf("failed to create row \"%s\" generator: %w", row, err)
			}

			generators = append(generators, *generatorBuffer)
		}

		fmt.Println("Generators created")

		queryTemplate = generateQueryTemplate(table.Name, generators)

		tx, err := m.db.Begin()
		if err != nil {
			return err
		}

		var (
			queryBuffer    string
			queryExporting string
		)

		for idx := 0; idx < table.BaseMultiplier*m.conf.BaseCount; idx++ {
			queryBuffer = fmt.Sprintf("%s (%s)", queryTemplate, paramsToValueString(generators))

			if _, err := tx.Exec(queryBuffer); err != nil {
				return err
			}

			if table.ExportQueries {
				queryExporting += queryBuffer + ";\n"
			}
		}

		fmt.Println("Transaction generated, applying")

		if err := tx.Commit(); err != nil {
			return err
		}

		if table.ExportQueries {
			if err := ioutil.WriteFile(table.ExportPath, []byte(queryExporting), 0644); err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func (m *FoxMonger) tagsToGenerator(tagsString, rowName string) (*tag.Generator, error) {
	var (
		tagsValues   = strings.Split(tagsString, TagSplitter)
		tagGenerator = tag.Generator{
			RowName: rowName,
		}
	)

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

func (m *FoxMonger) resolveTag(tagValue string) (any, error) {
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
		foreignTarget, err := resolveForeignKey(tagValue)
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

func (m *FoxMonger) getRows(tableName, rowName string) ([]string, error) {
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

func generateQueryTemplate(dbName string, tableParams []tag.Generator) string {
	if tableParams == nil {
		return ""
	}

	var rows []string

	for idx := 0; idx < len(tableParams); idx++ {
		rows = append(rows, tableParams[idx].RowName)
	}

	template := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES",
		dbName,
		strings.Join(rows, ","),
	)

	return template
}

func paramsToValueString(tableParams []tag.Generator) string {
	var rowValues []string

	for idx := range tableParams {
		rowValues = append(rowValues, "'"+tableParams[idx].Do()+"'")
	}

	return strings.Join(rowValues, ",")
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

func resolveForeignKey(foreignTag string) (*foreignKey, error) {
	foreignElements := strings.Split(foreignTag, TagSplitter)

	if len(foreignElements) != 2 {
		return nil, fmt.Errorf("tag does not contain foreign key description: %s", foreignTag)
	}

	var tableName, rowName string

	_, err := fmt.Sscanf(foreignElements[1], "%s(%s)", &tableName, &rowName)
	if err != nil {
		return nil, fmt.Errorf("failed parsing foreign key: %w", err)
	}

	return &foreignKey{
		TargetTable: tableName,
		TargetRow:   rowName,
	}, nil
}
