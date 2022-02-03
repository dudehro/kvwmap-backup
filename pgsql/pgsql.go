package pgsql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var globalConnection *sql.DB

func OpenConnection(host string, port string, user string, password string, dbname string) {

	if globalConnection == nil {

		conn_string := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
		fmt.Println(conn_string)
		db, err := sql.Open("postgres", conn_string)
		if err != nil {
			log.Println(err)
		}

		err = db.Ping()
		if err != nil {
			log.Println(err)
		}

		if db != nil {
			log.Println("Verbindung zur Datenbank erfolgreich hergestellt.")
		}

		globalConnection = db

		//        defer db.Close()

	}
}

func Select(query string) *sql.Rows {

	if globalConnection == nil {
		log.Println(errors.New("Select nicht möglich, es besteht keine Verbindung zur DB!"))
	}

	rows, err := globalConnection.Query(query)
	if err != nil {
		log.Println(err)
	}

	//    defer rows.Close()

	return rows
}

func ListSchemas() []string {

	list := make([]string, 0)

	rows := Select(`select s.nspname as table_schema
                     from pg_catalog.pg_namespace s
                join pg_catalog.pg_user u on u.usesysid = s.nspowner
                order by table_schema`)

	if rows == nil {
		return []string{}
	}

	for rows.Next() {
		var table_schema string
		err := rows.Scan(&table_schema)
		if err != nil {
			log.Println(err)
		}
		list = append(list, table_schema)
	}
	defer rows.Close()
	return list
}
