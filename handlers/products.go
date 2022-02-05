package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	protos "product.com/product-microservice/product-api/currency"
	"product.com/product-microservice/product-api/data"
)

// Products is a http.Handler
type Products  struct{
	l *log.Logger
	v *data.Validation
	cc protos.CurrencyClient
}

// NewProducts creates a products handler with the given logger
func NewProducts(l*log.Logger, v *data.Validation, cc protos.CurrencyClient) *Products {
	return &Products{l,v, cc}
}

type KeyProduct struct {}

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getProductID(r *http.Request) int {
	// parse the product id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}