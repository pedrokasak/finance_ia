import * as zod from "zod";

export const envSchema = zod.object({
  VITE_API_BASE_URL: zod.string().url(),
  VITE_API_TIMEOUT: zod.string().optional(),
});

export type Env = zod.infer<typeof envSchema>;

const parsed = envSchema.safeParse(import.meta.env);

if (!parsed.success) {
  throw new Error("Invalid environment variables");
}

export const env: Env = parsed.data;
