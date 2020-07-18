package lorm

import (
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Debug struct {
	ShowSql bool
}

var NonFatalErrors = []string{
	"bad connection",
	"broken pipe",
	"connection refused",
	"connection reset",
	"missing destination name",
	"open /opt/tada/cfg/yandex.crt",
	"read-only transaction",
	"the database system is in recovery mode",
}

var pool *pgxpool.Pool

func Pool() *pgxpool.Pool {
	if pool == nil {
		log.Panicln("lorm: connection not initialized")
	}
	return pool
}

func SetPool(p *pgxpool.Pool) {
	pool = p
}
