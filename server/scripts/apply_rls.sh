#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RLS_SQL="${ROOT_DIR}/db/security/rls.sql"

if ! command -v psql >/dev/null 2>&1; then
  echo "psql not found. Install PostgreSQL client tools first."
  exit 1
fi

DB_URL="${DATABASE_URL:-${TEST_DATABASE_URL:-}}"
if [[ -z "${DB_URL}" ]]; then
  echo "Set DATABASE_URL (or TEST_DATABASE_URL) before running."
  exit 1
fi

if [[ ! -f "${RLS_SQL}" ]]; then
  echo "RLS script not found at ${RLS_SQL}"
  exit 1
fi

echo "Applying RLS script to target database..."
psql "${DB_URL}" -v ON_ERROR_STOP=1 -f "${RLS_SQL}"
echo "RLS applied successfully."

