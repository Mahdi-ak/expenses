package bootstrap

import (
	handler "expenses/handler/http"
	"expenses/internal/application/expense"
	"log"
	"net/http"
)

func StartHTTPServer(port string, service expense.ServiceInterface) {
	h := handler.NewHandler(service)
	r := NewRouter(h)

	log.Println("HTTP server running", port)

	if err := http.ListenAndServe(port, r.Setup()); err != nil {
		log.Fatal(err)
	}
}
