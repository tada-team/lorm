package lorm

import "strings"

func IsCancelingStatementDueToLockTimeout(err error) bool {
	return err != nil && strings.Contains(err.Error(), "canceling statement due to lock timeout")
}

func IsDuplicateKeyViolatesUniqueConstraint(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func IsViolatesForeignKeyConstraint(err error) bool {
	return err != nil && strings.Contains(err.Error(), "violates foreign key constraint")
}

func IsDeadlockDetected(err error) bool {
	return err != nil && strings.Contains(err.Error(), "deadlock detected")
}
