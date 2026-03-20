package main

import (
    "context"
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "log"
)

func main() {
    db, err := sqlx.Connect("postgres", "postgres://appuser2:app123@localhost:5432/School-manajement?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    
    query := `
        SELECT pg_get_constraintdef(c.oid)
        FROM pg_constraint c
        JOIN pg_class t ON c.conrelid = t.oid
        WHERE t.relname = 'absensis' AND c.contype = 'c';
    `
    var defs []string
    err = db.SelectContext(context.Background(), &defs, query)
    if err != nil {
        log.Fatal(err)
    }
    for _, d := range defs {
        fmt.Println(d)
    }
}
