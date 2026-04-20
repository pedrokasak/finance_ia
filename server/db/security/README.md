# Postgres Production Hardening (Neon/Supabase)

## 1. Do not expose the database publicly
- Use a private network path when possible (VPC / private endpoint).
- If public endpoint is required, allowlist only backend egress IPs.
- Disable shared admin credentials; create a dedicated app role.
- Enforce TLS (`sslmode=require`).

## 2. Minimum privileges for app role
Run as admin:

```sql
REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON DATABASE postgres FROM PUBLIC;

-- Example role for backend app:
-- CREATE ROLE app_backend LOGIN PASSWORD '...';
-- GRANT CONNECT ON DATABASE postgres TO app_backend;
-- GRANT USAGE ON SCHEMA public TO app_backend;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_backend;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA public
--   GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO app_backend;
```

## 3. Enable RLS and policies
Run:

```sql
\i db/security/rls.sql
```

or from repo root:

```bash
cd server
DATABASE_URL="postgres://..." ./scripts/apply_rls.sh
```

The backend must set the session user id on each request/transaction:

```sql
SELECT set_config('app.user_id', '<jwt-user-id>', true);
```

## 4. Atomic operations / race conditions
- Keep unique constraints for ownership and idempotency.
- Use `INSERT ... ON CONFLICT` for upserts (already used for budgets).
- For password reset, consume token atomically (`used=false AND expires_at > now()`).

## 5. Security checks in app layer
- JWT strict validation (signing method + exp + user_id format).
- Owner checks on `/user/:id` style routes.
- Idempotency key scoped by user + route + payload hash.
- Avatar upload validation:
  - only PNG/JPEG/WEBP
  - max 256KB
  - magic-bytes verification
  - reject external URLs (prevents image trackers).

## 6. Integration tests
- Run backend tests with:

```bash
go test ./...
```

- The CI pipeline includes an integration job with Postgres + Newman:
  - repository-level Postgres tests (`internal/infrastructure/database/finance`)
  - API collection tests via Newman (`server/postman/*.json`)
