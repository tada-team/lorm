package lorm

import (
	"database/sql"

	"github.com/jackc/pgx/v4"
)

var (
	conn    *sql.DB
	pgxConn *pgx.Conn
)

var (
	NonFatalErrors []string
	ShowSql        bool
	MaxAttempts    = 10
)

func SetConn(v *sql.DB) {
	conn = v
}

func SetPgxConn(v *pgx.Conn) {
	pgxConn = v
}
