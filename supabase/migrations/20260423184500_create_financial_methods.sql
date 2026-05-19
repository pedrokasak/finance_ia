-- Ensure onboarding financial methods table exists before seed execution.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS public.financial_methods (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  key text UNIQUE NOT NULL,
  name text NOT NULL,
  tagline text NOT NULL,
  description text NOT NULL,
  for_who text NOT NULL,
  icon text NOT NULL,
  color text NOT NULL,
  bg text NOT NULL,
  split_raw text NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_financial_methods_is_active
  ON public.financial_methods (is_active);
