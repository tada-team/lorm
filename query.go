package lorm

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/tada-team/lorm/op"

	"github.com/jackc/pgx/v4"
)

func TxLock2(tx *Tx, k1, k2 int) error {
	args := op.NewArgs()
	query := op.Select(op.PgAdvisoryXactLock2(k1, k2))
	_, err := TxExec(tx, query, args)
	return err
}

func TxExec(tx *Tx, q op.Query, args op.Args) (res sql.Result, err error) {
	query := q.Query()
	defer trackQuery(tx, query, args)()
	err = retry(func() error {
		res, err = doExec(tx, query, args)
		return err
	})
	return res, err
}

func TxQuery(tx *Tx, q op.Query, args op.Args, each func(*sql.Rows) error) error {
	query := q.Query()
	defer trackQuery(tx, query, args)()
	return retry(func() error {
		rows, err := doQuery(tx, query, args)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			err := each(rows)
			if err != nil {
				return err
			}
		}
		return rows.Err()
	})
}

func TxQueryPgx(tx pgx.Tx, q op.Query, args op.Args, each func(pgx.Rows) error) (err error) {
	query := q.Query()
	defer trackQuery(nil, query, args)()
	err = retry(func() error {
		rows, err := doQueryPgx(tx, query, args)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			err := each(rows)
			if err != nil {
				return err
			}
		}
		return rows.Err()
	})
	return
}

func TxScan(tx *Tx, q op.Query, args op.Args, dest ...interface{}) error {
	query := q.Query()
	defer trackQuery(tx, query, args)()
	return retry(func() error { return doQueryRow(tx, query, args).Scan(dest...) })
}

func TxScanPgx(q op.Query, args op.Args, dest ...interface{}) error {
	query := q.Query()
	defer trackQuery(nil, query, args)()
	return retry(func() error { return doQueryRowPgx(query, args).Scan(dest...) })
}

func doExec(tx *Tx, query string, args op.Args) (sql.Result, error) {
	if tx == nil {
		return conn.Exec(query, args...)
	}
	return tx.Exec(query, args...)
}

func doQuery(tx *Tx, query string, args op.Args) (*sql.Rows, error) {
	if tx == nil {
		return conn.Query(query, args...)
	}
	return tx.Query(query, args...)
}

func doQueryPgx(tx pgx.Tx, query string, args op.Args) (pgx.Rows, error) {
	if tx == nil {
		return pgxConn.Query(context.Background(), query, args...)
	}
	return tx.Query(context.Background(), query, args...)
}

func doQueryRow(tx *Tx, query string, args op.Args) *sql.Row {
	if tx == nil {
		return conn.QueryRow(query, args...)
	}
	return tx.QueryRow(query, args...)
}

func doQueryRowPgx(query string, args op.Args) pgx.Row {
	return pgxConn.QueryRow(context.Background(), query, args...)
}

func retry(fn func() error) error {
	i := 0
	for {
		err := fn()
		if err != nil && nonFatalError(err) && i <= MaxAttempts {
			i++
			time.Sleep(time.Duration(i) * time.Second)
			continue
		}
		return err
	}
}

func nonFatalError(err error) bool {
	for _, s := range NonFatalErrors {
		if strings.Contains(err.Error(), s) {
			return true
		}
	}
	return false
}
