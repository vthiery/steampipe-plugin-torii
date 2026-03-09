# Table: torii_app_field

List custom and built-in fields defined for applications in your Torii organization.

## Example queries

**List all application fields:**

```sql
select
  id,
  name,
  system_key,
  type
from
  torii_app_field;
```

**List all dropdown fields with their options:**

```sql
select
  id,
  name,
  system_key,
  options
from
  torii_app_field
where
  type = 'dropdown';
```

**Find a field by name:**

```sql
select
  id,
  name,
  system_key,
  type,
  options
from
  torii_app_field
where
  name = 'Plan Name';
```
