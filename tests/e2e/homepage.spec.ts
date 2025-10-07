import { test, expect } from '@playwright/test'

test.describe('Beer Festival Homepage', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/')
  })

  test('should load the homepage successfully', async ({ page }) => {
    await expect(page).toHaveTitle(/Vite \+ Vue \+ TS/)
  })

  test('should display the site header', async ({ page }) => {
    const header = page.locator('header')
    await expect(header).toBeVisible()

    const title = page.getByTestId('site-title')
    await expect(title).toBeVisible()
    await expect(title).toHaveText('Festivals de Bière')
  })

  test('should display the tagline', async ({ page }) => {
    const tagline = page.getByTestId('site-tagline')
    await expect(tagline).toBeVisible()
    await expect(tagline).toContainText('Découvrez les meilleurs festivals de France')
  })

  test('should display the next festival section', async ({ page }) => {
    const nextFestivalSection = page.getByTestId('next-festival-section')
    await expect(nextFestivalSection).toBeVisible()

    const sectionTitle = nextFestivalSection.getByTestId('section-title')
    await expect(sectionTitle).toHaveText('Prochain Festival')
  })

  test('should display the map section', async ({ page }) => {
    const mapSection = page.getByTestId('festival-map-section')
    await expect(mapSection).toBeVisible()

    const sectionTitle = mapSection.getByTestId('section-title')
    await expect(sectionTitle).toHaveText('Carte des Festivals')
  })

  test('should display the festival list section', async ({ page }) => {
    const listSection = page.getByTestId('festival-list-section')
    await expect(listSection).toBeVisible()

    const sectionTitle = listSection.getByTestId('section-title')
    await expect(sectionTitle).toHaveText('Tous les Festivals')
  })

  test('should display multiple festival cards', async ({ page }) => {
    const festivalCards = page.locator('[data-testid^="festival-card-"]')
    const count = await festivalCards.count()
    expect(count).toBeGreaterThan(0)
  })

  test('should display festival card with details', async ({ page }) => {
    const firstCard = page.locator('[data-testid^="festival-card-"]').first()
    await expect(firstCard).toBeVisible()

    const festivalName = firstCard.getByTestId('festival-name')
    await expect(festivalName).toBeVisible()

    const festivalLocation = firstCard.getByTestId('festival-location')
    await expect(festivalLocation).toBeVisible()

    const festivalDates = firstCard.getByTestId('festival-dates')
    await expect(festivalDates).toBeVisible()
  })

  test('should display the footer', async ({ page }) => {
    const footer = page.locator('footer')
    await expect(footer).toBeVisible()

    const footerText = page.getByTestId('footer-text')
    await expect(footerText).toBeVisible()
    await expect(footerText).toContainText('© 2025 Festivals de Bière')
  })

  test('should have responsive layout', async ({ page }) => {
    await page.setViewportSize({ width: 1920, height: 1080 })
    const header = page.locator('header')
    await expect(header).toBeVisible()

    await page.setViewportSize({ width: 375, height: 667 })
    await expect(header).toBeVisible()
  })

  test('should scroll to festival list section', async ({ page }) => {
    const listSection = page.getByTestId('festival-list-section')
    await listSection.scrollIntoViewIfNeeded()
    await expect(listSection).toBeInViewport()
  })

  test('should display the logo icon', async ({ page }) => {
    const logo = page.getByTestId('logo-icon')
    await expect(logo).toBeVisible()
  })
})
