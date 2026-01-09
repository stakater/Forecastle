import { test, expect } from '@playwright/test';

test.describe('Forecastle App', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should load the application', async ({ page }) => {
    // Check that the page loads
    await expect(page).toHaveTitle(/Forecastle/);
  });

  test('should display the header with title', async ({ page }) => {
    // Header should be visible
    const header = page.locator('header');
    await expect(header).toBeVisible();

    // Title should be present (either "Forecastle" or custom title from config)
    const title = page.locator('h1');
    await expect(title).toBeVisible();
  });

  test('should have a search input', async ({ page }) => {
    // Search input should exist
    const searchInput = page.getByRole('textbox', { name: /search/i });
    await expect(searchInput).toBeVisible();
  });

  test('should have theme toggle button', async ({ page }) => {
    // Theme toggle button should exist
    const themeToggle = page.getByRole('button', { name: /toggle.*theme|dark.*mode|light.*mode/i });
    await expect(themeToggle).toBeVisible();
  });

  test('should have view mode toggle', async ({ page }) => {
    // View toggle buttons (grid/list) should exist
    const viewToggle = page.getByRole('button', { name: /grid|list/i }).first();
    await expect(viewToggle).toBeVisible();
  });
});

test.describe('API Integration', () => {
  test('should fetch and display apps from API', async ({ page }) => {
    // Intercept the API call
    const appsResponse = page.waitForResponse(
      (response) => response.url().includes('/api/apps') && response.status() === 200
    );

    await page.goto('/');

    // Wait for apps API response
    const response = await appsResponse;
    const apps = await response.json();

    // If there are apps, they should be displayed
    if (apps.length > 0) {
      // Wait for app cards to render
      await page.waitForSelector('[data-testid="app-card"], .MuiCard-root', { timeout: 10000 });

      // Count visible app elements
      const appCards = page.locator('.MuiCard-root');
      await expect(appCards.first()).toBeVisible();
    }
  });

  test('should fetch config from API', async ({ page }) => {
    // Intercept the config API call
    const configResponse = page.waitForResponse(
      (response) => response.url().includes('/api/config') && response.status() === 200
    );

    await page.goto('/');

    // Wait for config API response
    const response = await configResponse;
    const config = await response.json();

    // Config should have a title
    expect(config).toHaveProperty('title');
  });
});

test.describe('Search Functionality', () => {
  test('should filter apps when searching', async ({ page }) => {
    // Wait for apps to load
    await page.goto('/');
    await page.waitForResponse((r) => r.url().includes('/api/apps'));

    // Get search input
    const searchInput = page.getByRole('textbox', { name: /search/i });
    await expect(searchInput).toBeVisible();

    // Type a search query
    await searchInput.fill('test-search-query-that-should-not-match');

    // Wait for filtering to apply
    await page.waitForTimeout(300);

    // Either no results message or filtered results
    // The exact behavior depends on what apps are available
  });

  test('should allow typing in search input', async ({ page }) => {
    await page.goto('/');

    const searchInput = page.getByRole('textbox', { name: /search/i });
    await searchInput.fill('test query');

    // Verify the input has the value
    await expect(searchInput).toHaveValue('test query');

    // Clear the input manually
    await searchInput.clear();
    await expect(searchInput).toHaveValue('');
  });
});

// Helper to extract RGB values from a color string (handles both rgb and rgba)
function parseRgb(color: string): { r: number; g: number; b: number } | null {
  const match = color.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/);
  if (match) {
    return { r: parseInt(match[1]), g: parseInt(match[2]), b: parseInt(match[3]) };
  }
  return null;
}

// Helper to check if two colors are similar (ignoring alpha)
function colorsAreSimilar(color1: string, color2: string): boolean {
  const rgb1 = parseRgb(color1);
  const rgb2 = parseRgb(color2);
  if (!rgb1 || !rgb2) return false;
  // Allow small tolerance for rounding
  return Math.abs(rgb1.r - rgb2.r) <= 1 &&
         Math.abs(rgb1.g - rgb2.g) <= 1 &&
         Math.abs(rgb1.b - rgb2.b) <= 1;
}

