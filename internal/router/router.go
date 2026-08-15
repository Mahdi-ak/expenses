package router

import (
	"expenses/internal/expense"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	handler *expense.Handler
}

func New(handler *expense.Handler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) Setup() *chi.Mux {

	router := chi.NewRouter()

	router.Post("/expenses", r.handler.Create)
	router.Get("/expenses", r.handler.GetAll)
	router.Get("/expenses/{id}", r.handler.GetByID)
	router.Delete("/expenses", r.handler.Delete)

	return router
}
