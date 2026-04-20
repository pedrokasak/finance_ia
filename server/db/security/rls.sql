-- RLS hardening for Postgres (Neon/Supabase)
-- Run as database owner / admin.
-- Backend must set app.user_id per request:
--   SELECT set_config('app.user_id', '<uuid>', true);

BEGIN;

-- Ensure owner columns are indexed
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets(user_id);
CREATE INDEX IF NOT EXISTS idx_goals_user_id ON goals(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_insights_user_id ON ai_insights(user_id);

-- Enable and enforce RLS
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE users FORCE ROW LEVEL SECURITY;

ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
ALTER TABLE transactions FORCE ROW LEVEL SECURITY;

ALTER TABLE categories ENABLE ROW LEVEL SECURITY;
ALTER TABLE categories FORCE ROW LEVEL SECURITY;

ALTER TABLE budgets ENABLE ROW LEVEL SECURITY;
ALTER TABLE budgets FORCE ROW LEVEL SECURITY;

ALTER TABLE goals ENABLE ROW LEVEL SECURITY;
ALTER TABLE goals FORCE ROW LEVEL SECURITY;

ALTER TABLE subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions FORCE ROW LEVEL SECURITY;

ALTER TABLE ai_insights ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_insights FORCE ROW LEVEL SECURITY;

ALTER TABLE password_reset_tokens ENABLE ROW LEVEL SECURITY;
ALTER TABLE password_reset_tokens FORCE ROW LEVEL SECURITY;

-- Drop old policies if rerunning
DROP POLICY IF EXISTS users_owner_rw ON users;
DROP POLICY IF EXISTS transactions_owner_rw ON transactions;
DROP POLICY IF EXISTS categories_owner_or_default_ro_rw ON categories;
DROP POLICY IF EXISTS budgets_owner_rw ON budgets;
DROP POLICY IF EXISTS goals_owner_rw ON goals;
DROP POLICY IF EXISTS subscriptions_owner_rw ON subscriptions;
DROP POLICY IF EXISTS ai_insights_owner_rw ON ai_insights;
DROP POLICY IF EXISTS reset_tokens_owner_rw ON password_reset_tokens;

-- Users: only owner can read/update/delete own row
CREATE POLICY users_owner_rw ON users
FOR ALL
USING (id::text = current_setting('app.user_id', true))
WITH CHECK (id::text = current_setting('app.user_id', true));

-- Transactions
CREATE POLICY transactions_owner_rw ON transactions
FOR ALL
USING (user_id::text = current_setting('app.user_id', true))
WITH CHECK (user_id::text = current_setting('app.user_id', true));

-- Categories:
-- defaults (is_default=true) are readable by all authenticated users;
-- custom categories are private to owner.
CREATE POLICY categories_owner_or_default_ro_rw ON categories
FOR ALL
USING (
  is_default = true
  OR (user_id IS NOT NULL AND user_id::text = current_setting('app.user_id', true))
)
WITH CHECK (
  is_default = false
  AND user_id IS NOT NULL
  AND user_id::text = current_setting('app.user_id', true)
);

-- Budgets
CREATE POLICY budgets_owner_rw ON budgets
FOR ALL
USING (user_id::text = current_setting('app.user_id', true))
WITH CHECK (user_id::text = current_setting('app.user_id', true));

-- Goals
CREATE POLICY goals_owner_rw ON goals
FOR ALL
USING (user_id::text = current_setting('app.user_id', true))
WITH CHECK (user_id::text = current_setting('app.user_id', true));

-- Subscriptions
CREATE POLICY subscriptions_owner_rw ON subscriptions
FOR ALL
USING (user_id::text = current_setting('app.user_id', true))
WITH CHECK (user_id::text = current_setting('app.user_id', true));

-- AI insights
CREATE POLICY ai_insights_owner_rw ON ai_insights
FOR ALL
USING (user_id::text = current_setting('app.user_id', true))
WITH CHECK (user_id::text = current_setting('app.user_id', true));

-- Password reset tokens (ownership by user_id string)
CREATE POLICY reset_tokens_owner_rw ON password_reset_tokens
FOR ALL
USING (user_id = current_setting('app.user_id', true))
WITH CHECK (user_id = current_setting('app.user_id', true));

COMMIT;

