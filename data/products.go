package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	protos "product.com/product-microservice/product-api/currency"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

// Products defines a slice of Product
type Products []*Product

type ProductsDB struct {
	currency 	protos.CurrencyClient
	log 		hclog.Logger
	rates  		map[string]float64
	client  	protos.Currency_SubscribeRatesClient
}

func NewProductsDB(c protos.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{c, l, make(map[string]float64),nil}

	go pb.handleUpdates()

	return pb
}

func (pb *ProductsDB) handleUpdates() {
	sub, err := pb.currency.SubscribeRates(context.Background())
	if err != nil {
		pb.log.Error("Unable to subscribe for rates", "error", err)
		return
	}

	pb.client =  sub

	for {
		// Recv returns a StreamingRateResponse which can contain one of two messages
		// RateResponse or an Error.
		// We need to handle each case separately
		rr, err := sub.Recv()

		// handle connection errors
		// this is normally terminal requires a reconnect
		if err != nil {
			pb.log.Error("Error while waiting for message", "error", err)
			return
		}


		// handle a returned error message
		if ge := rr.GetError(); ge != nil {
			sre := status.FromProto(ge)

			if sre.Code() == codes.InvalidArgument {
				errDetails := ""
				// get the RateRequest serialized in the error response
				// Details is a collection but we are only returning a single item
				if d := sre.Details(); len(d) > 0 {
					pb.log.Error("Deets", "d", d)
					if rr, ok := d[0].(*protos.RateRequest); ok {
						errDetails = fmt.Sprintf("base: %s destination: %s", rr.GetBase().String(), rr.GetDestination().String())
					}
				}

				pb.log.Error("Received error from currency service rate subscription", "error", ge.GetMessage(), "details", errDetails)
			}
		}


		// handle a rate response
		if rr := rr.GetRateResponse(); rr != nil {
			pb.log.Info("Recieved updated rate from server", "dest", rr.GetDestination().String())
			pb.rates[rr.Destination.String()] = rr.Rate
		}
	}
}

// GetProducts returns all products from the database
func (pb *ProductsDB) GetProducts(currency string) (Products, error) {
	if (currency == "") {
		return productList, nil
	}

	resp, err := pb.getRate(currency)
	if err != nil {
		pb.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price *= resp
		pr = append(pr,&np)
	}

	return pr, nil
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (pb *ProductsDB) GetProductByID(id int, currency string) (*Product, error) {
	i := pb.findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if (currency == "") {
		return productList[i], nil
	}

	rate, err := pb.getRate(currency)
	if err != nil {
		pb.log.Error("Unable to get rate", "currency", currency, "error", err)
		return nil, err
	}

	np := *productList[i]
	np.Price *= rate

	return &np, nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (pb *ProductsDB) UpdateProduct(p Product) error {
	i := pb.findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &p

	return nil
}

// AddProduct adds a new product to the database
func (pb *ProductsDB) AddProduct(p Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, &p)
}

// DeleteProduct deletes a product from the database
func (pb *ProductsDB) DeleteProduct(id int) error {
	i := pb.findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func (pb *ProductsDB) findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (pb *ProductsDB) getRate(destinationRate string) (float64, error) {
	// if cached returns
	if r, ok := pb.rates[destinationRate]; ok {
		return r, nil
	}

	rr := &protos.RateRequest{
		Base: protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[destinationRate]),
	}

	// get initial rate
	resp, err := pb.currency.GetRate(context.Background(), rr)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0].(*protos.RateRequest)

			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("unable to get rate from server, destination and base currencies cannot be the same, base: %s, dest: %s", md.Base.String(), md.Destination.String())
			}

			return -1, fmt.Errorf("unable to get rate from server, base: %s, dest: %s", md.Base.String(), md.Destination.String())
		}
		return -1, err
	}

	pb.rates[destinationRate] = resp.Rate // update cached

	// subscribe for updates
	pb.client.Send(rr)
	
	return resp.Rate, err
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	&Product{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}