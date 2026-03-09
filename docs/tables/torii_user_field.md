# Table: torii_user_field

List custom user fields from connected integrations in your Torii organization.

## Example queries

**List all user fields:**

```sql
select
  id,
  name,
  key,
  type,
  source_id_app
from
  torii_user_field;
```

**List active (non-deleted) user fields:**

```sql
select
  id,
  name,
  key,
  type,
  source_id_app
from
  torii_user_field
where
  is_deleted = false;
```

**Search for a field by name:**

```sql
select
  id,
  name,
  key,
  type,
  source_id_app
from
  torii_user_field
where
  name = 'Department';
```
