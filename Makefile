default: build test

deps:
	go install github.com/hashicorp/terraform

build:
	go build -o terraform-provider-instana

test:
	go test -v

plan:
	@terraform plan