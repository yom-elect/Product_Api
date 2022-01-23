install-swagger:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: install-swagger
	swagger generate spec -o ./swagger.yaml --scan-models