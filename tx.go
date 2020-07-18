package lorm

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/pkg/errors"
)

type Transactional interface {
	Tx() *Tx
	SetTx(tx *Tx)
}

type BaseTransactional struct{ tx *Tx }

func (t BaseTransactional) Tx() *Tx       { return t.tx }
func (t *BaseTransactional) SetTx(tx *Tx) { t.tx = tx }

type Tx struct {
	tx      pgx.Tx
	num     int64
	objects []Transactional
}

func (tx Tx) String() string {
	return fmt.Sprintf("[tx:%d]", tx.num)
}

func NewTx(tx pgx.Tx, num int64) *Tx {
	return &Tx{
		tx:      tx,
		num:     num,
		objects: make([]Transactional, 0),
	}
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

	sqlTx, txErr := Pool().Begin(context.Background())
	if txErr != nil {
		return errors.Wrapf(txErr, "[tx:%d] begin failed", txNum)
	}

	tx := NewTx(sqlTx, txNum)
	if Debug.ShowSql {
		log.Printf("%s begin: %s", tx, breadcrumb())
	}

	err := fn(tx)

	for _, t := range tx.objects {
		t.SetTx(nil)
	}

	if err != nil {
		if Debug.ShowSql {
			log.Printf("%s rollback: %s", tx, err)
		}
		if err := tx.tx.Rollback(context.Background()); err != nil {
			return errors.Wrapf(err, "%s rollback failed", tx)
		}
		return err
	}

	startCommit := time.Now()
	if err := tx.tx.Commit(context.Background()); err != nil {
		return errors.Wrapf(err, "%s commit failed", tx)
	}

	if Debug.ShowSql {
		log.Printf(
			"%s commit (%s) total transaction: %s",
			tx,
			time.Since(startCommit).Truncate(time.Millisecond),
			time.Since(start).Truncate(time.Millisecond),
		)
	}

	return nil
}

func MaybeAddToTx(tx *Tx, t Transactional) {
	if tx != nil {
		tx.Add(t)
	}
}
