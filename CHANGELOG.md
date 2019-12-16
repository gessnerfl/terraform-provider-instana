# Change Log

## [v0.6.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.6.0) (2019-12-16)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.5.0...v0.6.0)

**Fixed bugs:**

- Threshold rule support window and rule together [\#24](https://github.com/gessnerfl/terraform-provider-instana/issues/24)

**Closed issues:**

- Support for label/name prefix and suffix [\#22](https://github.com/gessnerfl/terraform-provider-instana/issues/22)
- Update to terraform 0.12.x [\#20](https://github.com/gessnerfl/terraform-provider-instana/issues/20)

**Merged pull requests:**

- \#20: Update project to terraform 0.12.x [\#26](https://github.com/gessnerfl/terraform-provider-instana/pull/26) ([gessnerfl](https://github.com/gessnerfl))
- \#24: Support rollup and window in threshold rule together [\#25](https://github.com/gessnerfl/terraform-provider-instana/pull/25) ([gessnerfl](https://github.com/gessnerfl))

## [v0.5.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.5.0) (2019-10-15)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.4.0...v0.5.0)

**Merged pull requests:**

- \#22: migrate to customizable default name prefix and suffix [\#23](https://github.com/gessnerfl/terraform-provider-instana/pull/23) ([gessnerfl](https://github.com/gessnerfl))

## [v0.4.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.4.0) (2019-10-14)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.3.2...v0.4.0)

**Closed issues:**

- Add support to append terraform managed string [\#19](https://github.com/gessnerfl/terraform-provider-instana/issues/19)

**Merged pull requests:**

- Feature/19 append terraform managed string [\#21](https://github.com/gessnerfl/terraform-provider-instana/pull/21) ([gessnerfl](https://github.com/gessnerfl))

## [v0.3.2](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.3.2) (2019-06-19)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.3.1...v0.3.2)

**Fixed bugs:**

- Terraform provider should not have platform name in executable [\#17](https://github.com/gessnerfl/terraform-provider-instana/issues/17)

**Merged pull requests:**

- \#17: fix binary name [\#18](https://github.com/gessnerfl/terraform-provider-instana/pull/18) ([gessnerfl](https://github.com/gessnerfl))

## [v0.3.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.3.1) (2019-06-19)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.3.0...v0.3.1)

## [v0.3.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.3.0) (2019-06-19)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.2.2...v0.3.0)

**Closed issues:**

- Change release output [\#15](https://github.com/gessnerfl/terraform-provider-instana/issues/15)

**Merged pull requests:**

- Feature/15 release naming [\#16](https://github.com/gessnerfl/terraform-provider-instana/pull/16) ([gessnerfl](https://github.com/gessnerfl))

## [v0.2.2](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.2.2) (2019-05-07)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.2.1...v0.2.2)

**Closed issues:**

- Add support for dashes in tags in match\_specification [\#12](https://github.com/gessnerfl/terraform-provider-instana/issues/12)

**Merged pull requests:**

- Fix application configuration example in README [\#14](https://github.com/gessnerfl/terraform-provider-instana/pull/14) ([steinex](https://github.com/steinex))

## [v0.2.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.2.1) (2019-05-07)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.2.0...v0.2.1)

**Merged pull requests:**

- Feature/12 add support for dashes in identifiers [\#13](https://github.com/gessnerfl/terraform-provider-instana/pull/13) ([gessnerfl](https://github.com/gessnerfl))

## [v0.2.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.2.0) (2019-04-25)
[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.1.0...v0.2.0)

**Closed issues:**

- Severity for Rules are not user friendly [\#8](https://github.com/gessnerfl/terraform-provider-instana/issues/8)
- Add support to manage events [\#7](https://github.com/gessnerfl/terraform-provider-instana/issues/7)
- Migrate to OpenAPI [\#4](https://github.com/gessnerfl/terraform-provider-instana/issues/4)
- Add support to manage groups  [\#3](https://github.com/gessnerfl/terraform-provider-instana/issues/3)
- Add support to create Application Perspectives [\#1](https://github.com/gessnerfl/terraform-provider-instana/issues/1)

**Merged pull requests:**

- Feature/7 events [\#11](https://github.com/gessnerfl/terraform-provider-instana/pull/11) ([gessnerfl](https://github.com/gessnerfl))
- Feature/1 application perspective [\#10](https://github.com/gessnerfl/terraform-provider-instana/pull/10) ([gessnerfl](https://github.com/gessnerfl))
- \#8 Change severity to a user friendly text instead of int codes [\#9](https://github.com/gessnerfl/terraform-provider-instana/pull/9) ([gessnerfl](https://github.com/gessnerfl))
- Feature/3 manage groups [\#6](https://github.com/gessnerfl/terraform-provider-instana/pull/6) ([gessnerfl](https://github.com/gessnerfl))
- Feature/4 migrate to open api [\#5](https://github.com/gessnerfl/terraform-provider-instana/pull/5) ([gessnerfl](https://github.com/gessnerfl))

## [v0.1.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.1.0) (2019-03-14)


\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*