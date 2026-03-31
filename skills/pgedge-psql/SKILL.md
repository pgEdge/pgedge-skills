---
name: pgedge-psql
description: "Use this skill whenever the user wants to query, explore, or modify a PostgreSQL database using the psql command-line tool. This includes: running SQL queries against a Postgres database, exploring database schemas (listing tables, describing columns, checking indexes), inserting or updating data, running DDL statements, or any task where the user mentions psql, PostgreSQL, Postgres, or a database connection string starting with postgresql://. Trigger whenever the user mentions a database table, asks to 'check the database', 'run a query', 'look at the data', or references a Postgres connection — even casually. Do NOT trigger if the user is working with the pgedge-postgres MCP server (which has its own guardrails and knowledgebase), or if they're writing application code that uses a Postgres driver/ORM rather than the psql CLI directly."
---

# PostgreSQL psql Skill

This skill guides you through using the `psql` command-line tool to interact with PostgreSQL databases efficiently. The core principles are: minimize token usage by using compact output formats, default to read-only access unless the user explicitly needs writes, and always discover the schema before writing queries against unfamiliar tables.

## Connecting to PostgreSQL

psql supports several connection methods. Ask the user how they connect if it's not clear from context.

### Connection string (URI format)
```bash
psql "postgresql://user:password@host:5432/dbname"
```
This is the most portable format. If the user provides a connection string, use it directly.

### Individual flags
```bash
psql -h hostname -p 5432 -U username -d dbname
```

### Environment variables
The user may have these set already — check before asking for connection details:
```bash
echo $PGHOST $PGPORT $PGUSER $PGDATABASE
```
The standard libpq environment variables are: `PGHOST`, `PGPORT`, `PGUSER`, `PGDATABASE`, `PGPASSWORD` (though `.pgpass` is preferred over PGPASSWORD for security).

### .pgpass file
If the user has a `~/.pgpass` file, psql will authenticate automatically without prompting. The file format is:
```
hostname:port:database:username:password
```
You don't need to do anything special — psql reads it automatically. Just note that if authentication fails silently, checking for a `.pgpass` file (and its permissions — it must be `chmod 600`) is a good debugging step.

### Connection verification
After connecting, verify the connection works before doing anything else:
```bash
psql <connection> -c "SELECT version();"
```

## Token-Efficient Output Formats

This is critical. psql's default "aligned" output is human-readable but extremely wasteful for LLM consumption — it pads columns with spaces, draws borders, and adds headers/footers that consume tokens without adding information value.

### The default to use: `--csv -t`

For almost all queries, use `--csv` combined with `-t` (tuples-only):

```bash
psql <connection> --csv -t -c "SELECT id, name, email FROM users LIMIT 5;"
```

This produces compact CSV output with no headers and no row-count footer. It's the most token-efficient format that's still unambiguous (handles commas and quotes in data correctly via RFC 4180 quoting rules).

### When you need column headers: `--csv` (without `-t`)

If the user is asking about data and you need to see column names to understand what you're looking at (common during schema discovery), drop the `-t`:

```bash
psql <connection> --csv -c "SELECT * FROM users LIMIT 3;"
```

This adds a single header row — a small token cost that's worth it when you don't know the columns yet.

### When to use `-A -t` instead

The unaligned format (`-A -t`) is slightly more compact than CSV for simple data (no quoting overhead), but it doesn't handle embedded delimiters safely. Use it only when you're confident the data contains no pipe characters:

```bash
psql <connection> -A -t -c "SELECT count(*) FROM users;"
```

This is ideal for single-value queries (counts, existence checks, scalar lookups).

### Format summary

| Scenario | Flags | Why |
|---|---|---|
| Most queries | `--csv -t` | Compact, unambiguous, no wasted tokens |
| Need column names | `--csv` | One header row, still compact |
| Single scalar value | `-A -t` | Absolute minimum output |
| User wants readable output | (none) | Human-friendly, use only when displaying to user |

### Limiting output

Always use `LIMIT` when exploring data you haven't seen before. Start with `LIMIT 5` or `LIMIT 10` to understand the shape of the data before pulling more. Large result sets waste tokens and may not even fit in context.

```bash
psql <connection> --csv -c "SELECT * FROM large_table LIMIT 5;"
```

## Schema Discovery

Before writing queries against unfamiliar tables, discover the schema first. This prevents wasted round-trips from wrong column names or misunderstood table structures.

### Step 1: List schemas and tables

```bash
# List all user-created schemas (excludes pg_catalog, information_schema)
psql <connection> --csv -t -c "
  SELECT schemaname, tablename
  FROM pg_tables
  WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
  ORDER BY schemaname, tablename;
"
```

For a quick overview, the `\dt` meta-command also works, but its output is in the aligned format and harder to parse programmatically. The SQL query above is more token-efficient.

### Step 2: Describe a specific table

```bash
# Get column names, types, and nullability
psql <connection> --csv -c "
  SELECT column_name, data_type, is_nullable, column_default
  FROM information_schema.columns
  WHERE table_schema = 'public' AND table_name = 'users'
  ORDER BY ordinal_position;
"
```

This is preferable to `\d tablename` because the output is structured and compact.

### Step 3: Check indexes (when performance matters)

```bash
psql <connection> --csv -c "
  SELECT indexname, indexdef
  FROM pg_indexes
  WHERE tablename = 'users';
"
```

### Step 4: Check row count (approximate, fast)

```bash
# Fast approximate count from pg_stat — avoids full table scan
psql <connection> -A -t -c "
  SELECT reltuples::bigint FROM pg_class WHERE relname = 'users';
"
```

For an exact count, use `SELECT count(*) FROM users;` but be aware this performs a full table scan on large tables.

### Step 5: Sample a few rows

