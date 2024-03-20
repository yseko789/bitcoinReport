package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/transactions", app.showAllTransactionHandler)
	router.HandlerFunc(http.MethodPost, "/v1/transactions", app.createTransactionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/transactions/:id", app.showTransactionHandler)
	router.HandlerFunc(http.MethodPut, "/v1/transactions/:id", app.updateTransactionHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/transactions/:id", app.deleteTransactionHandler)

	return router
}
