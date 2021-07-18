package lorm

import (
	"database/sql"
	"log"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type Transactional interface {
	Tx() *Tx
	SetTx(tx *Tx)
}

type BaseTransactional struct{ tx *Tx }

func (t *BaseTransactional) Tx() *Tx { return t.tx }

func (t *BaseTransactional) SetTx(tx *Tx) { t.tx = tx }

func NewTx(tx *sql.Tx, num int64) *Tx {
	return &Tx{Tx: tx, num: num}
}

type Tx struct {
	*sql.Tx
	num       int64
	objects   []Transactional
	onSuccess []func() error
}

func (tx *Tx) String() string {
	return "[tx:" + strconv.FormatInt(tx.num, 10) + "]"
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
		log.Panicln("already in transaction!")
	}
	t.SetTx(tx)
	tx.objects = append(tx.objects, t)
}

var txNum int64

func Atomic(fn func(tx *Tx) error) error {
	start := time.Now()

	currentTxNum := atomic.AddInt64(&txNum, 1)
	sqlTx, txErr := conn.Begin()
	if txErr != nil {
		return errors.Wrapf(txErr, "[tx:"+strconv.FormatInt(currentTxNum, 10)+"] begin failed")
	}

	tx := NewTx(sqlTx, currentTxNum)
	if ShowSql {
		log.Println(tx.String() + "begin:" + breadcrumb())
	}

	err := fn(tx)

	for _, t := range tx.objects {
		t.SetTx(nil)
	}

	if err != nil {
		if ShowSql {
			log.Println(tx.String(), "rollback:", err.Error())
		}
		if err := tx.Rollback(); err != nil {
			return errors.Wrapf(err, tx.String()+" rollback failed")
		}
		return err
	}

	startCommit := time.Now()
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, tx.String()+" commit failed")
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
