    
export CGO_ENABLED:=0
export GO111MODULE=on
#export GOFLAGS=-mod=vendor

VERSION=$(shell git describe --tags --match "v*" --always --dirty)

.PHONY: all
all: build test vet lint fmt

.PHONY: build
build: clean bin/terraform-provider-instana

bin/terraform-provider-instana:
	@echo "+++++++++++  Run GO Build +++++++++++ "
	@go build -o $@ github.com/gessnerfl/terraform-provider-instana

.PHONY: test
test:
	@echo "+++++++++++  Run GO Test +++++++++++ "
	@go test ./... -cover

.PHONY: gosec
gosec:
	@echo "+++++++++++  Run GO SEC +++++++++++ "
	@gosec ./... 

.PHONY: vet
vet:
	@echo "+++++++++++  Run GO VET +++++++++++ "
	@go vet -all ./...

.PHONY: lint
lint:
	@echo "+++++++++++  Run GO Lint +++++++++++ "
	@golint -set_exit_status `go list ./...`

.PHONY: fmt
fmt:
	@echo "+++++++++++  Run GO FMT +++++++++++ "
	@test -z $$(go fmt ./...) 

.PHONY: update
update:
	@GOFLAGS="" go get -u
	@go mod tidy

.PHONY: vendor
vendor:
	@go mod vendor

.PHONY: clean
clean:
	@echo "+++++++++++  Clean up project +++++++++++ "
	@rm -rf bin