/// <reference types="vite/client" />
import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  timeout: 30_000,
  expect: { timeout: 5_000 },
  fullyParallel: true,
  forbidOnly: !!import.meta.env.CI,
  retries: import.meta.env.CI ? 2 : 0,
  workers: import.meta.env.CI ? 1 : undefined,

  reporter: [
    ["html", { outputFolder: "playwright-report", open: "never" }],
    ["junit", { outputFile: "playwright-results.xml" }],
    ["list"],
  ],

  use: {
    baseURL: import.meta.env.PLAYWRIGHT_BASE_URL || "http://localhost:5173",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },

  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    {
      name: "firefox",
      use: { ...devices["Desktop Firefox"] },
    },
    {
      name: "mobile-chrome",
      use: { ...devices["Pixel 5"] },
    },
  ],

  webServer: import.meta.env.CI
    ? undefined
    : {
        command: "npm run dev",
        url: "http://localhost:5173",
        reuseExistingServer: true,
        timeout: 30_000,
      },
});
