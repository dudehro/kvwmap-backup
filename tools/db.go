package tools

import (
	_ "database/sql"
	_ "github.com/lib/pg"
)

const (
	host     = "pgsql-server"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"

