import { test, expect } from '@playwright/test'

test('homepage displays article cards', async ({ page }) => {
  await page.goto('/')
  // wait for news grid or a loading timeout
  await page.waitForSelector('.news-grid .card', { timeout: 15000 })
  const count = await page.locator('.news-grid .card').count()
  expect(count).toBeGreaterThan(0)
})
