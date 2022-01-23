// Package classification of Product API
//
// Documentation for Product API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"example.com/mod/product-api/data"
)

type Products  struct{
	l *log.Logger
}

func NewProducts(l*log.Logger) *Products {
	return &Products{l}
}

// func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		p.getProducts(w,r)
// 		return
// 	}

// 	if r.Method == http.MethodPost {
// 		p.addProducts(w,r)
// 		return
// 	}

// 	// handle an update
// 	if r.Method == http.MethodPut {
// 		// expect the id in the URI
// 		reg := regexp.MustCompile(`/([0-9]+)`)
// 		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		
// 		if len(g) != 1 {
// 			http.Error(w, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		if len(g[0]) != 2 {
// 			http.Error(w, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		idString := g[0][1]
// 		id,err := strconv.Atoi(idString)

// 		if err != nil {
// 			http.Error(w, "Invalid URI", http.StatusBadRequest)
// 			return
// 		}

// 		p.updateProducts(id, w , r )
// 		return
// 	}

// 	// catch all 
// 	w.WriteHeader(http.StatusMethodNotAllowed)
// }

type KeyProduct struct {}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(r.Body)

		if err != nil {
			http.Error(w, "Unable to unMarshal json", http.StatusBadRequest)
			return
		} 

		// validate the product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				w, 
				fmt.Sprintf("Error validating product: %s", err),
				 http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(),KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(w, r)
	})
}