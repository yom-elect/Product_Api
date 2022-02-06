package handlers

import (
	"net/http"

	"product.com/product-microservice/product-api/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
// 		201: noContentResponse

// DeleteProducts deletes a product from the data store
// Delete handles DELETE requests and removes items from the database
func (p *Products) Delete(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Debug("Deleting record id", id)

	err := p.productDB.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("Unable to delete record id does not exist", "error", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	if err != nil {
		p.l.Error("Unable to delete record", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}