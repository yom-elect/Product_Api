package handlers

import (
	"net/http"
	"strconv"

	"example.com/mod/product-api/data"
	"github.com/gorilla/mux"
)


func (p *Products) UpdateProducts(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	idString := vars["id"]
	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}
	
	p.l.Println("Handle PUT  PRORUCT", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)

	if err == data.ErrProductNotFound {
		http.Error(w, "Product Not Found", http.StatusNotFound)
	}
 
	if err != nil {
		http.Error(w, "Product Not Found", http.StatusInternalServerError)
	}
}