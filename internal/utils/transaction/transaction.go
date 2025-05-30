package transaction

import (
	"github.com/wagecloud/wagecloud-server/internal/logger"
	"go.uber.org/zap"
)

type Op func() error
type UndoOp func() error

type Transaction struct {
	undoStack []UndoOp
	done      bool
}

func NewTransaction(done bool) *Transaction {
	return &Transaction{done: done}
}

// func (tx *Transaction) Add(op Op, undo UndoOp) error {
// 	tx.ops = append(tx.ops, op)
// 	tx.undoStack = append(tx.undoStack, undo)
// 	return nil
// }

func (tx *Transaction) Do(op Op, undo UndoOp) error {
	err := op()
	if err != nil {
		return err
	}
	tx.undoStack = append(tx.undoStack, undo)
	return nil
}

func (tx *Transaction) Commit() {
	// for i := len(tx.ops) - 1; i >= 0; i-- {
	// 	err := tx.ops[i]()
	// 	if err != nil {
	// 		logger.Log.Error("failed to commit op ", zap.Int("index", i), zap.Error(err))
	// 		return err
	// 	}
	// }
	// tx.ops = nil
	tx.done = true
}

// Rollback rolls back the transaction, safe to call even if there was no error
func (tx *Transaction) Rollback() error {
	// Only rollback if we haven't already done the transaction
	if tx.done {
		return nil
	}

	for i := len(tx.undoStack) - 1; i >= 0; i-- {
		err := tx.undoStack[i]()
		if err != nil {
			logger.Log.Error("failed to rollback op ", zap.Int("index", i), zap.Error(err))
			return err
		}
	}
	tx.undoStack = nil
	return nil
}
