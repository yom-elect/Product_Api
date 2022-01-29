install-swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: install-swagger
	swagger generate spec -o ./swagger.yaml --scan-models

generate_client:
	cd ../client && swagger generate client -f ../product-api/swagger.yaml -A product-api
# cd .. && mkdir client && cd client && go mod init example.com/swagger && swagger generate client -f ../product-api/swagger.yaml -A product-api