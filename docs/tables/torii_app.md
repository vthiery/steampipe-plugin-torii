# Table: torii_app

List applications discovered or managed in your Torii organization.

## Example queries

**List all managed applications:**

```sql
select
  id,
  name,
  category,
  active_users_count
from
  torii_app
where
  state = 'Managed';
```

**List all discovered (unreviewed) applications:**

```sql
select
  id,
  name,
  category,
  active_users_count,
  creation_time
from
  torii_app
where
  state = 'Discovered'
order by
  active_users_count desc;
```

**Find applications by category:**

```sql
select
  id,
  name,
  state,
  active_users_count
from
  torii_app
where
  category = 'Developer Tools';
```

**List custom (manually added) applications:**

```sql
select
  id,
  name,
  state,
  creation_time
from
  torii_app
where
  is_custom = true;
```
