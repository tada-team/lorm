package lorm

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgconn"

	"github.com/tada-team/lorm/op"

	"github.com/jackc/pgx/v4"
)

func TxLock2(tx *Tx, k1, k2 int) error {
	args := op.NewArgs()
	query := op.Select(op.PgAdvisoryXactLock2(k1, k2))
	_, err := TxExec(tx, query, args)
	return err
}

func TxExec(tx *Tx, q op.Query, args op.Args) (res pgconn.CommandTag, err error) {
	query := q.Query()
	defer trackQuery(tx, query, args)()
	err = retry(func() error {
		res, err = doExec(tx, query, args)
		return err
	})
	return res, err
}

func TxQuery(tx *Tx, q op.Query, args op.Args, each func(pgx.Rows) error) (err error) {
	query := q.Query()
	defer trackQuery(nil, query, args)()
	err = retry(func() error {
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
	return
}

func TxScan(tx *Tx, q op.Query, args op.Args, dest ...interface{}) error {
	query := q.Query()
	defer trackQuery(nil, query, args)()
	return retry(func() error {
		if tx == nil {
			conn, err := Pool().Acquire(context.Background())
			if err != nil {
				return err
			}
			defer conn.Release()
			return conn.QueryRow(context.Background(), query, args...).Scan(dest...)
		}
		return tx.tx.QueryRow(context.Background(), query, args...).Scan(dest...)
	})
}

func doExec(tx *Tx, query string, args op.Args) (pgconn.CommandTag, error) {
	if tx == nil {
		conn, err := Pool().Acquire(context.Background())
		if err != nil {
			return nil, err
		}
		defer conn.Release()
		return conn.Exec(context.Background(), query, args...)
	}
	return tx.tx.Exec(context.Background(), query, args...)
}

func doQuery(tx *Tx, query string, args op.Args) (pgx.Rows, error) {
	if tx == nil {
		conn, err := Pool().Acquire(context.Background())
		if err != nil {
			return nil, err
		}
		defer conn.Release()
		return conn.Query(context.Background(), query, args...)
	}
	return tx.tx.Query(context.Background(), query, args...)
}

func retry(fn func() error) error {
	const maxAttempts = 10
	i := 0
	for {
		err := fn()
		if err != nil && nonFatalError(err) && i <= maxAttempts {
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
