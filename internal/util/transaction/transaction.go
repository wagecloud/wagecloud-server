package transaction

import (
	"github.com/wagecloud/wagecloud-server/internal/logger"
	"go.uber.org/zap"
)

type Op func() error
type UndoOp func() error

type Transaction struct {
	ops       []Op
	undoStack []UndoOp
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (tx *Transaction) Add(op Op, undo UndoOp) error {
	tx.ops = append(tx.ops, op)
	tx.undoStack = append(tx.undoStack, undo)
	return nil
}

func (tx *Transaction) Commit() error {
	for i := len(tx.ops) - 1; i >= 0; i-- {
		err := tx.ops[i]()
		if err != nil {
			logger.Log.Error("failed to commit op ", zap.Int("index", i), zap.Error(err))
			return err
		}
	}
	tx.ops = nil
	return nil
}

func (tx *Transaction) Rollback() error {
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
