package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nvloff/f3_payments_service/gen/restapi"
	"github.com/nvloff/f3_payments_service/internal/api"
	"github.com/nvloff/f3_payments_service/internal/domain/mongodb"
)

func main() {
	url := os.Getenv("MONGO_HOST")

	if len(url) == 0 {
		// Default to localhost
		url = "localhost:27017"
	}

	d := mongodb.BuildDomain(mongodb.Config{
		URL:        url,
		Collection: "payments",
	})

	p := api.New(api.Config{
		CommandBus:  d.CommandBus,
		PaymentRepo: d.PaymentRepo,
	})

	// Initiate the http handler, with the objects that are implementing the business logic.
	h, err := restapi.Handler(restapi.Config{
		PaymentsAPI: p,
		Logger:      log.Printf,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting to serve, access server on http://localhost:8080")

	// Run the standard http server
	log.Fatal(http.ListenAndServe(":8080", h))
}