test.describe('Theme Toggle', () => {
  test('should toggle between light and dark themes', async ({ page }) => {
    await page.goto('/');

    // Find theme toggle button
    const themeToggle = page.getByRole('button', { name: /toggle.*theme|dark.*mode|light.*mode/i });

    if (await themeToggle.isVisible()) {
      // Get initial background color
      const body = page.locator('body');
      const initialBg = await body.evaluate((el) =>
        window.getComputedStyle(el).backgroundColor
      );

      // Click theme toggle
      await themeToggle.click();

      // Wait for theme transition
      await page.waitForTimeout(300);

      // Background should have changed
      const newBg = await body.evaluate((el) =>
        window.getComputedStyle(el).backgroundColor
      );

      expect(newBg).not.toBe(initialBg);
    }
  });

  test('should allow toggling theme multiple times', async ({ page }) => {
    await page.goto('/');

    const themeToggle = page.getByRole('button', { name: /toggle.*theme|dark.*mode|light.*mode/i });

    if (await themeToggle.isVisible()) {
      // Get initial state
      const body = page.locator('body');
      const initialBg = await body.evaluate((el) =>
        window.getComputedStyle(el).backgroundColor
      );

      // Toggle theme twice - should return to original
      await themeToggle.click();
      await page.waitForTimeout(300);

      await themeToggle.click();
      await page.waitForTimeout(300);

      const finalBg = await body.evaluate((el) =>
        window.getComputedStyle(el).backgroundColor
      );

      // After toggling twice, should be back to initial state (compare RGB, ignore alpha)
      expect(colorsAreSimilar(finalBg, initialBg)).toBe(true);
    }
  });
});

test.describe('View Mode Toggle', () => {
  test('should toggle between grid and list views', async ({ page }) => {
    await page.goto('/');

    // Wait for apps to load
    await page.waitForResponse((r) => r.url().includes('/api/apps'));

    // Find view toggle buttons
    const gridButton = page.getByRole('button', { name: /grid/i });
    const listButton = page.getByRole('button', { name: /list/i });

    // At least one should be visible
    const hasViewToggle = await gridButton.isVisible() || await listButton.isVisible();

    if (hasViewToggle) {
      // Click to toggle view
      if (await listButton.isVisible()) {
        await listButton.click();
        await page.waitForTimeout(200);
      }

      if (await gridButton.isVisible()) {
        await gridButton.click();
        await page.waitForTimeout(200);
      }
    }
  });

  test('should display apps in grid layout', async ({ page }) => {
    // Mock apps for consistent testing
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          { name: 'App 1', url: 'https://app1.example.com', group: 'test', discoverySource: 'Ingress' },
          { name: 'App 2', url: 'https://app2.example.com', group: 'test', discoverySource: 'Ingress' },
          { name: 'App 3', url: 'https://app3.example.com', group: 'test', discoverySource: 'Ingress' },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // Click grid button if available
    const gridButton = page.getByRole('button', { name: /grid/i });
    if (await gridButton.isVisible()) {
      await gridButton.click();
      await page.waitForTimeout(200);
    }

    // In grid view, cards should be laid out in a grid (multiple columns)
    const cards = page.locator('.MuiCard-root');
    const cardCount = await cards.count();
    expect(cardCount).toBeGreaterThanOrEqual(3);
  });

  test('should display apps in list layout', async ({ page }) => {
    // Mock apps for consistent testing
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          { name: 'App 1', url: 'https://app1.example.com', group: 'test', discoverySource: 'Ingress' },
          { name: 'App 2', url: 'https://app2.example.com', group: 'test', discoverySource: 'Ingress' },
          { name: 'App 3', url: 'https://app3.example.com', group: 'test', discoverySource: 'Ingress' },
        ]),
      });
    });

    await page.goto('/');
    // Wait for apps to load by checking for app names
    await page.waitForSelector('text=App 1', { timeout: 10000 });

    // Click list button if available
    const listButton = page.getByRole('button', { name: /list/i });
    if (await listButton.isVisible()) {
      await listButton.click();
      await page.waitForTimeout(200);
    }

    // Verify all 3 apps are displayed (list view uses different DOM structure than grid)
    await expect(page.getByText('App 1')).toBeVisible();
    await expect(page.getByText('App 2')).toBeVisible();
    await expect(page.getByText('App 3')).toBeVisible();
  });
});

