-- Enable and enforce RLS for app tables.
-- Policies use custom setting app.user_id, compatible with backend when set per request.

DO $$
BEGIN
  IF to_regclass('public.users') IS NOT NULL THEN
    EXECUTE 'ALTER TABLE public.users ENABLE ROW LEVEL SECURITY';
    EXECUTE 'ALTER TABLE public.users FORCE ROW LEVEL SECURITY';
    EXECUTE 'DROP POLICY IF EXISTS users_owner_rw ON public.users';
    EXECUTE '
      CREATE POLICY users_owner_rw ON public.users
      FOR ALL
      USING (id::text = current_setting(''app.user_id'', true))
      WITH CHECK (id::text = current_setting(''app.user_id'', true))
    ';
  END IF;

  IF to_regclass('public.transactions') IS NOT NULL THEN
    EXECUTE 'ALTER TABLE public.transactions ENABLE ROW LEVEL SECURITY';
    EXECUTE 'ALTER TABLE public.transactions FORCE ROW LEVEL SECURITY';
    EXECUTE 'DROP POLICY IF EXISTS transactions_owner_rw ON public.transactions';
    EXECUTE '
      CREATE POLICY transactions_owner_rw ON public.transactions
      FOR ALL
      USING (user_id::text = current_setting(''app.user_id'', true))
      WITH CHECK (user_id::text = current_setting(''app.user_id'', true))
    ';
  END IF;

  IF to_regclass('public.categories') IS NOT NULL THEN
    EXECUTE 'ALTER TABLE public.categories ENABLE ROW LEVEL SECURITY';
    EXECUTE 'ALTER TABLE public.categories FORCE ROW LEVEL SECURITY';
    EXECUTE 'DROP POLICY IF EXISTS categories_owner_or_default_ro_rw ON public.categories';
    EXECUTE '
      CREATE POLICY categories_owner_or_default_ro_rw ON public.categories
      FOR ALL
      USING (
        is_default = true
        OR (user_id IS NOT NULL AND user_id::text = current_setting(''app.user_id'', true))
      )
      WITH CHECK (
        is_default = false
        AND user_id IS NOT NULL
        AND user_id::text = current_setting(''app.user_id'', true)
      )
    ';
  END IF;

  IF to_regclass('public.budgets') IS NOT NULL THEN
    EXECUTE 'ALTER TABLE public.budgets ENABLE ROW LEVEL SECURITY';
    EXECUTE 'ALTER TABLE public.budgets FORCE ROW LEVEL SECURITY';
    EXECUTE 'DROP POLICY IF EXISTS budgets_owner_rw ON public.budgets';
    EXECUTE '
      CREATE POLICY budgets_owner_rw ON public.budgets
      FOR ALL
      USING (user_id::text = current_setting(''app.user_id'', true))
      WITH CHECK (user_id::text = current_setting(''app.user_id'', true))
    ';
  END IF;

  IF to_regclass('public.goals') IS NOT NULL THEN
    EXECUTE 'ALTER TABLE public.goals ENABLE ROW LEVEL SECURITY';
    EXECUTE 'ALTER TABLE public.goals FORCE ROW LEVEL SECURITY';
    EXECUTE 'DROP POLICY IF EXISTS goals_owner_rw ON public.goals';
    EXECUTE '
      CREATE POLICY goals_owner_rw ON public.goals
      FOR ALL
      USING (user_id::text = current_setting(''app.user_id'', true))
      WITH CHECK (user_id::text = current_setting(''app.user_id'', true))
    ';
  END IF;

  IF to_regclass('public.subscriptions') IS NOT NULL THEN
    EXECUTE 'ALTER TABLE public.subscriptions ENABLE ROW LEVEL SECURITY';
    EXECUTE 'ALTER TABLE public.subscriptions FORCE ROW LEVEL SECURITY';
    EXECUTE 'DROP POLICY IF EXISTS subscriptions_owner_rw ON public.subscriptions';
    EXECUTE '
      CREATE POLICY subscriptions_owner_rw ON public.subscriptions
      FOR ALL
      USING (user_id::text = current_setting(''app.user_id'', true))
      WITH CHECK (user_id::text = current_setting(''app.user_id'', true))
    ';
  END IF;

  IF to_regclass('public.ai_insights') IS NOT NULL THEN
    EXECUTE 'ALTER TABLE public.ai_insights ENABLE ROW LEVEL SECURITY';
    EXECUTE 'ALTER TABLE public.ai_insights FORCE ROW LEVEL SECURITY';
    EXECUTE 'DROP POLICY IF EXISTS ai_insights_owner_rw ON public.ai_insights';
    EXECUTE '
      CREATE POLICY ai_insights_owner_rw ON public.ai_insights
      FOR ALL
      USING (user_id::text = current_setting(''app.user_id'', true))
      WITH CHECK (user_id::text = current_setting(''app.user_id'', true))
    ';
  END IF;

  IF to_regclass('public.password_reset_tokens') IS NOT NULL THEN
    EXECUTE 'ALTER TABLE public.password_reset_tokens ENABLE ROW LEVEL SECURITY';
    EXECUTE 'ALTER TABLE public.password_reset_tokens FORCE ROW LEVEL SECURITY';
    EXECUTE 'DROP POLICY IF EXISTS reset_tokens_owner_rw ON public.password_reset_tokens';
    EXECUTE '
      CREATE POLICY reset_tokens_owner_rw ON public.password_reset_tokens
      FOR ALL
      USING (user_id = current_setting(''app.user_id'', true))
      WITH CHECK (user_id = current_setting(''app.user_id'', true))
    ';
  END IF;
END $$;

