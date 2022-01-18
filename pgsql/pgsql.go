package pgsql

import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
    "errors"
)

var globalConnection *sql.DB

func OpenConnection(host string, port string, user string, password string, dbname string) {

    if globalConnection == nil {

        conn_string := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
        fmt.Println(conn_string)
        db, err := sql.Open("postgres", conn_string)
        if err != nil {
            panic(err)
        }

        err = db.Ping()
        if err != nil {
            panic(err)
        }

        if db != nil {
            fmt.Println("Verbindung zur Datenbank erfolgreich hergestellt.")
        }

        globalConnection = db

//        defer db.Close()

    }
}

func Select(query string) (*sql.Rows){

    if globalConnection == nil {
        panic(errors.New("Select nicht möglich, es besteht keine Verbindung zur DB!"))
    }

    rows, err := globalConnection.Query(query)
    if err != nil {
        panic(err)
    }

//    defer rows.Close()

    return rows
}

func ListSchemas() ([]string) {

    list := make([]string, 0)

    rows := Select(`select s.nspname as table_schema
                     from pg_catalog.pg_namespace s
                join pg_catalog.pg_user u on u.usesysid = s.nspowner
                order by table_schema`)

    for rows.Next() {
        var table_schema string
        err := rows.Scan(&table_schema)
        if err != nil {
            panic(err)
        }
        list = append(list, table_schema)
    }
    defer rows.Close()
    return list
}