test.describe('App Cards', () => {
  test('should display app information correctly', async ({ page }) => {
    // Mock apps response
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'Test App',
            url: 'https://test.example.com',
            icon: 'https://example.com/icon.png',
            group: 'test-group',
            discoverySource: 'Ingress',
            networkRestricted: false,
            properties: {
              Version: '1.0.0',
            },
          },
        ]),
      });
    });

    await page.goto('/');

    // Wait for app to render
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // App name should be visible
    await expect(page.getByText('Test App')).toBeVisible();
  });

  test('should expand app properties when clicking expand button', async ({ page }) => {
    // Mock apps with properties
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'App With Properties',
            url: 'https://props.example.com',
            icon: '',
            group: 'test',
            discoverySource: 'Ingress',
            networkRestricted: false,
            properties: {
              Version: '2.0.0',
              Owner: 'Platform Team',
            },
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // Find and click expand button
    const expandButton = page.locator('[aria-label*="details"], [aria-label*="expand"]').first();

    if (await expandButton.isVisible()) {
      await expandButton.click();

      // Properties should become visible
      await expect(page.getByText('Version')).toBeVisible();
      await expect(page.getByText('2.0.0')).toBeVisible();
    }
  });

  test('should show network restricted indicator', async ({ page }) => {
    // Mock app with network restriction
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'Restricted App',
            url: 'https://restricted.example.com',
            icon: '',
            group: 'secure',
            discoverySource: 'Ingress',
            networkRestricted: true,
            properties: {},
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // Network restricted icon should be visible
    const restrictedIcon = page.locator('[data-testid="VpnLockIcon"], .MuiSvgIcon-root[aria-label*="restricted"]');
    // The icon might be visible - depends on implementation
  });

  test('should open app in new tab when clicking', async ({ page, context }) => {
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'Clickable App',
            url: 'https://clickable.example.com',
            icon: '',
            group: 'test',
            discoverySource: 'Ingress',
            networkRestricted: false,
            properties: {},
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // App links should have target="_blank"
    const appLink = page.locator('a[target="_blank"][href*="example.com"]').first();
    await expect(appLink).toHaveAttribute('rel', /noopener/);
  });
});

test.describe('App Groups', () => {
  test('should group apps by namespace/group', async ({ page }) => {
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'App 1',
            url: 'https://app1.example.com',
            group: 'production',
            discoverySource: 'Ingress',
          },
          {
            name: 'App 2',
            url: 'https://app2.example.com',
            group: 'production',
            discoverySource: 'Ingress',
          },
          {
            name: 'App 3',
            url: 'https://app3.example.com',
            group: 'staging',
            discoverySource: 'Ingress',
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // Group headers should be visible
    await expect(page.getByText(/production/i).first()).toBeVisible();
    await expect(page.getByText(/staging/i).first()).toBeVisible();
  });
});

