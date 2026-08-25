---
name: sql-review
description: Use when SQL is written, changed, or reviewed against the database standards — triggers on "/sql-review", "review this SQL", "check this migration", "audit this schema", "does this follow our conventions", when the user asks to write or alter a table, index, view, query, or migration, or when a diff touches sqlc.yaml or sqlc-generated code.
---

# SQL Review

Hold all SQL — new or existing — to the standards below. One checklist drives
both modes: reviewing flags deviations, writing avoids them from the start.

## Modes

**Review mode** — the user hands over existing SQL.

1. If a file path, diff, or migration was named, read it. Otherwise ask the
   user to paste the SQL.
2. Work through every checklist section below.
3. Report findings in the output format at the end.
4. If the SQL is clean, say so explicitly. Do not invent findings.

**Write mode** — the user asks for new or changed SQL.

1. Write the SQL with the checklist as constraints.
2. Self-review the result against every section.
3. Show the SQL. Name any standard you bent and why.

---

## Checklist

### Naming conventions

**Tables**
- [ ] `snake_case` only
- [ ] Plural nouns (`users`, `lab_orders`, not `user` or `labOrder`)
- [ ] No version suffixes (`v2`, `new_`, `_old`)
- [ ] Max 50 characters
- [ ] Business terminology, not source-system terminology

**Columns**
- [ ] `snake_case` only
- [ ] Primary key **must** be named `{table_singular}_id` — `id` is never
      acceptable (e.g. `organizations.organization_id`, `orders.order_id`,
      `ticket_tiers.ticket_tier_id`). Every table using `id` is a HIGH finding.
- [ ] Foreign keys named `{referenced_table_singular}_id`
- [ ] Timestamps use `{event}_at` past-tense format (`created_at`,
      `updated_at`, `deleted_at`)
- [ ] Dates use `{event}_date` format (`order_date`, not `date_order`)
- [ ] Booleans prefixed with `is_` or `has_` (`is_active`, `has_results`)
- [ ] Price/currency fields: decimal format (`19.99`) unless using cents
      suffix (`price_in_cents`)
- [ ] No abbreviations (`uid` → `user_id`, `usr_demo` → `user_demographics`)
- [ ] No reserved words as column names
- [ ] Max 50 characters

**Indexes**
- [ ] Regular indexes: `idx_{table}_{columns}`
- [ ] Unique indexes: `uidx_{table}_{columns}`
- [ ] Max 63 characters total

**Views / Materialized views**
- [ ] Views prefixed `v_`
- [ ] Materialized views prefixed `mv_`

---

### Column ordering

Columns appear in this order within a `CREATE TABLE`:

1. IDs (primary key first, then foreign keys)
2. Strings (text, enums, varchars)
3. Numerics (amounts, counts, measurements)
4. Booleans
5. Dates (date-only fields)
6. Timestamps (datetime fields)

---

### Data types

- [ ] All timestamps use `TIMESTAMPTZ`, never bare `TIMESTAMP`
- [ ] Primary keys are UUID (`gen_random_uuid()`), or SERIAL for internal
      lookup tables
- [ ] Categorical / constrained values use a custom `ENUM` type or a `CHECK`
      constraint — not bare `VARCHAR` with no validation
- [ ] ZIP codes: `VARCHAR(5)` with `CHECK (zip_code ~ '^[0-9]{5}$')`
- [ ] States: `us_state_type` enum
- [ ] Gender: `gender_type` enum
- [ ] Flexible structured data: `JSONB`, not `JSON` or `TEXT`

---

### Constraints and integrity

- [ ] Every table has a primary key
- [ ] Foreign keys are explicitly declared (or noted as analytics schema
      where omitted intentionally)
- [ ] NOT NULL applied to columns that must always have a value
- [ ] Default values provided where sensible (`DEFAULT NOW()`,
      `DEFAULT false`)
- [ ] JSON columns have a `CHECK (jsonb_typeof(...) = 'object')` constraint
      where structure matters

---

### Audit fields

Important business entities must include:

- [ ] `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- [ ] `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- [ ] `created_by` / `updated_by` where actor tracking matters
- [ ] `deleted_at TIMESTAMPTZ` (soft delete) instead of hard deletes on
      important entities

---

### Index coverage

- [ ] Every foreign key column has an index
- [ ] Time-series columns used in WHERE clauses are indexed (`created_at`,
      `order_date`)
- [ ] Composite indexes order columns to match query predicates — equality
      columns first, then range columns
- [ ] Partial indexes used where appropriate (e.g.
      `WHERE deleted_at IS NULL`)
- [ ] `CREATE INDEX CONCURRENTLY` used in production migrations to avoid
      locks — note it cannot run inside a transaction, so the migration tool
      must not wrap it in one
- [ ] Covering indexes (`INCLUDE`) used for high-frequency read patterns

---

### Queries

