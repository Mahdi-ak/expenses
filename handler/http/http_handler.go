package handler

import (
	"encoding/json"
	"errors"
	"expenses/internal/application/expense"
	"expenses/internal/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct {
	service expense.ServiceInterface
}

func NewHandler(service expense.ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

type CreateExpenseRequest struct {
	Amount      string `json:"amount"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Budget      string `json:"budget"`
}

type ExpenseSummaryResponse struct {
	TotalAmount decimal.Decimal `json:"total_amount"`
	TotalBudget decimal.Decimal `json:"total_budget"`
	Difference  decimal.Decimal `json:"difference"`
}

type ExpensesSummaryResponse struct {
	Expenses []domain.Expense       `json:"expenses"`
	Summary  ExpenseSummaryResponse `json:"summary"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateExpenseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	budget, err := decimal.NewFromString(req.Budget)
	if err != nil {
		http.Error(w, "invalid budget", http.StatusBadRequest)
		return
	}

	expense, err := domain.NewExpense(
		amount,
		req.Category,
		req.Description,
		budget,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.Create(r.Context(), *expense)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	var filter domain.Filter

	category := r.URL.Query().Get("category")
	if category != "" {
		filter.Category = &category
	}

	date := r.URL.Query().Get("date")
	if date != "" {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			http.Error(w, "invalid date", http.StatusBadRequest)
			return
		}

		filter.Date = &t
	}

	amount := r.URL.Query().Get("amount")
	if amount != "" {
		parsedAmount, err := decimal.NewFromString(amount)
		if err != nil {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}

		filter.Amount = &parsedAmount
	}

	expenses, err := h.service.GetAll(r.Context(), filter)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(expenses); err != nil {
		return
	}
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	expenseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := h.service.GetByID(
		r.Context(),
		expenseID,
	)

	if err != nil {
		if errors.Is(err, domain.ErrExpenseNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if errors.Is(err, expense.ErrInvalidExpenseID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	expenseID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(
		r.Context(),
		expenseID,
	)

	if err != nil {
		if errors.Is(err, domain.ErrExpenseNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if errors.Is(err, expense.ErrInvalidExpenseID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	expenses, summary, err := h.service.GetExpensesWithSummary(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := ExpensesSummaryResponse{
		Expenses: expenses,
		Summary: ExpenseSummaryResponse{
			TotalAmount: summary.TotalAmount,
			TotalBudget: summary.TotalBudget,
			Difference:  summary.Difference,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
