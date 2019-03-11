default: build 

deps:
	go install github.com/hashicorp/terraform

build:
	go build

test:
	go test -v

plan:
	@terraform plan