- [ ] No `SELECT *` in application queries or view definitions
- [ ] WHERE predicates are sargable — no function or cast wrapped around the
      indexed column
- [ ] Every `UPDATE` / `DELETE` has a WHERE clause, or a comment stating the
      full-table intent
- [ ] Explicit `JOIN ... ON`, never comma joins
- [ ] If the table uses soft deletes, queries filter
      `deleted_at IS NULL`

---

### Schema separation

- [ ] Operational (transactional) tables live in `public`
- [ ] Analytics / reporting tables live in `analytics` schema
- [ ] Analytics fact tables are denormalized with pre-calculated measures
      and time dimensions
- [ ] ETL metadata columns present on analytics tables
      (`source_updated_at`, `etl_batch_id`)

---

### Migrations

- [ ] Only safe operations without explicit review: `ADD COLUMN` (nullable),
      `CREATE INDEX CONCURRENTLY`, `ADD CONSTRAINT ... NOT VALID` + separate
      `VALIDATE CONSTRAINT`
- [ ] Column renames done in 3 steps: add → backfill → drop, each in a
      separate migration
- [ ] Rollback instructions included as a comment in the migration
- [ ] Migration file follows `NNN_descriptive_name.sql` naming

---

### sqlc configuration

If the project uses sqlc, compare `sqlc.yaml` against the canonical config
below. Flag every deviation as a finding.

- [ ] `version: "2"` — never version 1
- [ ] `engine: "postgresql"`
- [ ] Queries live in `database/queries.sql`; schema comes from
      `database/migrations`
- [ ] `sql_package: "pgx/v5"` — not `database/sql`, not pgx/v4
- [ ] Generated files carry the `.gen.go` suffix (`db.gen.go`,
      `models.gen.go`, `output_files_suffix: ".gen.go"`)
- [ ] `emit_db_tags`, `emit_sql_as_comment`, `emit_pointers_for_null_types`,
      and `emit_empty_slices` are all `true`
- [ ] `uuid` overrides to `github.com/google/uuid` `UUID`, pointer when
      nullable
- [ ] `pg_catalog.timestamp` and `pg_catalog.timestamptz` override to
      `time.Time`, pointer when nullable

If the config defines multiple databases, share these settings via YAML
anchors (`&default`, `&go_defaults`). A single-database config needs no
anchors — their absence is not a finding.

Canonical config:

```yaml
version: "2"
sql:
  - &default
    engine: "postgresql"
    queries:
      - "database/queries.sql"
    schema: "database/migrations"
    gen:
      go: &go_defaults
        package: "database"
        out: "database"
        sql_package: "pgx/v5"
        output_db_file_name: "db.gen.go"
        output_models_file_name: "models.gen.go"
        output_files_suffix: ".gen.go"
        emit_db_tags: true
        emit_sql_as_comment: true
        emit_pointers_for_null_types: true
        emit_empty_slices: true
        overrides:
          - db_type: "pg_catalog.timestamp"
            go_type:
              type: "time.Time"
          - db_type: "pg_catalog.timestamp"
            nullable: true
            go_type:
              type: "*time.Time"
          - db_type: "pg_catalog.timestamptz"
            go_type:
              type: "time.Time"
          - db_type: "pg_catalog.timestamptz"
            nullable: true
            go_type:
              type: "*time.Time"
          - db_type: "uuid"
            go_type:
              import: "github.com/google/uuid"
              type: "UUID"
          - db_type: "uuid"
            nullable: true
            go_type:
              type: "*uuid.UUID"
```

---

### Security

- [ ] Row Level Security (RLS) enabled on multi-tenant tables
- [ ] Policies created for `application_role` and `analytics_role` separately
- [ ] `REVOKE CREATE ON SCHEMA public FROM PUBLIC` present in initial setup
- [ ] PII columns identified in comments

---

## Output format

```
## SQL Review

### ✅ Passing
- <brief list of things that look correct>

### ⚠️ Findings

**[Category — CRITICAL | HIGH | MEDIUM | LOW]**
- Finding: <what's wrong>
  Snippet: `<offending SQL>`
  Fix: `<corrected SQL>`

### Summary
<one-line verdict and priority of fixes>
```

Severity guide:

- **CRITICAL**: data loss, corruption, or security exposure — no primary
  key, missing WHERE on UPDATE/DELETE, destructive migration, RLS absent on
  a multi-tenant table
- **HIGH**: wrong types or broken conventions — bare `TIMESTAMP`, `JSON`
  instead of `JSONB`, PK named `id`, missing FK declaration, unsafe
  migration operation
- **MEDIUM**: friction ahead — missing indexes, missing constraints or audit
  fields, bad column order
- **LOW**: style — abbreviations, missing comments

**If a rule in this checklist is broken, it is a finding at its listed
severity — never downgrade it to "advisory", "optional", or "if preferred".**

If there are no findings, say "No findings — SQL follows all standards."
