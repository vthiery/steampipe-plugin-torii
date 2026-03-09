# Table: torii_audit_log

List admin audit log events in your Torii organization.

## Example queries

**List the most recent audit log entries:**

```sql
select
  performed_by_email,
  type,
  creation_time
from
  torii_audit_log
order by
  creation_time desc
limit 50;
```

**Find all workflow-related audit events:**

```sql
select
  performed_by_email,
  type,
  creation_time,
  properties
from
  torii_audit_log
where
  type like '%workflow%'
order by
  creation_time desc;
```

**Find user lifecycle changes:**

```sql
select
  performed_by_email,
  type,
  creation_time,
  properties
from
  torii_audit_log
where
  type in ('update_user_lifecycle_status', 'update_user')
order by
  creation_time desc;
```

**Find contract-related audit events:**

```sql
select
  performed_by_email,
  type,
  creation_time,
  properties
from
  torii_audit_log
where
  type like '%contract%'
order by
  creation_time desc;
```
