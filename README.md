# terraform-provider-instana
Terraform provider implementation for Instana REST API


# Implementation
 Mocking:
 Tests are colocated in the package next to the implementation. We use gomock (https://github.com/golang/mock) for mocking. To generate mocks you need to use the package options to create the mocks in the same package:

```
mockgen -source=<source_file> -destination=mocks/<source_package>/<source_file_name>_mocks.go package=<source_package>_mocks -self_package=github.com/gessnerfl/terraform-provider-instana/<source_package>
```