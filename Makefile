.ONESHELL:
.PHONY: clean install tests run all
.EXPORT_ALL_VARIABLES:

clean:
	go clean

install:
	go install ./...

tests:
	go test -v -cover ./...

coverage:
	go test -cover -coverprofile=c.out ./...
	go tool cover -html=c.out

doc:
	godoc

run:
	go run *.go

all: clean install tests run
