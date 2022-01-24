.DEFAULT_GOAL := swagger

install-swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: install-swagger
	@echo Ensure you have the swagger CLI or this command will fail.
	@echo You can install the swagger CLI with: go get -u github.com/go-swagger/go-swagger/cmd/swagger
	@echo ....
	swagger generate spec -o ./swagger.yaml --scan-models