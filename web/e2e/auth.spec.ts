import { test, expect } from "@playwright/test";

test.describe("Authentication", () => {
  test("exibe tela de login por padrão", async ({ page }) => {
    await page.goto("/");
    await expect(
      page.getByRole("heading", { name: /entrar|login/i }),
    ).toBeVisible();
  });

  test("mostra erro com credenciais inválidas", async ({ page }) => {
    await page.goto("/");
    await page.getByPlaceholder(/email/i).fill("wrong@test.com");
    await page.getByPlaceholder(/senha|password/i).fill("wrongpassword");
    await page.getByRole("button", { name: /entrar|login/i }).click();
    await expect(page.getByText(/inválid|invalid|credencial/i)).toBeVisible({
      timeout: 5000,
    });
  });

  test("navega para signup", async ({ page }) => {
    await page.goto("/");
    await page.getByText(/criar conta|sign up|cadastrar/i).click();
    await expect(
      page.getByRole("heading", { name: /criar conta|cadastro|sign up/i }),
    ).toBeVisible();
  });

  test("navega para forgot password", async ({ page }) => {
    await page.goto("/");
    await page.getByText(/esqueceu|forgot/i).click();
    await expect(page.getByPlaceholder(/email/i)).toBeVisible();
  });
});
