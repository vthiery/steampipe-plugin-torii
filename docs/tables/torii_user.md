# Table: torii_user

List users in your Torii organization.

## Example queries

**List all active users:**

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

**List all external users:**

```sql
select
  id,
  email,
  lifecycle_status
from
  torii_user
where
  is_external = true;
```

**Find a user by email:**

```sql
select
  id,
  first_name,
  last_name,
  role,
  active_apps_count
from
  torii_user
where
  email = 'tony@stark.com';
```

**List users currently offboarding:**

```sql
select
  id,
  email,
  active_apps_count
from
  torii_user
where
  lifecycle_status = 'offboarding';
```
