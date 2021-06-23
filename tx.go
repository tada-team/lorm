package lorm

import (
	"database/sql"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type Transactional interface {
	Tx() *Tx
	SetTx(tx *Tx)
}

type BaseTransactional struct {
	tx       *Tx
	postSave func() error
}

func (t BaseTransactional) Tx() *Tx { return t.tx }

func (t *BaseTransactional) SetTx(tx *Tx, postSave func() error) {
	t.tx = tx
	t.postSave = postSave
}

func NewTx(tx *sql.Tx, num int64) *Tx {
	return &Tx{Tx: tx, num: num}
}

type Tx struct {
	*sql.Tx
	num       int64
	objects   []Transactional
	onSuccess []func() error
}

func (tx Tx) String() string {
	return fmt.Sprintf("[tx:%d]", tx.num)
}

func (tx *Tx) OnSuccess(fn func() error) error {
	tx.onSuccess = append(tx.onSuccess, fn)
	return nil
}

func (tx *Tx) Add(t Transactional) {
	if current := t.Tx(); current != nil {
		if current == t.Tx() {
			return
		}
		log.Panicln("switchTx() already in transaction!")
	}
	t.SetTx(tx)
	tx.objects = append(tx.objects, t)
}

var txNum int64

func Atomic(fn func(tx *Tx) error) error {
	start := time.Now()
	atomic.AddInt64(&txNum, 1)
	sqlTx, txErr := conn.Begin()
	if txErr != nil {
		return errors.Wrapf(txErr, "[tx:%d] begin failed", txNum)
	}

	tx := NewTx(sqlTx, txNum)
	if ShowSql {
		log.Printf("%s begin: %s", tx, breadcrumb())
	}

	err := fn(tx)

	for _, t := range tx.objects {
		t.SetTx(nil)
	}

	if err != nil {
		if ShowSql {
			log.Printf("%s rollback: %s", tx, err)
		}
		if err := tx.Rollback(); err != nil {
			return errors.Wrapf(err, "%s rollback failed", tx)
		}
		return err
	}

	startCommit := time.Now()
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "%s commit failed", tx)
	}

	if ShowSql {
		log.Printf(
			"%s commit (%s) total transaction: %s",
			tx,
			time.Since(startCommit).Truncate(time.Millisecond),
			time.Since(start).Truncate(time.Millisecond),
		)
	}

	for _, fn := range tx.onSuccess {
		if err := fn(); err != nil {
			return err
		}
	}
	tx.onSuccess = nil

	return nil
}

func MaybeAddToTx(tx *Tx, t Transactional) {
	if tx != nil {
		tx.Add(t)
	}
}
