package transactions

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/skykosiner/gobutar/pkg/budget"
	"github.com/skykosiner/gobutar/pkg/items"
	"github.com/skykosiner/gobutar/pkg/payee"
	"github.com/skykosiner/gobutar/pkg/templates"
)

type Transaction struct {
	ID           int     `sql:"id" json:"id"`
	PurchaseDate string  `sql:"purchase_date" json:"purchase_date"`
	Payee        string  `sql:"payee" json:"payee"`
	ItemID       string  `sql:"item_id" json:"item_id"`
	Outflow      float64 `sql:"outflow" json:"outflow"`
	Inflow       float64 `sql:"inflow" json:"inflow"`
	ItemName     string  `json:"item_name"`
}

func GetTransactions(db *sql.DB) ([]Transaction, error) {
	var transactions []Transaction

	rows, err := db.Query(`
		SELECT
			t.id AS transaction_id,
			t.purchase_date AS purchase_date,
			t.payee as payee,
			t.outflow AS outflow,
			t.inflow AS inflow,
			t.item_id AS spent_item_id,
			i.name AS item_name
		FROM
			transactions t
		LEFT JOIN
			items i ON t.item_id = i.id;
		`)

	if err != nil {
		return transactions, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.PurchaseDate, &transaction.Payee, &transaction.Outflow, &transaction.Inflow, &transaction.ItemID, &transaction.ItemName); err != nil {
			return transactions, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func NewTransaction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newTransaction struct {
			Payee        string `json:"payee"`
			PurchaseDate string `json:"purchase_date"`
			ItemID       string `json:"item_id"`
			Outflow      string `json:"outflow"`
			Inflow       string `json:"inflow"`
		}
		if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
			slog.Warn("Possible bad user input or error adding new transaction.", "error", err, "r.body", r.Body)
			http.Error(w, "Please make sure you entered all the information correctly.", http.StatusBadRequest)
			return
		}

		if err := payee.ProcessPayee(db, newTransaction.Payee); err != nil {
			slog.Error("Error processing payee info", "error", err, "new transaction", newTransaction)
			http.Error(w, "Sorry there was an error, please try again.", http.StatusInternalServerError)
			return
		}

		outflow, _ := strconv.ParseFloat(newTransaction.Outflow, 64)
		inflow, _ := strconv.ParseFloat(newTransaction.Inflow, 64)

		_, err := db.Exec("INSERT INTO transactions (payee, purchase_date, item_id, outflow, inflow) VALUES (?,?,?,?,?)", newTransaction.Payee, newTransaction.PurchaseDate, newTransaction.ItemID, outflow, inflow)
		if err != nil {
			slog.Error("Error creating new transaction in db", "error", err, "new transaction", newTransaction)
			http.Error(w, "Sorry there was an error, please try again.", http.StatusInternalServerError)
			return
		}

		if inflow > 0 {
			_, err := db.Exec("UPDATE budget SET unallocated = unallocated + ? AND current_balance = current_balance + ?", inflow, inflow)
			if err != nil {
				slog.Error("Error updating budget database.", "error", err, "new transaction", newTransaction)
				http.Error(w, "Sorry there was an error, please try again.", http.StatusInternalServerError)
				return
			}
		} else if outflow > 0 {
			b, err := budget.NewBudget(db)
			if err != nil {
				slog.Error("Error updating budget database.", "error", err, "new transaction", newTransaction)
				http.Error(w, "Sorry there was an error, please try again.", http.StatusInternalServerError)
				return
			}

			if outflow > b.CurrentBalance {
				http.Error(w,"You're outflow is greater then your current balance", http.StatusBadRequest)
				return
			}

			_, err = db.Exec("UPDATE budget SET current_balance = current_balance - ?", outflow)
			if err != nil {
				slog.Error("Error updating budget database.", "error", err, "new transaction", newTransaction)
				http.Error(w, "Sorry there was an error, please try again.", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}

func DeleteTransaction(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if len(id) == 0 {
			http.Error(w, "Please provide an id of the transaction to delete.", http.StatusBadRequest)
			return
		}

		_, err := db.Exec("DELETE FROM transactions WHERE id = ?", id)
		if err != nil {
			slog.Error("Depleting transaction.", "error", err, "transaction id", id)
			http.Error(w, "Error deleting that transaction please try again.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func SendNewTransactionForm(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newTransaction struct {
			Payees []string
			Items  []struct {
				ID   string
				Name string
			}
		}

		payees, err := payee.GetPayees(db)
		if err != nil {
			slog.Error("Error getting payees", "error", err)
			http.Error(w, "Sorry there was an error please try again.", http.StatusInternalServerError)
			return
		}

		for _, p := range payees {
			newTransaction.Payees = append(newTransaction.Payees, p.Name)
		}

		it, err := items.GetItems(db)
		if err != nil {
			slog.Error("Error getting items", "error", err)
			http.Error(w, "Sorry there was an error please try again.", http.StatusInternalServerError)
			return
		}

		for _, i := range it {
			newTransaction.Items = append(newTransaction.Items, struct {
				ID   string
				Name string
			}{
				i.ID,
				i.Name,
			})
		}

		templates.RenderTemplate(w, "new-transaction", newTransaction)
	}
}
