    
export CGO_ENABLED:=0
export GO111MODULE=on
#export GOFLAGS=-mod=vendor

VERSION=$(shell git describe --tags --match=v* --always --dirty)

.PHONY: all
all: build test vet lint fmt

.PHONY: build
build: clean bin/terraform-provider-instana

bin/terraform-provider-instana:
	@go build -o $@ github.com/gessnerfl/terraform-provider-instana

.PHONY: test
test:
	@go test ./... -cover

.PHONY: vet
vet:
	@go vet -all ./...

.PHONY: lint
lint:
	@golint -set_exit_status `go list ./...`

.PHONY: fmt
fmt:
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
	@rm -rf bin
	@rm -rf output

.PHONY: release
release: \
	clean \
	output/plugin-linux-amd64.tar.gz \
	output/plugin-darwin-amd64.tar.gz \
	output/plugin-windows-amd64.tar.gz

output/plugin-%.tar.gz: NAME=terraform-provider-instana-$(VERSION)-$*
output/plugin-%.tar.gz: DEST=output/$(NAME)
output/plugin-%.tar.gz: output/%/terraform-provider-instana
	@mkdir -p $(DEST)
	@cp output/$*/terraform-provider-instana $(DEST)
	@tar zcvf $(DEST).tar.gz -C output $(NAME)

output/linux-amd64/terraform-provider-instana: GOARGS = GOOS=linux GOARCH=amd64
output/darwin-amd64/terraform-provider-instana: GOARGS = GOOS=darwin GOARCH=amd64
output/windows-amd64/terraform-provider-instana: GOARGS = GOOS=windows GOARCH=amd64
output/%/terraform-provider-instana:
	$(GOARGS) go build -o $@ github.com/gessnerfl/terraform-provider-instana