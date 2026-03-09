# Torii Plugin for Steampipe

Use SQL to query users, applications, contracts, roles, audit logs, and more from [Torii](https://www.toriihq.com).

## Installation

Clone and build the plugin:

```sh
git clone https://github.com/vthiery/steampipe-plugin-torii.git
cd steampipe-plugin-torii
mkdir -p ~/.steampipe/plugins/local/torii
go build -o ~/.steampipe/plugins/local/torii/steampipe-plugin-torii.plugin .
```

## Configuration

Copy the sample config:

```sh
cp config/torii.spc ~/.steampipe/config/torii.spc
```

Edit `~/.steampipe/config/torii.spc`:

```hcl
connection "torii" {
  plugin  = "local/torii"

  # API key from the Torii dashboard → Settings → Security → API Keys.
  # Required for all table queries.
  api_key = "YOUR_API_KEY"
}
```

## Tables

| Table | Description |
|-------|-------------|
| [torii_user](tables/torii_user.md) | List users in your Torii organization. |
| [torii_app](tables/torii_app.md) | List applications discovered or managed in your organization. |
| [torii_app_user](tables/torii_app_user.md) | List users of a specific application. |
| [torii_contract](tables/torii_contract.md) | List software contracts managed in your organization. |
| [torii_role](tables/torii_role.md) | List roles defined in your organization. |
| [torii_audit_log](tables/torii_audit_log.md) | List admin audit log events. |
