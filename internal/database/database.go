package database

import (
	"context"
	"database/sql"
)

type Database interface {
	ReadTables(context.Context) ([]string, error)
	ReadRowsFromTable(context.Context, string) ([]string, error)
}

type database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) Database {
	return &database{db: db}
}

func (d *database) ReadTables(ctx context.Context) ([]string, error) {
	tablesRows, err := d.db.QueryContext(ctx, "SHOW TABLES;")
	if err != nil {
		return nil, err
	}

	var tables []string
	var tableBuffer string

	for tablesRows.Next() {
		if err := tablesRows.Scan(&tableBuffer); err != nil {
			return nil, err
		}

		tables = append(tables, tableBuffer)
	}

	return tables, nil
}

func (d *database) ReadRowsFromTable(ctx context.Context, table string) ([]string, error) {
	columnsRows, err := d.db.QueryContext(ctx, "SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = ?", table)
	if err != nil {
		return nil, err
	}

	var columns []string
	var columnBuffer string

	for columnsRows.Next() {
		if err := columnsRows.Scan(&columnBuffer); err != nil {
			return nil, err
		}

		columns = append(columns, columnBuffer)
	}

	return columns, nil
}
