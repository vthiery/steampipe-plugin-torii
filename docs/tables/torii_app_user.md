# Table: torii_app_user

List users of a specific application in your Torii organization. Requires `app_id` as a filter.

## Example queries

**List all active users of an application:**

```sql
select
  id_user,
  email,
  full_name,
  role,
  last_visit_time
from
  torii_app_user
where
  app_id = 1234
  and status = 'active';
```

**Find users who have been removed from an application:**

```sql
select
  id_user,
  email,
  full_name,
  is_user_removed_from_app
from
  torii_app_user
where
  app_id = 1234
  and is_user_removed_from_app = true;
```

**List external users of an application:**

```sql
select
  id_user,
  email,
  full_name,
  status
from
  torii_app_user
where
  app_id = 1234
  and is_external = true;
```

**Find users pending offboarding from an application:**

```sql
select
  id_user,
  email,
  lifecycle_status,
  is_owner_notified_to_remove_user
from
  torii_app_user
where
  app_id = 1234
  and lifecycle_status = 'offboarding';
```
