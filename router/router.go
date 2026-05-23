package router

import (
	"net/http"

	"go_http-rest-api/handlers"
)

func Setup() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAccounts(w, r)
		case http.MethodPost:
			handlers.CreateAccounts(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAccounts(w, r)
		case http.MethodDelete:
			handlers.DeleteAccount(w, r)
		}
	})

	mux.HandleFunc("/accounts/{id}/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTransactions(w, r)
		case http.MethodPost:
			handlers.CreateTransactions(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/accounts/{id}/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetSummary(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
