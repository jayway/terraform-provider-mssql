package sql

import (
	"fmt"
)

// CreateLogin connects to the SQL Database to create a login with the provided
// credentials
func (c Connector) CreateLogin(username string, password string) error {
	return c.Execute(
		// SQL Server does not support creating logins with parameters, so we build the string directly
		fmt.Sprintf("CREATE LOGIN %s WITH PASSWORD = '%s'", username, password),
	)
}
