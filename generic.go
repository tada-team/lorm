package lorm

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"
	"github.com/tada-team/lorm/op"
)

func DoCount(f Filter, table op.Table) int {
	var res int
	if !f.IsEmpty() {
		query := op.Select(op.Count(table.Pk())).From(table).Where(f.GetConds())
		err := TxScan(f.Tx(), query, f.GetArgs(), &res)
		if err != nil {
			log.Panicln(errors.Wrapf(err, "%s.Count() fail on %s and %v", table, query, f.GetArgs()))
		}
	}
	return res
}

func DoDeleteFiltered(f Filter, table op.Table) error {
	if !f.IsEmpty() {
		query := op.Delete(table).Where(f.GetConds())
		if _, err := TxExec(f.Tx(), query, f.GetArgs()); err != nil {
			return err
		}
	}
	return nil
}

func DoExists(f Filter, table op.Table) bool {
	var res bool
	if !f.IsEmpty() {
		query := op.RawQuery("SELECT", op.Exists(op.Select(op.One).From(table).Where(f.GetConds())))
		err := TxScan(f.Tx(), query, f.GetArgs(), &res)
		if err != nil {
			log.Panicln(errors.Wrapf(err, "%s.Exists() fail on %s and %v", table, query, f.GetArgs()))
		}
	}
	return res
}

func DoSaveall(err error, r Record, t op.Table) error {
	if err != nil {
		return err
	}
	if r.HasPk() {
		return DoUpdate(r, t)
	}
	return DoInsert(r, t)
}

func DoUpdate(r Record, t op.Table) error {
	kv := make(op.Set)
	values := r.GetAllFields()
	args := op.NewArgs()
	for i, f := range t.GetAllFields() {
		if f.BareName() != t.Pk().BareName() {
			kv[f] = args.Next(values[i])
		}
	}
	query := op.Update(t, kv).Where(r.PkCond(&args))
	_, err := TxExec(r.Tx(), query, args)
	return err
}

func DoInsert(r Record, t op.Table) error {
	kv := make(op.Set)
	values := r.GetAllFields()
	args := op.NewArgs()
	for i, f := range t.GetAllFields() {
		if f.BareName() != t.Pk().BareName() || !isEmpty(values[0]) {
			kv[f] = args.Next(values[i])
		}
	}
	query := op.Insert(t, kv).Returning(t.Pk().BareName())
	if err := TxScan(r.Tx(), query, args, values[0]); err != nil {
		return err
	}
	if !r.HasPk() {
		return errors.New("programming error: no pk after insert")
	}
	return nil
}

func DoDelete(r Record, t op.Table) error {
	args := op.NewArgs()
	query := op.Delete(t).Where(r.PkCond(&args))
	_, err := TxExec(r.Tx(), query, args)
	return err
}

func DoGet(f Filter, r Record, t op.Table) (bool, error) {
	if f.IsEmpty() {
		return false, nil
	}
	query := op.Select().From(t).Where(f.GetConds()).Lock(f.GetLock()).OrderBy(f.GetOrderBy())
	err := TxScan(f.Tx(), query, f.GetArgs(), r.GetAllFields()...)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	MaybeAddToTx(f.Tx(), r)
	return true, nil
}

var ReloadError = errors.New("lorm: reload error")

func DoReload(r Record, t op.Table) error {
	args := op.NewArgs()
	query := op.Select().From(t).Where(r.PkCond(&args))
	err := TxScan(r.Tx(), query, args, r.GetAllFields()...)
	if err == sql.ErrNoRows {
		return ReloadError
	}
	return err
}