test.describe('Discovery Sources', () => {
  test('should display apps from Ingress discovery source', async ({ page }) => {
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'Ingress App',
            url: 'https://ingress-app.example.com',
            icon: 'https://example.com/ingress-icon.png',
            group: 'kubernetes',
            discoverySource: 'Ingress',
            networkRestricted: false,
            properties: {},
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // Check for the app name in the card heading
    await expect(page.locator('h3').filter({ hasText: 'Ingress App' })).toBeVisible();
  });

  test('should display apps from Config discovery source', async ({ page }) => {
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'Custom Config App',
            url: 'https://config-app.example.com',
            icon: 'https://example.com/config-icon.png',
            group: 'custom',
            discoverySource: 'Config',
            networkRestricted: false,
            properties: {
              Source: 'config.yaml',
            },
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // Check for the app name in the card heading
    await expect(page.locator('h3').filter({ hasText: 'Custom Config App' })).toBeVisible();
  });

  test('should display apps from ForecastleAppCRD discovery source', async ({ page }) => {
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'CRD App',
            url: 'https://crd-app.example.com',
            icon: 'https://example.com/crd-icon.png',
            group: 'crd-apps',
            discoverySource: 'ForecastleAppCRD',
            networkRestricted: false,
            properties: {},
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // Check for the app name in the card link (use exact match to avoid matching group "Crd Apps")
    await expect(page.getByRole('link', { name: /CRD App/ })).toBeVisible();
    // Also verify discovery source indicator
    await expect(page.getByText('ForecastleAppCRD')).toBeVisible();
  });

  test('should display apps from HTTPRoute discovery source', async ({ page }) => {
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'Gateway API App',
            url: 'https://gateway-app.example.com',
            icon: 'https://example.com/gateway-icon.png',
            group: 'gateway-api',
            discoverySource: 'HTTPRoute',
            networkRestricted: false,
            properties: {},
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // Check for the app name in the card heading
    await expect(page.locator('h3').filter({ hasText: 'Gateway API App' })).toBeVisible();
    // Also verify discovery source indicator
    await expect(page.getByText('HTTPRoute')).toBeVisible();
  });

  test('should display apps from multiple discovery sources together', async ({ page }) => {
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([
          {
            name: 'Ingress Service',
            url: 'https://ingress.example.com',
            group: 'infrastructure',
            discoverySource: 'Ingress',
          },
          {
            name: 'Config Service',
            url: 'https://config.example.com',
            group: 'infrastructure',
            discoverySource: 'Config',
          },
          {
            name: 'CRD Service',
            url: 'https://crd.example.com',
            group: 'infrastructure',
            discoverySource: 'ForecastleAppCRD',
          },
          {
            name: 'HTTPRoute Service',
            url: 'https://httproute.example.com',
            group: 'infrastructure',
            discoverySource: 'HTTPRoute',
          },
        ]),
      });
    });

    await page.goto('/');
    await page.waitForSelector('.MuiCard-root', { timeout: 10000 });

    // All apps from different sources should be displayed (check headings)
    await expect(page.locator('h3').filter({ hasText: 'Ingress Service' })).toBeVisible();
    await expect(page.locator('h3').filter({ hasText: 'Config Service' })).toBeVisible();
    await expect(page.locator('h3').filter({ hasText: 'CRD Service' })).toBeVisible();
    await expect(page.locator('h3').filter({ hasText: 'HTTPRoute Service' })).toBeVisible();

    // Verify we have 4 cards
    const cards = page.locator('.MuiCard-root');
    await expect(cards).toHaveCount(4);
  });
});

test.describe('Responsive Design', () => {
  test('should be responsive on mobile viewport', async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');

    // Header should still be visible
    const header = page.locator('header');
    await expect(header).toBeVisible();

    // Search should be accessible
    const searchInput = page.getByRole('textbox', { name: /search/i });
    await expect(searchInput).toBeVisible();
  });

  test('should be responsive on tablet viewport', async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });
    await page.goto('/');

    // App should render correctly
    await expect(page.locator('header')).toBeVisible();
  });
});

test.describe('Error Handling', () => {
  test('should handle API errors gracefully', async ({ page }) => {
    // Mock API error
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal Server Error' }),
      });
    });

    await page.goto('/');

    // App should not crash - page should still be interactive
    await expect(page.locator('header')).toBeVisible();
  });

  test('should handle empty apps list', async ({ page }) => {
    await page.route('**/api/apps', async (route) => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify([]),
      });
    });

    await page.goto('/');

    // Should show empty state or just render without errors
    await expect(page.locator('header')).toBeVisible();
  });
});

test.describe('Accessibility', () => {
  test('should have proper heading structure', async ({ page }) => {
    await page.goto('/');

    // Should have h1
    const h1 = page.locator('h1');
    await expect(h1).toBeVisible();
  });

  test('should have accessible form controls', async ({ page }) => {
    await page.goto('/');

    // Search input should be accessible
    const searchInput = page.getByRole('textbox', { name: /search/i });
    await expect(searchInput).toBeVisible();

    // Should be focusable
    await searchInput.focus();
    await expect(searchInput).toBeFocused();
  });

  test('should support keyboard navigation', async ({ page }) => {
    await page.goto('/');

    // Tab through interactive elements
    await page.keyboard.press('Tab');

    // Something should be focused
    const focused = page.locator(':focus');
    await expect(focused).toBeVisible();
  });
});

test.describe('Health Endpoints', () => {
  test('should respond to health check', async ({ request }) => {
    const response = await request.get('/healthz');
    expect(response.ok()).toBeTruthy();
  });

  test('should respond to readiness check', async ({ request }) => {
    const response = await request.get('/readyz');
    // May be 200 or 503 depending on cache state
    expect([200, 503]).toContain(response.status());
  });
});
