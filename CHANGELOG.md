## v0.0.3 [2026-03-13]

_Bug fixes_

- `torii_audit_log`: fixed the `entity` key column — it was incorrectly mapped to `type`, is now required (matching the Torii API requirement), and is populated on each returned row.

_Dependencies_

- Bumped `github.com/hashicorp/go-getter` from 1.7.5 to 1.7.9
- Bumped `github.com/ulikunitz/xz` from 0.5.10 to 0.5.14
- Bumped `github.com/golang/glog` from 1.2.0 to 1.2.4
- Bumped `go.opentelemetry.io/otel/sdk` from 1.26.0 to 1.40.0
- Bumped `golang.org/x/crypto` from 0.21.0 to 0.45.0
- Bumped `golang.org/x/net` from 0.23.0 to 0.38.0
- Bumped `golang.org/x/oauth2` from 0.17.0 to 0.27.0

## v0.0.2 [2026-03-09]

_What's new?_

- New tables added:
  - [torii_app_field](docs/tables/torii_app_field.md) — list custom and built-in fields defined for applications
  - [torii_contract_field](docs/tables/torii_contract_field.md) — list custom and built-in fields defined for contracts
  - [torii_organization](docs/tables/torii_organization.md) — get the organization profile
  - [torii_user_app](docs/tables/torii_user_app.md) — list applications used by a specific user
  - [torii_user_field](docs/tables/torii_user_field.md) — list custom user fields from connected integrations

## v0.0.1 [2026-03-09]

_What's new?_

- Initial release of the plugin with 6 tables:
  - [torii_app](docs/tables/torii_app.md)
  - [torii_app_user](docs/tables/torii_app_user.md)
  - [torii_audit_log](docs/tables/torii_audit_log.md)
  - [torii_contract](docs/tables/torii_contract.md)
  - [torii_role](docs/tables/torii_role.md)
  - [torii_user](docs/tables/torii_user.md)
