# Table: torii_organization

Get the organization profile for your Torii account. Returns a single row.

## Example queries

**Get organization details:**

```sql
select
  id,
  company_name,
  domain,
  creation_time
from
  torii_organization;
```

**Get the organization domain:**

```sql
select
  domain
from
  torii_organization;
```
