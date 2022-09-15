package parser

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //nolint
	"github.com/pkg/errors"
)

// GetCreateTableFromDB get create table info from mysql
func GetCreateTableFromDB(dsn, tableName string) (string, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", errors.WithMessage(err, "open db error")
	}
	defer db.Close() //nolint

	rows, err := db.Query("SHOW CREATE TABLE " + tableName)
	if err != nil {
		return "", errors.WithMessage(err, "query show create table error")
	}

	defer rows.Close() //nolint
	if !rows.Next() {
		return "", errors.Errorf("table(%s) not found", tableName)
	}

	var table string
	var createSQL string
	err = rows.Scan(&table, &createSQL)
	if err != nil {
		return "", err
	}

	return createSQL, nil
}
