package main

import (
	"errors"
	"net/http"

	"bitcoinreport.yutaseko.net/internal/data"
)

func (app *application) showAllTransactionHandler(w http.ResponseWriter, r *http.Request) {
	transactions, err := app.models.Transactions.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"transactions": transactions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	transaction, err := app.models.Transactions.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transaction": transaction}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		AmountBTC       float64 `json:"amount_btc"`
		PricePerBTC     int64   `json:"price_per_btc"`
		TransactionType int8    `json:"transaction_type"`
		Note            string  `json:"note"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	transaction := &data.Transaction{
		AmountBTC:       input.AmountBTC,
		PricePerBTC:     input.PricePerBTC,
		TransactionType: input.TransactionType,
		Note:            input.Note,
	}

	// add validator

	err = app.models.Transactions.Insert(transaction)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// if successful, add the address of where to see the new transaction to the header

	err = app.writeJSON(w, http.StatusCreated, envelope{"transaction": transaction}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	transaction, err := app.models.Transactions.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		AmountBTC       *float64 `json:"amount_btc"`
		PricePerBTC     *int64   `json:"price_per_btc"`
		TransactionType *int8    `json:"transaction_type"`
		Note            *string  `json:"note"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.AmountBTC != nil {
		transaction.AmountBTC = *input.AmountBTC
	}
	if input.PricePerBTC != nil {
		transaction.PricePerBTC = *input.PricePerBTC
	}
	if input.TransactionType != nil {
		transaction.TransactionType = *input.TransactionType
	}
	if input.Note != nil {
		transaction.Note = *input.Note
	}
	err = app.models.Transactions.Update(transaction)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"transaction": transaction}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Transactions.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "task successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
