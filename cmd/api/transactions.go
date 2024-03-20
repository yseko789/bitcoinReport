package main

import (
	"net/http"

	"bitcoinreport.yutaseko.net/internal/data"
)

func (app *application) showAllTransactionHandler(w http.ResponseWriter, r *http.Request) {
	transactions := []data.Transaction{
		{
			ID:              1,
			AmountBTC:       1.0,
			PricePerBTC:     60000,
			Note:            "First purchase",
			TransactionType: 1,
			Version:         1,
		},
		{
			ID:              2,
			AmountBTC:       1.4,
			PricePerBTC:     56000,
			Note:            "Second Purchase",
			TransactionType: 1,
			Version:         1,
		},
		{
			ID:              3,
			AmountBTC:       0.8,
			PricePerBTC:     61000,
			Note:            "First Sell",
			TransactionType: 2,
			Version:         1,
		},
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"transactions": transactions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"find id": id}, nil)
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
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"update id": id}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"delete id": id}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
