    
export CGO_ENABLED:=0
export GO111MODULE=on
#export GOFLAGS=-mod=vendor

ifdef TRAVIS_TAG
	VERSION=$(TRAVIS_TAG)
else
	VERSION=$(shell git describe --tags --match "v*" --always --dirty)
endif

.PHONY: all
all: build test vet lint fmt

.PHONY: ci
ci: build test_with_report vet_with_report lint_with_report fmt sonar release

.PHONY: build
build: clean bin/terraform-provider-instana

bin/terraform-provider-instana:
	@echo "+++++++++++  Run GO Build +++++++++++ "
	@go build -o $@ github.com/gessnerfl/terraform-provider-instana

.PHONY: test
test:
	@echo "+++++++++++  Run GO Test +++++++++++ "
	@go test ./... -cover

.PHONY: test_with_report
test_with_report:
	@echo "+++++++++++  Run GO Test (with report) +++++++++++ "
	@mkdir -p output
	@go test ./... -cover -coverprofile=output/coverage.out -json > output/unit-test-report.json
	#@echo "Result:"
	#@cat output/unit-test-report.json

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
ifdef TRAVIS_TAG
	@sonar-scanner -X -Dsonar.projectVersion=$(VERSION) -Dsonar.branch.name=master
endif

ifeq ($(TRAVIS_BRANCH), master)
	@sonar-scanner -X -Dsonar.projectVersion=$(VERSION)
else ifdef TRAVIS_PULL_REQUEST_BRANCH
	@sonar-scanner -X -Dsonar.projectVersion=$(VERSION) -Dsonar.pullrequest.key=$(TRAVIS_PULL_REQUEST) -Dsonar.pullrequest.branch=$(TRAVIS_PULL_REQUEST_BRANCH) -Dsonar.pullrequest.base=$(TRAVIS_BRANCH)
else
	@sonar-scanner -X -Dsonar.projectVersion=$(VERSION) -Dsonar.branch.name=$(TRAVIS_BRANCH) -Dsonar.branch.target=master
endif

.PHONY: release
release: \
	clean \
	output/plugin-linux-amd64.zip \
	output/plugin-darwin-amd64.zip \
	output/plugin-windows-amd64.zip

output/plugin-%.zip: TARGET_PLATFORM=$*
output/plugin-%.zip: NAME=terraform-provider-instana_$(VERSION)_$(TARGET_PLATFORM)
output/plugin-%.zip: DEST=output/$(NAME)
output/plugin-%.zip: output/%/terraform-provider-instana
	@zip -j $(DEST).zip output/$(TARGET_PLATFORM)/$(NAME)

output/linux-amd64/terraform-provider-instana: GOARGS = GOOS=linux GOARCH=amd64
output/darwin-amd64/terraform-provider-instana: GOARGS = GOOS=darwin GOARCH=amd64
output/windows-amd64/terraform-provider-instana: GOARGS = GOOS=windows GOARCH=amd64
output/%/terraform-provider-instana:
	@echo "+++++++++++ Build Release $@ +++++++++++ "
	$(GOARGS) go build -o $@_$(VERSION)_$(TARGET_PLATFORM) github.com/gessnerfl/terraform-provider-instana