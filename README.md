# Torii Plugin for Steampipe

Use SQL to query users, applications, contracts, roles, audit logs, and more from [Torii](https://www.toriihq.com) SaaS management platform.

- **[Get started →](https://github.com/vthiery/steampipe-plugin-torii)**
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
- Get involved: [Issues](https://github.com/vthiery/steampipe-plugin-torii/issues)

## Quick start

### Prerequisites

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

### Build and Install

```sh
git clone https://github.com/vthiery/steampipe-plugin-torii.git
cd steampipe-plugin-torii
make install
```

### Configuration

Copy the sample config:

```sh
cp config/torii.spc ~/.steampipe/config/torii.spc
```

Edit `~/.steampipe/config/torii.spc`:

```hcl
connection "torii" {
  plugin  = "local/torii"

  # API key from Torii dashboard → Settings → Security → API Keys.
  api_key = "YOUR_API_KEY"
}
```

### Run a query

```shell
steampipe query
```

List all active users:

```sql
select
  id,
  first_name,
  last_name,
  email,
  role
from
  torii_user
where
  lifecycle_status = 'active';
```

## Tables

| Table | Description |
|-------|-------------|
| [torii_user](torii/table_torii_user.go) | List users in your Torii organization. |
| [torii_app](torii/table_torii_app.go) | List applications discovered or managed in your Torii organization. |
| [torii_app_user](torii/table_torii_app_user.go) | List users of a specific application. |
| [torii_contract](torii/table_torii_contract.go) | List software contracts managed in your Torii organization. |
| [torii_role](torii/table_torii_role.go) | List roles defined in your Torii organization. |
| [torii_audit_log](torii/table_torii_audit_log.go) | List admin audit log events. |

## Development

### Prerequisites

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

### Build and Install

```sh
make install
```

Configure the plugin:

```sh
cp config/torii.spc ~/.steampipe/config/torii.spc
vi ~/.steampipe/config/torii.spc
```

## Testing

Run a smoke query against every table:

```sh
make test
```

The test script ([scripts/test_tables.sh](scripts/test_tables.sh)) builds the plugin, queries each table, and reports pass/fail/skip.

### Further reading

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Torii API reference](https://developers.toriihq.com/reference/introduction-1)
