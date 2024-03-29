package lorm

import (
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/tada-team/lorm/op"
)

func TxLock2(tx *Tx, k1, k2 int) error {
	args := op.NewArgs()
	query := op.Select(op.PgAdvisoryXactLock2(k1, k2))
	_, err := TxExec(tx, nil, query, args)
	return err
}

func TxExec(tx *Tx, locker sync.Locker, q op.Query, args op.Args) (res sql.Result, err error) {
	query := q.Query()
	defer trackQuery(tx, query, args)()
	err = retry(func() error {
		res, err = doExec(tx, locker, query, args)
		return err
	})
	return res, err
}

func TxQuery(tx *Tx, q op.Query, args op.Args, each func(*sql.Rows) error) error {
	query := q.Query()
	defer trackQuery(tx, query, args)()
	return retry(func() error {
		rows, err := doQuery(tx, nil, query, args)
		if err != nil {
			return err
		}
		defer func() { _ = rows.Close() }()
		for rows.Next() {
			err := each(rows)
			if err != nil {
				return err
			}
		}
		return rows.Err()
	})
}

func TxScan(tx *Tx, locker sync.Locker, q op.Query, args op.Args, dest ...interface{}) error {
	query := q.Query()
	defer trackQuery(tx, query, args)()
	return retry(func() error { return doQueryRow(tx, locker, query, args).Scan(dest...) })
}

func doExec(tx *Tx, locker sync.Locker, query string, args op.Args) (sql.Result, error) {
	if locker != nil && !disableLocks {
		locker.Lock()
		defer locker.Unlock()
	}
	if tx == nil {
		return conn.Exec(query, args...)
	}
	return tx.Exec(query, args...)
}

func doQuery(tx *Tx, locker sync.Locker, query string, args op.Args) (*sql.Rows, error) {
	if locker != nil && !disableLocks {
		locker.Lock()
		defer locker.Unlock()
	}
	if tx == nil {
		return conn.Query(query, args...)
	}
	return tx.Query(query, args...)
}

func doQueryRow(tx *Tx, locker sync.Locker, query string, args op.Args) *sql.Row {
	if locker != nil && !disableLocks {
		locker.Lock()
		defer locker.Unlock()
	}
	if tx == nil {
		return conn.QueryRow(query, args...)
	}
	return tx.QueryRow(query, args...)
}

func retry(fn func() error) error {
	i := 0
	for {
		err := fn()
		if err != nil && nonFatalError(err) && i <= MaxAttempts {
			i++
			log.Println("lorm: warn:", err, "retry:", i)
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