```bash
psql <connection> --csv -c "SELECT * FROM users LIMIT 5;"
```

This gives you real data to understand formats, typical values, and relationships.

## Read-Only Access Pattern

When the user's task is purely about reading data (querying, exploring, reporting), use a read-only transaction to prevent accidental modifications. This is a safety net — it costs nothing and protects the database.

### For single queries

```bash
psql <connection> --csv -t -c "
  BEGIN TRANSACTION READ ONLY;
  SELECT * FROM orders WHERE status = 'pending' LIMIT 10;
  COMMIT;
"
```

### For multi-query sessions via a SQL file

Write a `.sql` file with `BEGIN TRANSACTION READ ONLY;` at the top:

```bash
cat > /tmp/query.sql << 'EOF'
BEGIN TRANSACTION READ ONLY;

SELECT count(*) FROM orders;
SELECT status, count(*) FROM orders GROUP BY status;

COMMIT;
EOF

psql <connection> --csv -f /tmp/query.sql
```

### When NOT to use read-only

Skip the read-only wrapper when the user explicitly asks to insert, update, delete, or modify schema. In those cases, follow the Write Operations section below.

## Write Operations

When the user explicitly asks to modify data or schema, proceed carefully. The key principle: understand before you change.

### Pre-write checklist

1. **Confirm the target** — verify the table name, schema, and what's currently in it
2. **Show the plan** — tell the user what you're about to do and ask for confirmation before running destructive operations (UPDATE, DELETE, DROP, TRUNCATE, ALTER)
3. **Use transactions** — wrap writes in explicit transactions so they can be rolled back if something looks wrong

### Safe write pattern

```bash
psql <connection> -c "
  BEGIN;

  -- Show what will be affected first
  SELECT id, status FROM orders WHERE status = 'stale';

  -- Then make the change
  UPDATE orders SET status = 'archived' WHERE status = 'stale';

  COMMIT;
"
```

For DDL operations (CREATE TABLE, ALTER TABLE, DROP TABLE), always confirm with the user first, and prefer `IF EXISTS` / `IF NOT EXISTS` clauses to avoid errors on re-runs:

```bash
psql <connection> -c "CREATE TABLE IF NOT EXISTS audit_log (
  id serial PRIMARY KEY,
  action text NOT NULL,
  created_at timestamptz DEFAULT now()
);"
```

### Destructive operations — extra caution

For DELETE, DROP, or TRUNCATE, always:
1. First run a SELECT to show what will be affected
2. Show the user the count and a sample of rows
3. Get explicit confirmation before proceeding
4. Use a transaction so it can be rolled back

## Calling psql Efficiently

### One-shot queries with `-c`

For single queries, use `-c` to avoid an interactive session:

```bash
psql <connection> --csv -t -c "SELECT count(*) FROM users;"
```

### Multi-statement scripts with `-f`

For multiple related queries, write them to a file and use `-f`. This is a single psql invocation — much more efficient than multiple `-c` calls:

```bash
cat > /tmp/analysis.sql << 'EOF'
\timing on

SELECT count(*) AS total_users FROM users;
SELECT status, count(*) FROM orders GROUP BY status ORDER BY count DESC;
SELECT avg(total) FROM orders WHERE created_at > now() - interval '30 days';
EOF

psql <connection> --csv -f /tmp/analysis.sql
```

### Avoid multiple separate invocations

Bad (3 connection round-trips):
```bash
psql <connection> --csv -t -c "SELECT count(*) FROM users;"
psql <connection> --csv -t -c "SELECT count(*) FROM orders;"
psql <connection> --csv -t -c "SELECT count(*) FROM products;"
```

Good (1 connection, 1 invocation):
```bash
psql <connection> --csv -t -c "
  SELECT 'users', count(*) FROM users
  UNION ALL
  SELECT 'orders', count(*) FROM orders
  UNION ALL
  SELECT 'products', count(*) FROM products;
"
```

Or use a script file for more complex multi-query work.

### Using psql variables

For parameterized queries, use psql's `-v` flag to pass variables:

```bash
psql <connection> --csv -t -v target_status='pending' -c "
  SELECT * FROM orders WHERE status = :'target_status' LIMIT 10;
"
```

Note the `:'varname'` syntax — the colon-quote form safely quotes the value as a string literal.

## Error Handling

### Connection failures
If psql can't connect, check:
1. Is the host reachable? (`pg_isready -h hostname -p 5432`)
2. Are credentials correct? (check env vars or .pgpass)
3. Is the database name right? (`psql <connection> -l` lists available databases)
4. Is there a firewall or pg_hba.conf issue? (the error message usually says)

### Query errors
When a query fails, psql prints the error to stderr. Read the error message carefully — PostgreSQL error messages are usually quite specific about what went wrong and where.

Common mistakes:
- Wrong column name → run schema discovery again
- Permission denied → the user may need a different role, or the table may be in a different schema
- Syntax error → check for missing commas, mismatched quotes, or PostgreSQL-specific syntax differences

## Workflow Summary

When a user asks you to work with a PostgreSQL database:

1. **Establish connection** — figure out how to connect (ask if needed, check env vars)
2. **Verify connection** — run a quick `SELECT version();` or `SELECT 1;`
3. **Discover schema** — unless the user has already told you the table structure, list tables and describe the relevant ones
4. **Use compact output** — default to `--csv -t` for queries, `--csv` when you need headers, `-A -t` for scalar values
5. **Limit results** — use LIMIT when exploring, especially on unfamiliar tables
6. **Read-only by default** — wrap read queries in `BEGIN TRANSACTION READ ONLY` unless the user needs writes
7. **Confirm before writing** — for any data modification, show what will change and get user confirmation
8. **Batch queries** — use `-f` with a script file or combine queries into one `-c` invocation rather than running psql multiple times
