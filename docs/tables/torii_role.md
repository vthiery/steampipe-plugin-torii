# Table: torii_role

List roles defined in your Torii organization.

## Example queries

**List all roles:**

```sql
select
  id,
  name,
  system_key,
  is_admin,
  users_count
from
  torii_role;
```

**List admin roles:**

```sql
select
  id,
  name,
  description,
  users_count
from
  torii_role
where
  is_admin = true;
```

**Find roles with no users:**

```sql
select
  id,
  name,
  description
from
  torii_role
where
  users_count = 0;
```
