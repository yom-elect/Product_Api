package handlers

import (
	"net/http"
	"strconv"

	"example.com/mod/product-api/data"
	"github.com/gorilla/mux"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
// 		201: noContent

// DeleteProducts deletes a product from the data store
func (p *Products) DeleteProducts(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	idString := vars["id"]
	id, _ := strconv.Atoi(idString)
	
	p.l.Println("Handle Delete  PRORUCT", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(w, "Product Not Found", http.StatusNotFound)
	}
 
	if err != nil {
		http.Error(w, "Product Not Found", http.StatusInternalServerError)
	}
}