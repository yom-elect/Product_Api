package handlers

import (
	"net/http"

	"example.com/mod/product-api/data"
)

func (p *Products) GetProducts(w http.ResponseWriter, r *http.Request){
	p.l.Println("Handle GET  PRORUCT")

	lp := data.GetProducts()
	err := lp.ToJSON(w)
	// d, err := json.Marshal(lp)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}