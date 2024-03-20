package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Transaction struct {
	ID              int64     `json:"id"`
	CreatedAt       time.Time `json:"-"`
	AmountBTC       float64   `json:"amount_btc"`
	PricePerBTC     int64     `json:"price_per_btc"`
	TransactionType int8      `json:"transaction_type"`
	Note            string    `json:"note,omitempty"`
	Version         int32     `json:"version"`
}

type TransactionModel struct {
	DB *sql.DB
}

func (t TransactionModel) Insert(transaction *Transaction) error {
	query := `
	INSERT INTO transactions (amount_btc, price_per_btc, transaction_type, note)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	args := []any{transaction.AmountBTC, transaction.PricePerBTC, transaction.TransactionType, transaction.Note}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&transaction.ID, &transaction.CreatedAt, &transaction.Version)
}

func (t TransactionModel) Get(id int64) (*Transaction, error) {
	query := `
	SELECT id, created_at, amount_btc, price_per_btc, transaction_type, note, version
	FROM transactions
	WHERE id = $1`

	var transaction Transaction

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := t.DB.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.CreatedAt,
		&transaction.AmountBTC,
		&transaction.PricePerBTC,
		&transaction.TransactionType,
		&transaction.Note,
		&transaction.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &transaction, nil
}

func (t TransactionModel) GetAll() ([]*Transaction, error) {
	query := `
	SELECT id, created_at, amount_btc, price_per_btc, transaction_type, note, version
	FROM transactions
	ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	transactions := []*Transaction{}

	for rows.Next() {
		var transaction Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.CreatedAt,
			&transaction.AmountBTC,
			&transaction.PricePerBTC,
			&transaction.TransactionType,
			&transaction.Note,
			&transaction.Version,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}

func (t TransactionModel) Update(transaction *Transaction) error {
	query := `
	UPDATE transactions
	SET amount_btc = $1, price_per_btc = $2, transaction_type = $3, note = $4, version = version + 1
	WHERE id = $5 and version = $6
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		transaction.AmountBTC,
		transaction.PricePerBTC,
		transaction.TransactionType,
		transaction.Note,
		transaction.ID,
		transaction.Version,
	}

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&transaction.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}
	return nil
}

func (t TransactionModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM transactions
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
