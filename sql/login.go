package sql

import (
	"database/sql"
)

// CreateLogin connects to the SQL Database to create a login with the provided
// credentials
func (c Connector) CreateLogin(username string, password string) error {
	cmd := `DECLARE @sql nvarchar(max)
					SET @sql = 'CREATE LOGIN ' + QuoteName(@username) + ' ' +
										 'WITH PASSWORD = ' + QuoteName(@password, '''')
					EXEC (@sql)`
	return c.Execute(cmd, sql.Named("username", username), sql.Named("password", password))
}

// DeleteLogin connects to the SQL Database and removes a login with the provided
// username, if it exists. If it does not exist, this is a noop.
func (c Connector) DeleteLogin(username string) error {
	cmd := `DECLARE @sql nvarchar(max);
					SET @sql = 'IF EXISTS (SELECT 1 FROM [master].[sys].[server_principals] WHERE [name] = ' + QuoteName(@username, '''') + ') ' +
										 'DROP LOGIN ' + QuoteName(@username);
					EXEC (@sql)`
	err := c.killSessionsForLogin(username)
	if err != nil {
		return err
	}
	return c.Execute(cmd, sql.Named("username", username))
}

// Login represents an SQL Server Login
type Login struct {
	Username    string
	PrincipalID int64
}

// GetLogin reads a login from the SQL Database, if it exists. If it does not,
// no error is returned, but the returned Login is nil
func (c Connector) GetLogin(username string) (*Login, error) {
	var principalID int64 = -1

	err := c.Query(
		"SELECT principal_id FROM [master].[sys].[server_principals] WHERE [name] = @username",
		func(r *sql.Rows) error {
			for r.Next() {
				err := r.Scan(&principalID)
				if err != nil {
					return err
				}
			}
			return nil
		},
		sql.Named("username", username),
	)

	if err != nil {
		return nil, err
	}
	if principalID != -1 {
		return &Login{Username: username, PrincipalID: principalID}, nil
	}
	return nil, nil
}

// UpdateLogin updates the password of a login, if it exists.
func (c Connector) UpdateLogin(username string, password string) error {
	cmd := `DECLARE @sql nvarchar(max)
					SET @sql = 'IF EXISTS (SELECT 1 FROM [master].[sys].[server_principals] WHERE [name] = ' + QuoteName(@username, '''') + ') ' +
										 'ALTER LOGIN ' + QuoteName(@username) + ' ' +
										 'WITH PASSWORD = ' + QuoteName(@password, '''')
					EXEC (@sql)`

	return c.Execute(cmd, sql.Named("username", username), sql.Named("password", password))
}

func (c Connector) killSessionsForLogin(username string) error {
	cmd := ` -- adapted from https://stackoverflow.com/a/5178097/38055
	DECLARE sessionsToKill CURSOR FAST_FORWARD FOR
			SELECT session_id
			FROM sys.dm_exec_sessions
			WHERE login_name = @username
	OPEN sessionsToKill

	DECLARE @sessionId INT
	DECLARE @statement NVARCHAR(200)

	FETCH NEXT FROM sessionsToKill INTO @sessionId

	WHILE @@FETCH_STATUS = 0
	BEGIN
			PRINT 'Killing session ' + CAST(@sessionId AS NVARCHAR(20)) + ' for login ' + @username

			SET @statement = 'KILL ' + CAST(@sessionId AS NVARCHAR(20))
			EXEC sp_executesql @statement

			FETCH NEXT FROM sessionsToKill INTO @sessionId
	END

	CLOSE sessionsToKill
	DEALLOCATE sessionsToKill`

	return c.Execute(cmd, sql.Named("username", username))
}
