package lorm

import (
	"github.com/tada-team/lorm/op"
)

type BaseUpdater struct {
	ch op.Changes
}

func NewUpdater() BaseUpdater { return BaseUpdater{ch: make(op.Changes)} }

func (u *BaseUpdater) Ch(f op.Column, v interface{}) { u.ch[f] = v }

func (u *BaseUpdater) HasChanges() bool { return len(u.ch) > 0 }

func (u *BaseUpdater) DoSave(r Record, t op.Table) error {
	if !r.HasPk() {
		u.ch = make(op.Changes)
		return r.Save()
	}

	if len(u.ch) == 0 {
		return nil
	}

	args := op.NewArgs()
	kv := make(op.Set, len(u.ch))
	for k, v := range u.ch {
		kv[k] = args.Next(v)
	}

	query := op.Update(t, kv).Where(r.PkCond(&args))
	if _, err := TxExec(r.Tx(), r, query, args); err != nil {
		return err
	}

	if err := r.PostSave(); err != nil {
		return err
	}

	u.ch = make(op.Changes)
	return nil
}
