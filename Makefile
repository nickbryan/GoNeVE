-include .env

GOFOLDER=./../../..

# proto: Build the protocol
proto:
	@protoc -I=. --go_out=plugins=grpc:$(GOFOLDER) pkg/network/protocol/domain/*.proto && \
	 protoc -I=. --go_out=plugins=grpc:$(GOFOLDER) pkg/network/protocol/service/*.proto

.PHONY: help

all: help

help: Makefile
	@echo
	@echo "Available commands:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
