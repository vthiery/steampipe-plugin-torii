# Table: torii_contract

List software contracts managed in your Torii organization.

## Example queries

**List all active contracts:**

```sql
select
  id,
  name,
  id_app,
  owner
from
  torii_contract
where
  status = 'active';
```

**List all contracts:**

```sql
select
  id,
  name,
  id_app,
  owner,
  status
from
  torii_contract
order by
  name;
```

**Find contracts for a specific application:**

```sql
select
  id,
  name,
  owner,
  status
from
  torii_contract
where
  id_app = 1234;
```
