package sql

import (
	"context"
	"database/sql"

	// import mssql driver
	_ "github.com/denisenkom/go-mssqldb"
)

// Connector stores connection information for the SQL Server instance
// and contains instance methods that can run SQL queries and commands
type Connector struct {
	ConnectionString string
}

// Execute an SQL statement and ignore the results
func (c Connector) Execute(command string, args ...interface{}) error {
	ctx := context.Background()
	conn, err := sql.Open("sqlserver", c.ConnectionString)
	defer conn.Close()

	if err != nil {
		return err
	}

	_, err = conn.ExecContext(ctx, command, args...)
	if err != nil {
		return err
	}
	return nil
}

// Query the database
func (c Connector) Query(query string, scanner func(*sql.Rows) error, args ...interface{}) error {
	ctx := context.Background()
	conn, err := sql.Open("sqlserver", c.ConnectionString)
	if err != nil {
		return err
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = scanner(rows)
	if err != nil {
		return err
	}

	return nil
}
