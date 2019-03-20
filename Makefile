    
export CGO_ENABLED:=0
export GO111MODULE=on
#export GOFLAGS=-mod=vendor

VERSION=$(shell git describe --tags --match=v* --always --dirty)

.PHONY: all
all: build test vet lint fmt

.PHONY: ci
ci: build test_with_report vet_with_report lint_with_report fmt release

.PHONY: build
build: clean bin/terraform-provider-instana

bin/terraform-provider-instana:
	@echo "+++++++++++  Run GO Build +++++++++++ "
	@go build -o $@ github.com/gessnerfl/terraform-provider-instana

.PHONY: test
test:
	@echo "+++++++++++  Run GO Test +++++++++++ "
	@go test ./... -cover -v

.PHONY: test_with_report
test_with_report:
	@echo "+++++++++++  Run GO Test (with report) +++++++++++ "
	@mkdir -p output
	@go test ./... -cover -v -coverprofile=output/coverage.out -json > output/unit-test-report.json

.PHONY: vet
vet:
	@echo "+++++++++++  Run GO VET +++++++++++ "
	@go vet -all ./...

.PHONY: vet_with_report
vet_with_report:
	@echo "+++++++++++  Run GO VET (with report) +++++++++++ "
	@mkdir -p output
	@go vet -all ./... 2> output/govet-report.out

.PHONY: lint
lint:
	@echo "+++++++++++  Run GO Lint +++++++++++ "
	@golint -set_exit_status `go list ./...`

.PHONY: lint_with_report
lint_with_report:
	@echo "+++++++++++  Run GO Lint (with report) +++++++++++ "
	@mkdir -p output
	@golint -set_exit_status `go list ./...` > output/golint-report.out

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
	@rm -rf output

.PHONY: sonar
sonar:
	@echo "+++++++++++  Run Sonar Scanner +++++++++++ "
	@sonar-scanner -X -Dsonar.projectVersion=$(VERSION)

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
	@echo "+++++++++++ Build Release $@ +++++++++++ "
	$(GOARGS) go build -o $@ github.com/gessnerfl/terraform-provider-instana