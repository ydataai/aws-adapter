-include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

.PHONY: build fmt vet debug server test mock

build: ### Build
	go build -a cmd

fmt: ### Run go fmt against code
	go fmt ./...

vet: ### Run go vet against code
	go vet ./...

debug: ### Run debug
	go run ./debug

server:	### Run main package
	go run ./cmd/server

test: ### Runs application's tests in verbose mode
	go test -v -cover ./...

mock:
	@ rm mock/*.go || true && \
		$(GOPATH)/bin/mockgen -source=pkg/service/rest_service.go -destination=mock/rest_service_mock.go -package=mock && \
		$(GOPATH)/bin/mockgen -source=pkg/clients/ec2_client.go -destination=mock/ec2_client_mock.go -package=mock && \
		$(GOPATH)/bin/mockgen -source=pkg/clients/service_quota_client.go -destination=mock/service_quota_client_mock.go -package=mock 
		
