# Changelog

## [v1.6.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.6.1) (2023-07-19)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.6.0...v1.6.1)

**Closed issues:**

- Support to manage Synthetic endpoints [\#78](https://github.com/gessnerfl/terraform-provider-instana/issues/78)

**Merged pull requests:**

- Support env vars in provider [\#149](https://github.com/gessnerfl/terraform-provider-instana/pull/149) ([mzwennes](https://github.com/mzwennes))

## [v1.6.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.6.0) (2023-07-12)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.5.3...v1.6.0)

**Closed issues:**

- Removing the instana terraform code, resulted in an error when removing the resources [\#139](https://github.com/gessnerfl/terraform-provider-instana/issues/139)
- Tag call.http.header is independent of source or destination [\#137](https://github.com/gessnerfl/terraform-provider-instana/issues/137)

**Merged pull requests:**

- Feature/147 syntetic test and location [\#148](https://github.com/gessnerfl/terraform-provider-instana/pull/148) ([gessnerfl](https://github.com/gessnerfl))
- Add Synthetic monitors Resource & Synthetic Location Datasource [\#147](https://github.com/gessnerfl/terraform-provider-instana/pull/147) ([mzwennes](https://github.com/mzwennes))
- Update OpenAPI specification to latest version [\#146](https://github.com/gessnerfl/terraform-provider-instana/pull/146) ([mzwennes](https://github.com/mzwennes))
- fixed typos [\#144](https://github.com/gessnerfl/terraform-provider-instana/pull/144) ([CayoM](https://github.com/CayoM))

## [v1.5.3](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.5.3) (2022-11-01)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.5.2...v1.5.3)

**Fixed bugs:**

-  Application\_alert\_configuration Resource is duplicated each time terraform code is applied [\#141](https://github.com/gessnerfl/terraform-provider-instana/issues/141)

**Closed issues:**

- Updating deleted event config ID [\#140](https://github.com/gessnerfl/terraform-provider-instana/issues/140)
- Make id not auto generated so I can pass it to instana\_alerting\_config [\#138](https://github.com/gessnerfl/terraform-provider-instana/issues/138)

**Merged pull requests:**

- \#141: Fix mapping of empty TagFilterExpression [\#142](https://github.com/gessnerfl/terraform-provider-instana/pull/142) ([gessnerfl](https://github.com/gessnerfl))

## [v1.5.2](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.5.2) (2022-07-26)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.5.1...v1.5.2)

**Implemented enhancements:**

- allow to skip certificate validation [\#134](https://github.com/gessnerfl/terraform-provider-instana/issues/134)
- Add gosec to project [\#67](https://github.com/gessnerfl/terraform-provider-instana/issues/67)
- Maintenance configuration [\#31](https://github.com/gessnerfl/terraform-provider-instana/issues/31)

**Merged pull requests:**

- \#67: activate gosec via golangci-lint [\#136](https://github.com/gessnerfl/terraform-provider-instana/pull/136) ([gessnerfl](https://github.com/gessnerfl))
- Feature/134 skip tls verify [\#135](https://github.com/gessnerfl/terraform-provider-instana/pull/135) ([gessnerfl](https://github.com/gessnerfl))

## [v1.5.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.5.1) (2022-06-22)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.5.0...v1.5.1)

**Implemented enhancements:**

- Tag filters does not support bracketing [\#131](https://github.com/gessnerfl/terraform-provider-instana/issues/131)

**Merged pull requests:**

- \#131: Fix mapping of nested bracketed expressions from instana model … [\#133](https://github.com/gessnerfl/terraform-provider-instana/pull/133) ([gessnerfl](https://github.com/gessnerfl))
- Feature/131 bracketing in tagfilters [\#132](https://github.com/gessnerfl/terraform-provider-instana/pull/132) ([gessnerfl](https://github.com/gessnerfl))

## [v1.5.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.5.0) (2022-06-01)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.4.0...v1.5.0)

**Implemented enhancements:**

- Support for WebSite and Application smart alerts [\#102](https://github.com/gessnerfl/terraform-provider-instana/issues/102)

**Closed issues:**

- switch to golangci-lint [\#128](https://github.com/gessnerfl/terraform-provider-instana/issues/128)
- Dashboard resource  [\#106](https://github.com/gessnerfl/terraform-provider-instana/issues/106)

**Merged pull requests:**

- Feature/106 dashboards [\#130](https://github.com/gessnerfl/terraform-provider-instana/pull/130) ([gessnerfl](https://github.com/gessnerfl))
- Feature/128 switch to golangci lint [\#129](https://github.com/gessnerfl/terraform-provider-instana/pull/129) ([gessnerfl](https://github.com/gessnerfl))

## [v1.4.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.4.0) (2022-04-20)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.3.0...v1.4.0)

**Implemented enhancements:**

- Support for website alert config [\#126](https://github.com/gessnerfl/terraform-provider-instana/issues/126)

**Merged pull requests:**

- Feature/126 website alert config [\#127](https://github.com/gessnerfl/terraform-provider-instana/pull/127) ([gessnerfl](https://github.com/gessnerfl))

## [v1.3.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.3.0) (2022-04-15)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.2.1...v1.3.0)

**Implemented enhancements:**

- Support Application Alert Configs [\#74](https://github.com/gessnerfl/terraform-provider-instana/issues/74)

**Merged pull requests:**

- Feature/74 application alert config [\#125](https://github.com/gessnerfl/terraform-provider-instana/pull/125) ([gessnerfl](https://github.com/gessnerfl))
- Fix heading for Alerting Channel Webhook Resource [\#124](https://github.com/gessnerfl/terraform-provider-instana/pull/124) ([ppuschmann](https://github.com/ppuschmann))

## [v1.2.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.2.1) (2022-03-08)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.2.0...v1.2.1)

**Merged pull requests:**

- \#122: trim whitespaces of api token and endpoint url of provider config [\#123](https://github.com/gessnerfl/terraform-provider-instana/pull/123) ([gessnerfl](https://github.com/gessnerfl))

## [v1.2.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.2.0) (2022-01-10)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.1.5...v1.2.0)

**Implemented enhancements:**

- ARM64 support [\#119](https://github.com/gessnerfl/terraform-provider-instana/issues/119)
- Address vulnerabilities [\#110](https://github.com/gessnerfl/terraform-provider-instana/issues/110)

**Merged pull requests:**

- \#119: add arm64 support [\#120](https://github.com/gessnerfl/terraform-provider-instana/pull/120) ([gessnerfl](https://github.com/gessnerfl))

## [v1.1.5](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.1.5) (2021-12-28)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.1.4...v1.1.5)

**Fixed bugs:**

-  Instana Terraform provider has wrong documentation paths on registry.terraform.io [\#117](https://github.com/gessnerfl/terraform-provider-instana/issues/117)

**Closed issues:**

- at least two elements are expected for a tag filter expression [\#116](https://github.com/gessnerfl/terraform-provider-instana/issues/116)

**Merged pull requests:**

- \#117: fix naming of resources and data sources in documentation [\#118](https://github.com/gessnerfl/terraform-provider-instana/pull/118) ([gessnerfl](https://github.com/gessnerfl))

## [v1.1.4](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.1.4) (2021-10-05)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.1.3...v1.1.4)

**Breaking changes:**

- Tag Filter does not support unary expressions for tags [\#114](https://github.com/gessnerfl/terraform-provider-instana/issues/114)

**Fixed bugs:**

- Provider vs. Instana-Server v195: "Error: value not allowed for unary operator expression" [\#100](https://github.com/gessnerfl/terraform-provider-instana/issues/100)

**Merged pull requests:**

- Bugfix/114 tags for unary expressions [\#115](https://github.com/gessnerfl/terraform-provider-instana/pull/115) ([gessnerfl](https://github.com/gessnerfl))

## [v1.1.3](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.1.3) (2021-10-04)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.1.2...v1.1.3)

**Merged pull requests:**

- \#111: value check for tag filter should consider value parameter only [\#113](https://github.com/gessnerfl/terraform-provider-instana/pull/113) ([gessnerfl](https://github.com/gessnerfl))

## [v1.1.2](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.1.2) (2021-10-04)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.1.1...v1.1.2)

**Fixed bugs:**

- Terraform Provider is considering tag keys as value [\#111](https://github.com/gessnerfl/terraform-provider-instana/issues/111)

**Merged pull requests:**

- \#111: fix handling of tag keys in unary expressions [\#112](https://github.com/gessnerfl/terraform-provider-instana/pull/112) ([gessnerfl](https://github.com/gessnerfl))

## [v1.1.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.1.1) (2021-09-21)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.1.0...v1.1.1)

**Closed issues:**

- Documentation for RBAC groups is missing [\#108](https://github.com/gessnerfl/terraform-provider-instana/issues/108)

**Merged pull requests:**

- Update changelog after release 1.1.0 [\#109](https://github.com/gessnerfl/terraform-provider-instana/pull/109) ([gessnerfl](https://github.com/gessnerfl))

## [v1.1.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.1.0) (2021-09-20)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v1.0.0...v1.1.0)

**Implemented enhancements:**

- Add support to manage groups [\#93](https://github.com/gessnerfl/terraform-provider-instana/issues/93)

**Closed issues:**

- Support tag filter for application config [\#104](https://github.com/gessnerfl/terraform-provider-instana/issues/104)
- invalid preface while using terraform provider for imports and creating slack channel [\#99](https://github.com/gessnerfl/terraform-provider-instana/issues/99)

**Merged pull requests:**

- Feature/104 tag filter [\#107](https://github.com/gessnerfl/terraform-provider-instana/pull/107) ([gessnerfl](https://github.com/gessnerfl))
- Feature/93 groups [\#103](https://github.com/gessnerfl/terraform-provider-instana/pull/103) ([gessnerfl](https://github.com/gessnerfl))

## [v1.0.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v1.0.0) (2021-04-25)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.11.0...v1.0.0)

**Breaking changes:**

- Resource instana\_user\_role not working anymore [\#92](https://github.com/gessnerfl/terraform-provider-instana/issues/92)

**Implemented enhancements:**

- Support Import for resources [\#95](https://github.com/gessnerfl/terraform-provider-instana/issues/95)
- Support API Tokens [\#94](https://github.com/gessnerfl/terraform-provider-instana/issues/94)
- Migrate to Terraform Plugin SDK v2 [\#87](https://github.com/gessnerfl/terraform-provider-instana/issues/87)

**Merged pull requests:**

- Feature/95 import support [\#98](https://github.com/gessnerfl/terraform-provider-instana/pull/98) ([gessnerfl](https://github.com/gessnerfl))
- Feature/94 api tokens [\#97](https://github.com/gessnerfl/terraform-provider-instana/pull/97) ([gessnerfl](https://github.com/gessnerfl))
- Feature/87 sdkv2 [\#96](https://github.com/gessnerfl/terraform-provider-instana/pull/96) ([gessnerfl](https://github.com/gessnerfl))

## [v0.11.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.11.0) (2021-02-09)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.10.1...v0.11.0)

**Implemented enhancements:**

- Use structs to encapsulate resource implementation [\#90](https://github.com/gessnerfl/terraform-provider-instana/issues/90)
- Support Builtin Events as a data source e.g. for built in events [\#86](https://github.com/gessnerfl/terraform-provider-instana/issues/86)

**Fixed bugs:**

- release 0.10.0 with mismatching checksums/SHA [\#88](https://github.com/gessnerfl/terraform-provider-instana/issues/88)

**Merged pull requests:**

- Feature/90 resources as structs [\#91](https://github.com/gessnerfl/terraform-provider-instana/pull/91) ([gessnerfl](https://github.com/gessnerfl))
- Feature/86 data source builtin events [\#89](https://github.com/gessnerfl/terraform-provider-instana/pull/89) ([gessnerfl](https://github.com/gessnerfl))

## [v0.10.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.10.1) (2021-02-01)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.10.0...v0.10.1)

**Implemented enhancements:**

- Add support for application config -\> matchspecification -\> entity? [\#84](https://github.com/gessnerfl/terraform-provider-instana/issues/84)
- Add support to website monitoring [\#72](https://github.com/gessnerfl/terraform-provider-instana/issues/72)

## [v0.10.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.10.0) (2021-02-01)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.9.2...v0.10.0)

**Implemented enhancements:**

- Improve system to support rest resources with separated create and update endpoints [\#79](https://github.com/gessnerfl/terraform-provider-instana/issues/79)

**Closed issues:**

- Support SLI Configuration [\#77](https://github.com/gessnerfl/terraform-provider-instana/issues/77)

**Merged pull requests:**

- Feature/84 support entity origin in app config match spec [\#85](https://github.com/gessnerfl/terraform-provider-instana/pull/85) ([gessnerfl](https://github.com/gessnerfl))
- Feature/72 website monitoring [\#82](https://github.com/gessnerfl/terraform-provider-instana/pull/82) ([gessnerfl](https://github.com/gessnerfl))
- Feature/77 sli configuration [\#81](https://github.com/gessnerfl/terraform-provider-instana/pull/81) ([jonassiefker](https://github.com/jonassiefker))
- Feature/79 rest resource improvement [\#80](https://github.com/gessnerfl/terraform-provider-instana/pull/80) ([gessnerfl](https://github.com/gessnerfl))

## [v0.9.2](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.9.2) (2020-11-10)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.9.1...v0.9.2)

**Implemented enhancements:**

- Match Specification for Application Configuration is to strict [\#75](https://github.com/gessnerfl/terraform-provider-instana/issues/75)
- Mocks are out of sync / readme does not match [\#71](https://github.com/gessnerfl/terraform-provider-instana/issues/71)

**Merged pull requests:**

- \#75: Allow slashes in identifiers [\#76](https://github.com/gessnerfl/terraform-provider-instana/pull/76) ([gessnerfl](https://github.com/gessnerfl))
- Feature/71 update mocks and documentation for mocking [\#73](https://github.com/gessnerfl/terraform-provider-instana/pull/73) ([gessnerfl](https://github.com/gessnerfl))

## [v0.9.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.9.1) (2020-10-14)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.9.0...v0.9.1)

**Fixed bugs:**

- Breaking Change of EQUAL Condition Operator in Release 188 [\#69](https://github.com/gessnerfl/terraform-provider-instana/issues/69)

**Merged pull requests:**

- Bugfix/69 equal condition operator [\#70](https://github.com/gessnerfl/terraform-provider-instana/pull/70) ([gessnerfl](https://github.com/gessnerfl))

## [v0.9.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.9.0) (2020-08-27)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.8.2...v0.9.0)

**Implemented enhancements:**

- Adjust user role to changes in REST API [\#65](https://github.com/gessnerfl/terraform-provider-instana/issues/65)
- Migrate to Github Actions [\#63](https://github.com/gessnerfl/terraform-provider-instana/issues/63)
- Add github pages based documentation [\#61](https://github.com/gessnerfl/terraform-provider-instana/issues/61)
- Allow Instana API Representation for Threshold Rule Matching Operator also in Terraform [\#54](https://github.com/gessnerfl/terraform-provider-instana/issues/54)
- Add boundary scope to application configuration [\#51](https://github.com/gessnerfl/terraform-provider-instana/issues/51)
- Add new operators to application config [\#50](https://github.com/gessnerfl/terraform-provider-instana/issues/50)
- Add Event Type `agent_monitoring_issue` [\#49](https://github.com/gessnerfl/terraform-provider-instana/issues/49)
- Add support for custom event configurations of rule type "Threshold Rule using a dynamic built-in metric by pattern" [\#47](https://github.com/gessnerfl/terraform-provider-instana/issues/47)

**Fixed bugs:**

- Breaking change of matching\_operator in Instana API [\#48](https://github.com/gessnerfl/terraform-provider-instana/issues/48)
- Idempotence in "instana\_alerting\_config" only with priorities/event\_filter\_event\_types in alphabetical order [\#43](https://github.com/gessnerfl/terraform-provider-instana/issues/43)

**Merged pull requests:**

- \#65: add new and remove deprecated fields of user role [\#66](https://github.com/gessnerfl/terraform-provider-instana/pull/66) ([gessnerfl](https://github.com/gessnerfl))
- Feature/63 GitHub actions [\#64](https://github.com/gessnerfl/terraform-provider-instana/pull/64) ([gessnerfl](https://github.com/gessnerfl))
- Create github pages based documentation [\#62](https://github.com/gessnerfl/terraform-provider-instana/pull/62) ([gessnerfl](https://github.com/gessnerfl))
- Bugfix/47 fix dynamic built in metrics [\#60](https://github.com/gessnerfl/terraform-provider-instana/pull/60) ([gessnerfl](https://github.com/gessnerfl))
- Feature/51 app config boundary scope [\#59](https://github.com/gessnerfl/terraform-provider-instana/pull/59) ([gessnerfl](https://github.com/gessnerfl))
- \#49 add new supported event type agent\_monitoring\_issue [\#57](https://github.com/gessnerfl/terraform-provider-instana/pull/57) ([gessnerfl](https://github.com/gessnerfl))
- \#50 Add additional comparison operators for application configurations [\#56](https://github.com/gessnerfl/terraform-provider-instana/pull/56) ([gessnerfl](https://github.com/gessnerfl))
- \#54: allow instana representation in addition to terraform for matchi… [\#55](https://github.com/gessnerfl/terraform-provider-instana/pull/55) ([gessnerfl](https://github.com/gessnerfl))
- Feature/47 custom events dynamic filter [\#53](https://github.com/gessnerfl/terraform-provider-instana/pull/53) ([gessnerfl](https://github.com/gessnerfl))
- Feature/48 fix api changes [\#52](https://github.com/gessnerfl/terraform-provider-instana/pull/52) ([gessnerfl](https://github.com/gessnerfl))

## [v0.8.2](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.8.2) (2020-03-02)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.8.1...v0.8.2)

**Implemented enhancements:**

- Consistently use testify assert [\#41](https://github.com/gessnerfl/terraform-provider-instana/issues/41)
- Allow up to 1024 rule-ids per instana\_alerting\_config [\#39](https://github.com/gessnerfl/terraform-provider-instana/issues/39)

**Fixed bugs:**

- Downstream integration Ids not supported anymore for Custom Event [\#44](https://github.com/gessnerfl/terraform-provider-instana/issues/44)

**Merged pull requests:**

- Bugfix/43 order of event types [\#46](https://github.com/gessnerfl/terraform-provider-instana/pull/46) ([gessnerfl](https://github.com/gessnerfl))
- \#44: Remove downstream integration IDs from custom event specs [\#45](https://github.com/gessnerfl/terraform-provider-instana/pull/45) ([gessnerfl](https://github.com/gessnerfl))
- \#41: use testify instead of plain go if checks in tests [\#42](https://github.com/gessnerfl/terraform-provider-instana/pull/42) ([gessnerfl](https://github.com/gessnerfl))

## [v0.8.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.8.1) (2020-02-21)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.8.0...v0.8.1)

**Merged pull requests:**

- \#39 Allow 1024 rule-ids per instana\_alerting\_config [\#40](https://github.com/gessnerfl/terraform-provider-instana/pull/40) ([ppuschmann](https://github.com/ppuschmann))

## [v0.8.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.8.0) (2020-02-14)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.7.0...v0.8.0)

**Implemented enhancements:**

- Migrate to new resource approach [\#35](https://github.com/gessnerfl/terraform-provider-instana/issues/35)
- REST Client should support retries [\#32](https://github.com/gessnerfl/terraform-provider-instana/issues/32)
- Alerting Configuration [\#30](https://github.com/gessnerfl/terraform-provider-instana/issues/30)
- Alerting Channel Configuration [\#29](https://github.com/gessnerfl/terraform-provider-instana/issues/29)

**Merged pull requests:**

- Feature/32 rest throttling [\#38](https://github.com/gessnerfl/terraform-provider-instana/pull/38) ([gessnerfl](https://github.com/gessnerfl))
- Feature/30 alerting configuration [\#37](https://github.com/gessnerfl/terraform-provider-instana/pull/37) ([gessnerfl](https://github.com/gessnerfl))
- Feature/35 new resource approach [\#36](https://github.com/gessnerfl/terraform-provider-instana/pull/36) ([gessnerfl](https://github.com/gessnerfl))
- Feature/29 altering channels [\#34](https://github.com/gessnerfl/terraform-provider-instana/pull/34) ([gessnerfl](https://github.com/gessnerfl))

## [v0.7.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.7.0) (2019-12-17)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.6.0...v0.7.0)

**Implemented enhancements:**

- Add support for Entity Verification Rule Type [\#27](https://github.com/gessnerfl/terraform-provider-instana/issues/27)

**Merged pull requests:**

- Feature/27 entity verification events [\#28](https://github.com/gessnerfl/terraform-provider-instana/pull/28) ([gessnerfl](https://github.com/gessnerfl))

## [v0.6.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.6.0) (2019-12-16)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.5.0...v0.6.0)

**Implemented enhancements:**

- Update to terraform 0.12.x [\#20](https://github.com/gessnerfl/terraform-provider-instana/issues/20)

**Fixed bugs:**

- Threshold rule support window and rule together [\#24](https://github.com/gessnerfl/terraform-provider-instana/issues/24)

**Merged pull requests:**

- \#20: Update project to terraform 0.12.x [\#26](https://github.com/gessnerfl/terraform-provider-instana/pull/26) ([gessnerfl](https://github.com/gessnerfl))
- \#24: Support rollup and window in threshold rule together [\#25](https://github.com/gessnerfl/terraform-provider-instana/pull/25) ([gessnerfl](https://github.com/gessnerfl))

## [v0.5.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.5.0) (2019-10-15)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.4.0...v0.5.0)

**Implemented enhancements:**

- Support for label/name prefix and suffix [\#22](https://github.com/gessnerfl/terraform-provider-instana/issues/22)

**Merged pull requests:**

- \#22: migrate to customizable default name prefix and suffix [\#23](https://github.com/gessnerfl/terraform-provider-instana/pull/23) ([gessnerfl](https://github.com/gessnerfl))

## [v0.4.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.4.0) (2019-10-14)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.3.2...v0.4.0)

**Implemented enhancements:**

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

**Implemented enhancements:**

- Change release output [\#15](https://github.com/gessnerfl/terraform-provider-instana/issues/15)

**Merged pull requests:**

- Feature/15 release naming [\#16](https://github.com/gessnerfl/terraform-provider-instana/pull/16) ([gessnerfl](https://github.com/gessnerfl))

## [v0.2.2](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.2.2) (2019-05-07)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.2.1...v0.2.2)

**Implemented enhancements:**

- Add support for dashes in tags in match\_specification [\#12](https://github.com/gessnerfl/terraform-provider-instana/issues/12)

**Merged pull requests:**

- Fix application configuration example in README [\#14](https://github.com/gessnerfl/terraform-provider-instana/pull/14) ([steinex](https://github.com/steinex))

## [v0.2.1](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.2.1) (2019-05-07)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.2.0...v0.2.1)

**Merged pull requests:**

- Feature/12 add support for dashes in identifiers [\#13](https://github.com/gessnerfl/terraform-provider-instana/pull/13) ([gessnerfl](https://github.com/gessnerfl))

## [v0.2.0](https://github.com/gessnerfl/terraform-provider-instana/tree/v0.2.0) (2019-04-25)

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/v0.1.0...v0.2.0)

**Implemented enhancements:**

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

[Full Changelog](https://github.com/gessnerfl/terraform-provider-instana/compare/627e6874cfda8cf8e5d5793f016aaf60b5285e6f...v0.1.0)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
