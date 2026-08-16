package handler

import (
	"encoding/json"
	"expenses/internal/expense"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

type Handler struct {
	service *expense.Service
}

func NewHandler(service *expense.Service) *Handler {
	return &Handler{
		service: service,
	}
}

type CreateExpenseRequest struct {
	Amount      string `json:"amount"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

	var req CreateExpenseRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return
	}

	expense := expense.Expense{
		Amount:      amount,
		Category:    req.Category,
		Description: req.Description,
	}

	result, err := h.service.Create(r.Context(), expense)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {

	var filter expense.Filter

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

	expenses, err := h.service.GetAll(r.Context(), filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(expenses)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	expenseID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	expense, err := h.service.GetByID(
		r.Context(),
		expenseID,
	)

	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	json.NewEncoder(w).Encode(expense)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	expenseID, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	err = h.service.Delete(
		r.Context(),
		expenseID,
	)

	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
