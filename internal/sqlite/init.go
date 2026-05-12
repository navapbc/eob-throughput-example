package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
	// _ "github.com/tursodatabase/turso-go"
)

// With mattn
// const SQLITE_DRIVER = "sqlite3"
// With modernc
const SQLITE_DRIVER = "sqlite"

// With Turso
// const SQLITE_DRIVER = "turso"

//go:embed schema.sql
var ddl string

type SQLite struct {
	DB   *sql.DB
	N    int
	Name string
	Tx   *sql.Tx
	Q    *Queries
	QTx  *Queries
}

func NewSQLite(n int) *SQLite {
	name := fmt.Sprintf("./tput-%d.sqlite", n)
	db, err := sql.Open(SQLITE_DRIVER, name)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1)

	if _, err := db.ExecContext(context.Background(), ddl); err != nil {
		log.Println(err)
	}

	return &SQLite{
		DB:   db,
		N:    n,
		Name: name,
		Q:    New(db),
	}
}

func (s *SQLite) BeginTx() {
	tx, err := s.DB.BeginTx(context.Background(), nil)
	if err != nil {
		log.Println(err)
	}
	s.Tx = tx
	s.QTx = s.Q.WithTx(tx)
}

func (s *SQLite) EndTx() {
	err := s.Tx.Commit()
	if err != nil {
		log.Println(err)
	}
}
