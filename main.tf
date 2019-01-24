provider "mssql" {
  connection_string = "Data Source=localhost:1433;user id=sa;password=$up3R$3cR37"
}

provider "random" {}

resource "random_string" "pwd" {
  keepers = {
    username = "terraform"
  }

  upper   = true
  lower   = true
  number  = true
  special = true
  length  = 23
}

resource "mssql_login" "test" {
  username = "${random_string.pwd.keepers.username}"
  password = "${random_string.pwd.result}"
}
