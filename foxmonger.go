package foxmonger

import (
	"database/sql"
	"fmt"
	"regexp"

	"github.com/FoxFurry/foxmonger/internal/tag"
	"github.com/FoxFurry/foxmonger/internal/util"
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
	conf Config
}

func NewMonger(conf Config) FoxMonger {
	return &monger{
		conf: conf,
	}
}

func (m *monger) PopulateDatabase() error {
	var generators []tag.Generator

	for _, table := range m.conf.Tables {
		fmt.Printf("Working on table: %s\n", table.Name)

		for row, tagString := range table.Data {
			generator, err := m.tagsToGenerator(tagString)
			if err != nil {
				return fmt.Errorf("failed to create row %s generator: %w", row, err)
			}

			generators = append(generators, *generator)
		}
	}

	return nil
}

func (m *monger) tagsToGenerator(tagsString string) (*tag.Generator, error) {
	tagsValues, err := util.SplitTags(tagsString)
	if err != nil {
		return nil, fmt.Errorf("could not split tags: %w", err)
	}

	tagGenerator := tag.Generator{}

	for _, tagValue := range tagsValues {
		resolvedTag, err := m.resolveTag(tagValue)
		if err != nil {
			return nil, err
		}

		if resolvedTag == nil && len(tagsValues) != 1 {
			return nil, fmt.Errorf("auto tag cannot have any additional tags")
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
	case tagValue == "auto": // Auto is used to keep history of autoincrement rows
		return nil, nil

	case tagValue == "fullname":
		return tag.NewFullNameProducer(), nil

	case tagValue == "email":
		return tag.NewEmailProducer(), nil

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

	db, err := m.openConnection()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query)
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

func (m *monger) openConnection() (*sql.DB, error) {
	switch m.conf.DBType {
	case MySQLType:
		return sql.Open(MySQLType, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.conf.DBUser, m.conf.DBPass, m.conf.DBHost, m.conf.DBPort, m.conf.DBName))
	case PostgreSQL:
		return nil, fmt.Errorf("%s is not supported yet", m.conf.DBType)
	default:
		return nil, fmt.Errorf("unknown db type: %s", m.conf.DBType)
	}
}

//func generateQuery(tableName string, tableParams []rowParameter) (string, error) {
//	if tableParams == nil {
//		return "", fmt.Errorf("bruh")
//	}
//
//	output := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
//		tableName,
//		paramsToRowsString(tableParams),
//		paramsToValueString(tableParams))
//
//	return output, nil
//}
//
//func paramsToRowsString(tableParams []rowParameter) string {
//	if tableParams == nil {
//		return ""
//	}
//
//	rowsString := tableParams[0].RowName
//
//	for idx := 1; idx < len(tableParams); idx++ {
//		rowsString += fmt.Sprintf(", %s", tableParams[idx].RowName)
//	}
//
//	return rowsString
//}
//
//func paramsToValueString(tableParams []rowParameter) string {
//	if tableParams == nil {
//		return ""
//	}
//
//	rowsString := tableParams[0].RowGenerator.do()
//
//	for idx := 1; idx < len(tableParams); idx++ {
//		rowsString += fmt.Sprintf(", %s", tableParams[idx].RowGenerator.do())
//	}
//	return rowsString
//}
