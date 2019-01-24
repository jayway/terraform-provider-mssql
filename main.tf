provider "mssql" {
  connection_string = "Data Source=localhost:1433;user id=sa;password=$up3R$3cR37"
}

resource "mssql_login" "test" {
  username = "terraform"
  password = "P@ssw0rd"
}
