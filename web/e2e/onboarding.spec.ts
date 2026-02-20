import { test, expect, Page } from '@playwright/test';

// Helper: login and get to onboarding state
async function loginAndGetToken(page: Page) {
    await page.goto('/');
    // Remove onboarding_completed to force onboarding flow
    await page.evaluate(() => {
        localStorage.removeItem('onboarding_completed');
        localStorage.setItem('token', 'mock-jwt-token-for-e2e');
    });
    await page.reload();
}

test.describe('Onboarding Flow', () => {
    test.beforeEach(async ({ page }) => {
        await loginAndGetToken(page);
    });

    test('exibe Step 1 — boas-vindas', async ({ page }) => {
        await expect(page.getByText(/bem-vindo|welcome/i)).toBeVisible();
        await expect(page.getByRole('button', { name: /começar|start/i })).toBeVisible();
    });

    test('avança para Step 2 após clicar em começar', async ({ page }) => {
        await page.getByRole('button', { name: /começar|start/i }).click();
        await expect(page.getByText(/método financeiro/i)).toBeVisible();
    });

    test('seleciona método 50-30-20', async ({ page }) => {
        await page.getByRole('button', { name: /começar|start/i }).click();
        await page.getByText('Regra 50-30-20').click();
        // Continue button should be enabled
        await expect(page.getByRole('button', { name: /continuar/i })).toBeEnabled();
    });

    test('mostra modal de detalhes do método', async ({ page }) => {
        await page.getByRole('button', { name: /começar|start/i }).click();
        // Click the info button on the first card
        await page.locator('button').filter({ has: page.getByRole('img', { name: /info/i }) }).first().click();
        await expect(page.getByRole('dialog')).toBeVisible();
    });

    test('avança para Step 3 — renda', async ({ page }) => {
        await page.getByRole('button', { name: /começar|start/i }).click();
        await page.getByText('Regra 50-30-20').click();
        await page.getByRole('button', { name: /continuar/i }).click();
        await expect(page.getByText(/renda mensal/i)).toBeVisible();
    });

    test('mostra preview de distribuição ao digitar renda', async ({ page }) => {
        await page.getByRole('button', { name: /começar|start/i }).click();
        await page.getByText('Regra 50-30-20').click();
        await page.getByRole('button', { name: /continuar/i }).click();
        await page.locator('input[type="number"]').fill('5000');
        await expect(page.getByText('R$ 2.500')).toBeVisible();
    });
});
