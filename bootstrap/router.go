package bootstrap

import (
	handler "expenses/handler/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	handler *handler.Handler
}

func NewRouter(handler *handler.Handler) *Router {
	return &Router{
		handler: handler,
	}
}

func (r *Router) Setup() *chi.Mux {

	router := chi.NewRouter()

	router.Post("/expenses", r.handler.Create)
	router.Get("/expenses", r.handler.GetAll)
	router.Get("/expenses/{id}", r.handler.GetByID)
	router.Delete("/expenses/{id}", r.handler.Delete)
	router.Get("/exspenses/summary", r.handler.GetSummary)
	return router
}
