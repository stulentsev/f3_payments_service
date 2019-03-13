default: build

build:
	@go build -v cmd/f3paymentsd/main.go

swagger:
	script/swagger.sh generate server -t gen -f ./swagger/swagger.yaml -A f3_payments_service

test:
	@go test -v ./...

docker_test:
	@docker-compose run runner make test

cover:
	@go test -v -cover ./...

imports:
	goimports -l -w .

fmt:
	go fmt ./...

.PHONY: fmt test swagger docker_test
