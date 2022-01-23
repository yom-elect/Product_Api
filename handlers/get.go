package handlers

import (
	"net/http"

	"example.com/mod/product-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products from the database
// responses:
// 		200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET  PRORUCT")
	
	// fetch the products from the data store
	lp := data.GetProducts()

	// serializes the list to JSON
	err := lp.ToJSON(w)
	// d, err := json.Marshal(lp)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}