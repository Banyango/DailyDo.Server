package util

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type StoreResult struct {
	Data      interface{}
	Total     int
	Err       error
}

type StoreChannel chan StoreResult
type SQLContextKey string

var TransactionContextKey SQLContextKey = "Transaction"
var TransactionWaitGroup SQLContextKey = "WaitGroup"

type SqlStore struct {
	Db *sqlx.DB
}

func (self *SqlStore) Execute(ctx context.Context, fn func(c context.Context) error) error {
	tx, err := self.Db.BeginTxx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	err = fn(context.WithValue(ctx, TransactionContextKey, tx))

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}