package db

import (
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

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
