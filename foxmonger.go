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
	db   *sql.DB
	conf *Config

	fakerInstance faker.Faker
}

func NewMonger(conf *Config, database *sql.DB) *FoxMonger {
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

		var (
			queryBuilder strings.Builder
			bufferSize   = 1000
			numOfQueries = table.BaseMultiplier * m.conf.BaseCount
		)

		queryBuilder.WriteString(queryTemplate + " " + paramsToValueString(generators))
		for idx := 2; idx <= numOfQueries; idx++ {
			queryBuilder.WriteString(",\n" + paramsToValueString(generators))

			if idx%bufferSize == 0 || idx == numOfQueries {
				if _, err := m.db.Exec(queryBuilder.String()); err != nil {
					log.Printf("Writing batch")
					panic(err)
				}

				queryBuilder.Reset()
				queryBuilder.WriteString(queryTemplate + " " + paramsToValueString(generators))
			}
		}

		fmt.Println("Transaction generated, applying")

		//if !table.Dummy {
		//	if _, err := m.db.Exec(queryBuilder.String()); err != nil {
		//		panic(err)
		//	}
		//}

		if table.ExportQueries {
			if err := ioutil.WriteFile(table.ExportPath, []byte(queryBuilder.String()), 0644); err != nil {
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

	return "(" + strings.Join(rowValues, ",") + ")"
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
