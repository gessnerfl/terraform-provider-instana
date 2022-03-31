# Terraform Provider Instana

![CI/CD Build](https://github.com/gessnerfl/terraform-provider-instana/workflows/CICD/badge.svg)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=de.gessnerfl.terraform-provider-instana&metric=alert_status)](https://sonarcloud.io/dashboard?id=de.gessnerfl.terraform-provider-instana)

Terraform provider implementation for the Instana REST API.

Terraform Registry: <https://registry.terraform.io/providers/gessnerfl/instana/latest>

Changes Log: **[CHANGELOG.md](https://github.com/gessnerfl/terraform-provider-instana/blob/master/CHANGELOG.md)**

## Documentation

The documentation of the provider can be found on the Terraform Registry Page <https://registry.terraform.io/providers/gessnerfl/instana/latest>.

## Implementation Details

### Testing

 Mocking:
 Tests are co-located in the package next to the implementation. We use gomock (<https://github.com/golang/mock>) for mocking. Mocks are 
 created using the *source mode*. All mocks are create in the `mock` package. To generate mocks you can use the helper script 
 `generate-mock-for-file <source-file>` from the root directory of this project.

 Alternatively you can manually execute `mockgen` as follows

```bash
mockgen -source=<source_file> -destination=mocks/<source_file_name>_mocks.go -package=mocks
```

### Release a new version

1. Create a new tag follow semantic versioning approach
2. Update changelog before creating a new release by using [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator)
3. Push the tag to the remote to build the new release
