package lorm

import (
	"database/sql"
	"strings"
)

var (
	conn *sql.DB
)

var (
	NonFatalErrors []string
	ShowSql        bool
	MaxAttempts    = 10
)

func SetConn(v *sql.DB) {
	conn = v
}

var disableLocks = false

func DisableLocks() { disableLocks = true }

func EnableLocks() { disableLocks = false }

func SetDbParam(tx *Tx, arg, value string) (err error) {
	value = strings.ReplaceAll(value, "'", "")
	// Prepared statement doesn't work with SET. FIXME: add more sql safety
	query := "SET " + arg + " = '" + value + "'"
	if tx == nil {
		_, err = conn.Exec(query)
	} else {
		_, err = tx.Exec(query)
	}
	return
}

