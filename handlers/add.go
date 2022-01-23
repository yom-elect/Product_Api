package handlers

import (
	"net/http"

	"example.com/mod/product-api/data"
)


func (p *Products) AddProducts(w http.ResponseWriter, r *http.Request){
	p.l.Println("Handle POST  PRORUCT")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	//prod := &data.Product{}
	// err := prod.FromJSON(r.Body)

	// if err != nil {
	// 	http.Error(w, "Unable to unMarshal json", http.StatusBadRequest)
	// }

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(&prod)
}