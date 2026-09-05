package handler

import (
	"expenses/internal/domain"
	pb "expenses/proto/pb"
	"time"
)

func toProtoExpense(e domain.Expense) *pb.Expense {
	createdAt := ""
	if !e.CreatedAt.IsZero() {
		createdAt = e.CreatedAt.Format(time.RFC3339)
	}

	amount := e.Amount.String()
	budget := e.Budget.String()

	return pb.Expense_builder{
		Id:          &e.ID,
		Amount:      &amount,
		Category:    &e.Category,
		Description: &e.Description,
		Budget:      &budget,
		CreatedAt:   &createdAt,
	}.Build()
}
