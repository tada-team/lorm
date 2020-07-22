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
	MaxAttempts    = 10
	ShowSql        = false
	NotFatalErrors = []string{
		"bad connection",
		"broken pipe",
		"connection refused",
		"connection reset",
		"missing destination name",
		"open /opt/tada/cfg/yandex.crt",
		"read-only transaction",
		"the database system is in recovery mode",
	}
)

func SetConn(v *sql.DB) {
	conn = v
}

func SetPgxConn(v *pgx.Conn) {
	pgxConn = v
}
