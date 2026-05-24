package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go_http-rest-api/db"
	"go_http-rest-api/models"
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	accountID, err := extractAccountID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(
		"SELECT id, account_id, amount, type, description, created_at FROM transactions WHERE account_id = $1 ORDER BY created_at DESC",
		accountID,
	)
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	transactions := []models.Transaction{}
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.AccountID, &t.Amount, &t.Type, &t.Description, &t.CreatedAt); err != nil {
			http.Error(w, "Failed to scan transaction", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}

	writeJSON(w, http.StatusOK, transactions)
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	accountID, err := extractAccountID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Check account exists
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	var t models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if t.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if t.Type != "income" && t.Type != "expense" {
		http.Error(w, "Type must be 'income' or 'expense'", http.StatusBadRequest)
		return
	}

	err = db.DB.QueryRow(
		"INSERT INTO transactions (account_id, amount, type, description) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		accountID, t.Amount, t.Type, t.Description,
	).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	t.AccountID = accountID
	writeJSON(w, http.StatusCreated, t)
}

func GetSummary(w http.ResponseWriter, r *http.Request) {
	accountID, err := extractAccountID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Check account exists
	var exists bool
	err = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)", accountID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	var summary models.Summary
	summary.AccountID = accountID

	err = db.DB.QueryRow(`
		SELECT
			COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) AS total_income,
			COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) AS total_expenses
		FROM transactions
		WHERE account_id = $1
	`, accountID).Scan(&summary.TotalIncome, &summary.TotalExpenses)
	if err != nil {
		http.Error(w, "Failed to fetch summary", http.StatusInternalServerError)
		return
	}

	summary.Balance = summary.TotalIncome - summary.TotalExpenses
	writeJSON(w, http.StatusOK, summary)
}

// --- Helper ---

func extractAccountID(path string) (int, error) {
	// Path: /accounts/3/transactions or /accounts/3/summary
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return 0, strconv.ErrSyntax
	}
	return strconv.Atoi(parts[1])
}
