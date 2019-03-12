# terraform-provider-instana
Terraform provider implementation for Instana REST API


# Implementation
 Mocking:
 Tests are colocated in the package next to the implementation. We use pegomock for mocking. To generate mocks you need to use the package options to create the mocks in the same package:

```
pegomock generate -package <impl_package_name> github.com/gessnerfl/terraform-provider-instana/<package> <interface_to_mock>
```