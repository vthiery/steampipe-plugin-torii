# Table: torii_contract_field

List custom and built-in fields defined for contracts in your Torii organization.

## Example queries

**List all contract fields:**

```sql
select
  id,
  name,
  system_key,
  type
from
  torii_contract_field;
```

**List all dropdown fields with their options:**

```sql
select
  id,
  name,
  system_key,
  options
from
  torii_contract_field
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
  torii_contract_field
where
  name ilike '%status%';
```
