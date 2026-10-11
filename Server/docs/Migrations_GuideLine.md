# BE FIT: Database Migration Guide

Step-by-step guide for creating and applying database migrations in the BE FIT backend.

- **Tool:** [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- **Database:** PostgreSQL 18 in Docker (container `BE-FIT`, database `be_fit`)
- **Shell:** PowerShell (VS Code terminal on Windows)
- **Migration folder:** `migrations/`
- **Baseline:** `000001_baseline` (the full schema). Older migrations live in `migrations/_archive/` for reference only.

---

## 1. Golden Rules

1. **Never edit a migration that has already been applied** anywhere other than your own throwaway DB. Write a new migration instead.
2. **One migration = one small, focused change.** Easy to review, easy to undo.
3. **Every migration has two files:** `.up.sql` (apply) and `.down.sql` (undo).
4. **Test on a scratch database first**, then apply to the real one.
5. **Never run `down` on production.** Fix forward with a new migration.
6. **Commit both files** (`.up.sql` and `.down.sql`) together with the code that needs them.

---

## 2. One-Time Setup

### 2.1 Check the tools

```powershell
migrate -version
docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Ports}}"
```

You should see `migrate` print a version and the `BE-FIT` container running on port 5432. If the container is stopped:

```powershell
docker start BE-FIT
```

### 2.2 Set helper variables (every time you open a new terminal)

```powershell
$env:DATABASE_URL = "postgres://tuann:<PASSWORD>@localhost:5432/be_fit?sslmode=disable"
$c = "BE-FIT"
```

Replace `<PASSWORD>` with your real password. Do not paste the real password into chats, screenshots, or commits.

### 2.3 Check the current state

```powershell
migrate -path migrations -database $env:DATABASE_URL version
```

Expected: a version number with **no** `dirty` flag (right now `1`).

---

## 3. Create a New Migration

### Step 1: Generate the file pair

```powershell
migrate create -ext sql -dir migrations -seq add_workout_plans
```

This creates:

```
migrations/000002_add_workout_plans.up.sql
migrations/000002_add_workout_plans.down.sql
```

Naming convention: `verb_object`, all lowercase with underscores, for example `add_avatar_to_users`, `create_exercises_table`, `drop_legacy_goal_column`.

### Step 2: Write the `up` SQL

Open `000002_add_workout_plans.up.sql` and write the change.

### Step 3: Write the `down` SQL

Open `000002_add_workout_plans.down.sql` and write the exact reverse of `up`.

---

## 4. SQL Recipes

### 4.1 Create a table

`up`:
```sql
CREATE TABLE workout_plans (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(150) NOT NULL,
    is_active  BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_workout_plans_user_id ON workout_plans(user_id);
```

`down`:
```sql
DROP TABLE IF EXISTS workout_plans;
```

Match the id type and style of your existing tables (check `000001_baseline.up.sql`) so foreign keys line up.

### 4.2 Add a column

`up`:
```sql
ALTER TABLE users ADD COLUMN avatar_url TEXT;
```

`down`:
```sql
ALTER TABLE users DROP COLUMN IF EXISTS avatar_url;
```

### 4.3 Add a NOT NULL column to a table that already has rows

Do it in safe stages, in one migration:

`up`:
```sql
ALTER TABLE users ADD COLUMN experience_level VARCHAR(20);
UPDATE users SET experience_level = 'beginner' WHERE experience_level IS NULL;
ALTER TABLE users ALTER COLUMN experience_level SET NOT NULL;
```

`down`:
```sql
ALTER TABLE users DROP COLUMN IF EXISTS experience_level;
```

### 4.4 Add an index

`up`:
```sql
CREATE INDEX idx_workout_sessions_user_id ON workout_sessions(user_id);
```

`down`:
```sql
DROP INDEX IF EXISTS idx_workout_sessions_user_id;
```

### 4.5 Rename a column

`up`:
```sql
ALTER TABLE users RENAME COLUMN full_name TO display_name;
```

`down`:
```sql
ALTER TABLE users RENAME COLUMN display_name TO full_name;
```

Renames break running code. Deploy the matching Go code change together with the migration.

### 4.6 Add or change an enum value

`up`:
```sql
ALTER TYPE goal_type ADD VALUE IF NOT EXISTS 'recomp';
```

`down`: Postgres cannot remove an enum value. Leave a comment explaining that, or use a lookup table instead of an enum for values that change often.

---

## 5. Test Before Applying to the Real DB

### Step 1: Create a scratch database

```powershell
docker exec $c psql -U tuann -d be_fit -c "CREATE DATABASE be_fit_verify;"
$verify = "postgres://tuann:<PASSWORD>@localhost:5432/be_fit_verify?sslmode=disable"
```

### Step 2: Run the full chain from empty (baseline plus all migrations)

```powershell
migrate -path migrations -database $verify up
```

### Step 3: Test that `down` and `up` both work

```powershell
migrate -path migrations -database $verify down 1
migrate -path migrations -database $verify up
```

All three commands must finish without errors.

### Step 4: Delete the scratch database

```powershell
docker exec $c psql -U tuann -d be_fit -c "DROP DATABASE be_fit_verify;"
```

If it complains about active connections, close DBeaver, pgAdmin, or any other tool connected to it.

---

## 6. Apply to Your Dev Database

```powershell
migrate -path migrations -database $env:DATABASE_URL up
migrate -path migrations -database $env:DATABASE_URL version
```

Then:

1. Run the Go app and confirm it works (`go run ./cmd/api`).
2. Run tests (`go test ./...`).
3. Commit.

```powershell
git add migrations/ 
git commit -m "db: add workout_plans table"
```

---

## 7. Rolling Back

### Undo the last migration (dev only)

```powershell
migrate -path migrations -database $env:DATABASE_URL down 1
```

### Jump to a specific version

```powershell
migrate -path migrations -database $env:DATABASE_URL goto 2
```

On shared or production databases, do not use `down`. Create a new migration that reverses the change instead.

---

## 8. Fixing a "dirty" Database

If a migration fails halfway, `migrate` marks the version as dirty and refuses to continue:

```
error: Dirty database version 2. Fix and force version.
```

Fix it:

1. Read the error message and find the SQL problem.
2. Inspect the database to see what was partially applied:
   ```powershell
   docker exec $c psql -U tuann -d be_fit -c "\dt"
   docker exec $c psql -U tuann -d be_fit -c "\d users"
   ```
3. Undo any partial changes by hand (Postgres runs most DDL in a transaction, so often nothing is left over).
4. Reset the version to the **last good one** (the version before the failed migration):
   ```powershell
   migrate -path migrations -database $env:DATABASE_URL force 1
   ```
5. Fix the `.up.sql` file and run `up` again.

`force` only rewrites the version number in `schema_migrations`. It does not run any SQL.

---

## 9. Setting Up a Fresh Database (new machine, CI, teammate)

1. Start Postgres:
   ```powershell
   docker compose up -d
   ```
2. Set `$env:DATABASE_URL` for the new empty database.
3. Run all migrations:
   ```powershell
   migrate -path migrations -database $env:DATABASE_URL up
   ```

The baseline plus every later migration builds the whole schema.

If your database already contains the baseline schema but no `schema_migrations` record (for example, restored from a dump), mark it as applied instead of running it:

```powershell
migrate -path migrations -database $env:DATABASE_URL force 1
```

---

## 10. Working With a Team

- After `git pull`, if new files appeared in `migrations/`, run `up`.
- If two branches both created `000002_...`, the second person to merge must rename theirs to `000003_...` (both files) before merging, then re-test.
- Review migrations in pull requests like code: check the `down`, check locks and data-loss risk.

---

## 11. Production Checklist

Before running migrations on staging or production:

- [ ] Backup taken: `pg_dump -Fc -U tuann -d be_fit -f backup.dump`
- [ ] Migration tested on a scratch DB with a copy of realistic data
- [ ] No destructive step (`DROP COLUMN`, `DROP TABLE`, type change) without a verified backup
- [ ] Large table changes checked for locking (adding an index with `CREATE INDEX CONCURRENTLY` must run outside a transaction)
- [ ] App code is compatible with both the old and new schema during deploy
- [ ] Run via CI/CD as a separate step, not on every app instance start
- [ ] Rollback plan = restore backup or a new forward-fix migration

---

## 12. Command Cheat Sheet

| Task | Command |
|---|---|
| Create migration | `migrate create -ext sql -dir migrations -seq <name>` |
| Apply all pending | `migrate -path migrations -database $env:DATABASE_URL up` |
| Apply next 1 | `migrate -path migrations -database $env:DATABASE_URL up 1` |
| Undo last 1 (dev only) | `migrate -path migrations -database $env:DATABASE_URL down 1` |
| Go to version N | `migrate -path migrations -database $env:DATABASE_URL goto N` |
| Show version | `migrate -path migrations -database $env:DATABASE_URL version` |
| Fix dirty state | `migrate -path migrations -database $env:DATABASE_URL force N` |
| List tables | `docker exec BE-FIT psql -U tuann -d be_fit -c "\dt"` |
| Describe a table | `docker exec BE-FIT psql -U tuann -d be_fit -c "\d table_name"` |

---

## 13. Troubleshooting

| Problem | Cause and fix |
|---|---|
| `Usage: migrate OPTIONS COMMAND` printed | `$env:DATABASE_URL` is empty. Set it again (section 2.2). PowerShell uses `$env:NAME`, not `$NAME`. |
| `No such container` | Wrong container name. Check with `docker ps`. Yours is `BE-FIT`. |
| `Dirty database version N` | A migration failed halfway. See section 8. |
| `no change` | Nothing pending. Normal. |
| `file does not exist` or `no migration found for version N` | The database version is higher than your newest file, or you pulled/moved files. Check `version` against the files in `migrations/`. |
| `relation already exists` | The change was already applied (maybe by hand). Remove the duplicate statement or use `IF NOT EXISTS`, then `force` the right version. |
| `password authentication failed` | Wrong user or password in `DATABASE_URL`. |
| Syntax error near `\restrict` | `pg_dump` 18 adds `\restrict` lines. Delete lines starting with `\restrict` or `\unrestrict` from any dumped SQL. |
| Garbled characters in an SQL file | File saved as UTF-16 or with BOM by PowerShell redirection (`>`). Save as UTF-8 (no BOM) in VS Code, and avoid `>` for SQL files. |

---

## 14. Quick Workflow (TL;DR)

```powershell
# 1. set env
$env:DATABASE_URL = "postgres://tuann:<PASSWORD>@localhost:5432/be_fit?sslmode=disable"

# 2. create files
migrate create -ext sql -dir migrations -seq <what_you_are_changing>

# 3. write up.sql and down.sql in VS Code

# 4. apply
migrate -path migrations -database $env:DATABASE_URL up

# 5. verify
migrate -path migrations -database $env:DATABASE_URL version
go test ./...

# 6. commit both files
git add migrations/
git commit -m "db: <what_you_changed>"
```