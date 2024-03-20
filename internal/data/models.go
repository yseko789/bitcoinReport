package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Transactions TransactionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Transactions: TransactionModel{DB: db},
	}
}
