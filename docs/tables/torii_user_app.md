# Table: torii_user_app

List applications used by a specific user in your Torii organization. Requires `user_id` as a filter.

## Example queries

**List all applications for a user:**

```sql
select
  id,
  name,
  state,
  category,
  is_user_removed_from_app
from
  torii_user_app
where
  user_id = 12345;
```

**List active (non-removed) applications for a user:**

```sql
select
  id,
  name,
  category,
  last_usage_time,
  annual_cost,
  currency
from
  torii_user_app
where
  user_id = 12345
  and is_user_removed_from_app = false;
```

**Find managed applications for a user:**

```sql
select
  id,
  name,
  category,
  annual_cost,
  currency
from
  torii_user_app
where
  user_id = 12345
  and state = 'Managed';
```

**Count applications per user (join with torii_user):**

```sql
select
  u.email,
  count(a.id) as app_count
from
  torii_user as u
  join torii_user_app as a on a.user_id = u.id
where
  u.lifecycle_status = 'active'
  and a.is_user_removed_from_app = false
group by
  u.email
order by
  app_count desc;
```
