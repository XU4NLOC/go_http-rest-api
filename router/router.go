package router

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go_http-rest-api/db"
	"go_http-rest-api/handlers"
)

func Setup() http.Handler {
	mux := http.NewServeMux()

	// Health and readiness endpoints used by Docker, reverse proxies and Kubernetes later.
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/healthz", health)
	mux.HandleFunc("/ready", ready)
	mux.HandleFunc("/readyz", ready)

	// Account routes
	mux.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAccounts(w, r)
		case http.MethodPost:
			handlers.CreateAccount(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/accounts/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetAccount(w, r)
		case http.MethodDelete:
			handlers.DeleteAccount(w, r)
		}
	})

	// Transaction routes
	mux.HandleFunc("/accounts/{id}/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTransactions(w, r)
		case http.MethodPost:
			handlers.CreateTransaction(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Summary route
	mux.HandleFunc("/accounts/{id}/summary", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetSummary(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

func health(w http.ResponseWriter, r *http.Request) {
	writeStatus(w, http.StatusOK, map[string]string{"status": "ok"})
}

func ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		writeStatus(w, http.StatusServiceUnavailable, map[string]string{
			"status": "not_ready",
		})
		return
	}

	writeStatus(w, http.StatusOK, map[string]string{"status": "ready"})
}

func writeStatus(w http.ResponseWriter, status int, payload map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
