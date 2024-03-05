package utils

import (
	"errors"
	"github.com/go-sql-driver/mysql"
)

func SQLDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return false
}

func SQLForeignKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1452 {
		return true
	}
	return false
}

func SQLNoRowsError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1146 {
		return true
	}
	return false
}

func SQLInvalidDataError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1366 {
		return true
	}
	return false
}
