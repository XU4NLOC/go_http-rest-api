package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go_http-rest-api/db"
	"go_http-rest-api/models"
)

func GetAccounts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, name, type, created_at FROM accounts ORDER BY id")
	if err != nil {
		http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	accounts := []models.Account{}
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.CreatedAt); err != nil {
			http.Error(w, "Failed to scan account", http.StatusInternalServerError)
			return
		}
		accounts = append(accounts, a)
	}

	writeJSON(w, http.StatusOK, accounts)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var a models.Account
	row := db.DB.QueryRow("SELECT id, name, type, created_at FROM accounts WHERE id = $1", id)
	if err := row.Scan(&a.ID, &a.Name, &a.Type, &a.CreatedAt); err != nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, a)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var a models.Account
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if a.Name == "" || a.Type == "" {
		http.Error(w, "Name and type are required", http.StatusBadRequest)
		return
	}

	err := db.DB.QueryRow(
		"INSERT INTO accounts (name, type) VALUES ($1, $2) RETURNING id, created_at",
		a.Name, a.Type,
	).Scan(&a.ID, &a.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create account", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, a)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec("DELETE FROM accounts WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func extractID(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return strconv.Atoi(parts[1])
}
