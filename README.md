# Terraform Provider Instana

![CI/CD Build](https://github.com/gessnerfl/terraform-provider-instana/workflows/CICD/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=de.gessnerfl.terraform-provider-instana&metric=alert_status)](https://sonarcloud.io/dashboard?id=de.gessnerfl.terraform-provider-instana)

Terraform provider implementation for the Instana REST API.

Changes Log: **[CHANGELOG.md](https://github.com/gessnerfl/terraform-provider-instana/blob/master/CHANGELOG.md)**

## Documentation

The documentation of the provider can be found on the Github Page <https://gessnerfl.github.io/terraform-provider-instana>.

## Implementation Details

### Testing

 Mocking:
 Tests are co-located in the package next to the implementation. We use gomock (<https://github.com/golang/mock)> for mocking. To generate mocks you need to use the package options to create the mocks in the same package:

```bash
mockgen -source=<source_file> -destination=mocks/<source_package>/<source_file_name>_mocks.go package=<source_package>_mocks -self_package=github.com/gessnerfl/terraform-provider-instana/<source_package>
```

### Release a new version

1. Create a new tag follow semantic versioning approach
2. Update changelog before creating a new release by using [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)
3. Push the tag to the remote to build the new release
