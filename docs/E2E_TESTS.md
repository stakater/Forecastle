# End-to-End (E2E) Tests

This document describes the end-to-end test suites for Forecastle, covering both backend (Go) and frontend (Playwright) tests.

## Overview

Forecastle has two e2e test suites:

| Suite | Technology | Location | Purpose |
|-------|------------|----------|---------|
| Backend | Go testing | `e2e/` | API endpoints, Kubernetes resource discovery |
| Frontend | Playwright | `frontend/e2e/` | UI functionality, user interactions |

---

## Backend E2E Tests (Go)

**Location:** `e2e/e2e_test.go`

### Prerequisites

- Running Forecastle instance (default: `http://localhost:3000`)
- Optional: Kubernetes cluster (Kind recommended) for discovery tests

### Running the Tests

```bash
# Start Forecastle first
./forecastle --port 3000 &

# Run all e2e tests
make test-e2e

# Or with custom URL
FORECASTLE_URL=http://localhost:3000 go test -v ./e2e/...
```

### Test Categories

#### 1. Health & Readiness Endpoints

| Test | Description |
|------|-------------|
| `TestHealthEndpoint` | Verifies `/healthz` returns 200 OK |
| `TestReadinessEndpoint` | Verifies `/readyz` returns 200 when cache is populated (waits up to 30s) |

#### 2. API Endpoints

| Test | Description |
|------|-------------|
| `TestConfigEndpoint` | Verifies `/api/config` returns valid config with `title` field |
| `TestAppsEndpoint` | Verifies `/api/apps` returns a valid array (may be empty) |

#### 3. Ingress Discovery (Requires Kubernetes Cluster)

| Test | Description |
|------|-------------|
| `TestIngressDiscovery` | Creates an Ingress with Forecastle annotations, verifies it appears in `/api/apps` with correct URL, group, icon, and discovery source |
| `TestIngressWithTLS` | Verifies TLS-enabled Ingresses get `https://` URLs |
| `TestIngressWithCustomAppName` | Verifies `forecastle.stakater.com/appName` annotation overrides the app name |
| `TestIngressWithURLOverride` | Verifies `forecastle.stakater.com/url` annotation overrides the discovered URL |
| `TestAppWithoutExposeAnnotation` | Verifies Ingresses without `forecastle.stakater.com/expose: "true"` are NOT discovered |

#### 4. Multiple Apps

| Test | Description |
|------|-------------|
| `TestMultipleAppsInSameGroup` | Creates 3 Ingresses in the same group, verifies all are discovered and grouped correctly |

#### 5. ForecastleApp CRD (Requires CRD Installed)

| Test | Description |
|------|-------------|
| `TestForecastleAppCRD` | Creates a ForecastleApp custom resource, verifies it appears with `ForecastleAppCRD` discovery source |

#### 6. Custom Apps from Config

| Test | Description |
|------|-------------|
| `TestCustomAppsFromConfig` | Verifies apps defined in `config.yaml` are discovered with `Config` discovery source |
| `TestAllDiscoverySourcesHaveValidFormat` | Validates all apps have valid discovery sources and required fields, logs distribution |

### Cluster-Dependent Tests

Tests that require a Kubernetes cluster use `skipIfNoCluster(t)` to gracefully skip when no cluster is available. This allows API-only tests to run in any environment.

---

## Frontend E2E Tests (Playwright)

**Location:** `frontend/e2e/app.spec.ts`

### Prerequisites

- Running Forecastle instance (default: `http://localhost:3000`)
- Playwright browsers installed (`npx playwright install`)

### Running the Tests

```bash
cd frontend

# Run all tests (headless)
yarn test:e2e

# Run with visible browser
yarn test:e2e:headed

# Run with Playwright UI
yarn test:e2e:ui

# Run specific test file
npx playwright test e2e/app.spec.ts
```

### Browser Coverage

Tests run on multiple browsers and viewports:
- Chromium (Desktop)
- Firefox (Desktop)
- WebKit (Desktop)
- Mobile Chrome (Pixel 5)
- Mobile Safari (iPhone 12)

### Test Categories

#### 1. Forecastle App (Basic Loading)

| Test | Description |
|------|-------------|
| `should load the application` | Verifies page loads with "Forecastle" in title |
| `should display the header with title` | Verifies header and h1 title are visible |
| `should have a search input` | Verifies search textbox is present |
| `should have theme toggle button` | Verifies theme toggle button exists |
| `should have view mode toggle` | Verifies grid/list view toggle buttons exist |

#### 2. API Integration

| Test | Description |
|------|-------------|
| `should fetch and display apps from API` | Intercepts `/api/apps` call, verifies apps render as cards |
| `should fetch config from API` | Intercepts `/api/config` call, verifies it contains `title` property |

#### 3. Search Functionality

| Test | Description |
|------|-------------|
| `should filter apps when searching` | Types a search query, verifies filtering behavior |
| `should allow typing in search input` | Verifies input accepts text, can be cleared |

#### 4. Theme Toggle

| Test | Description |
|------|-------------|
| `should toggle between light and dark themes` | Clicks theme toggle, verifies background color changes |
| `should allow toggling theme multiple times` | Toggles twice, verifies return to initial state |

#### 5. View Mode Toggle

| Test | Description |
|------|-------------|
| `should toggle between grid and list views` | Clicks grid/list buttons, verifies view mode changes |
| `should display apps in grid layout` | Mocks apps, switches to grid view, verifies cards render |
| `should display apps in list layout` | Mocks apps, switches to list view, verifies cards render |

#### 6. App Cards

| Test | Description |
|------|-------------|
| `should display app information correctly` | Mocks API response, verifies app name is displayed |
| `should expand app properties when clicking expand button` | Mocks app with properties, verifies expand reveals Version/Owner |
| `should show network restricted indicator` | Mocks app with `networkRestricted: true`, verifies indicator |
| `should open app in new tab when clicking` | Verifies links have `target="_blank"` and `rel="noopener"` |

#### 7. App Groups

| Test | Description |
|------|-------------|
| `should group apps by namespace/group` | Mocks multiple apps, verifies group headers (production, staging) are visible |

#### 8. Discovery Sources

| Test | Description |
|------|-------------|
| `should display apps from Ingress discovery source` | Mocks Ingress app, verifies it renders correctly |
| `should display apps from Config discovery source` | Mocks Config app (custom app from config.yaml), verifies rendering |
| `should display apps from ForecastleAppCRD discovery source` | Mocks CRD app, verifies rendering and discovery source indicator |
| `should display apps from HTTPRoute discovery source` | Mocks HTTPRoute app (Gateway API), verifies rendering |
| `should display apps from multiple discovery sources together` | Mocks 4 apps from all sources, verifies all render correctly |

#### 9. Responsive Design

| Test | Description |
|------|-------------|
| `should be responsive on mobile viewport` | Sets 375x667 viewport, verifies header and search visible |
| `should be responsive on tablet viewport` | Sets 768x1024 viewport, verifies app renders correctly |

#### 10. Error Handling

| Test | Description |
|------|-------------|
| `should handle API errors gracefully` | Mocks 500 error from `/api/apps`, verifies app doesn't crash |
| `should handle empty apps list` | Mocks empty array from `/api/apps`, verifies no errors |

#### 11. Accessibility

| Test | Description |
|------|-------------|
| `should have proper heading structure` | Verifies h1 is visible |
| `should have accessible form controls` | Verifies search input is focusable |
| `should support keyboard navigation` | Tabs through elements, verifies focus is visible |

#### 12. Health Endpoints

| Test | Description |
|------|-------------|
| `should respond to health check` | Direct request to `/healthz`, verifies 200 OK |
| `should respond to readiness check` | Direct request to `/readyz`, verifies 200 or 503 |

---

## CI/CD Integration

The e2e tests are configured to run in GitHub Actions via `.github/workflows/e2e-tests.yaml`.

### Workflow Triggers
- Pull requests to `master` branch with `ok-to-test` label (same as `run-tests-on-pr`)

### Pipeline Steps

1. **Setup** - Go, Node.js, Kind cluster (via `helm/kind-action`)
2. **Build** - Forecastle binary, frontend dependencies, Playwright browsers
3. **Configure** - Create `config.yaml` with test apps
4. **Start** - Launch Forecastle, wait for health/readiness
5. **Test Backend** - Run `make test-e2e` (API, Ingress, CRD, Config discovery)
6. **Test Frontend** - Run Playwright on Chromium (UI, all discovery sources, accessibility)
7. **Artifacts** - Upload Playwright report (always) and screenshots (on failure)

### Reusable Actions

| Action | Purpose |
|--------|---------|
| `helm/kind-action@v1` | Creates Kind cluster (replaces manual kubectl/kind install) |
| `actions/setup-go@v5` | Go environment |
| `actions/setup-node@v4` | Node.js with yarn cache |
| `actions/upload-artifact@v4` | Test reports |

---

## Test Data

### Mock App Structure

Frontend tests use mocked API responses with this structure:

```json
{
  "name": "Test App",
  "url": "https://test.example.com",
  "icon": "https://example.com/icon.png",
  "group": "test-group",
  "discoverySource": "Ingress",
  "networkRestricted": false,
  "properties": {
    "Version": "1.0.0",
    "Owner": "Platform Team"
  }
}
```

### Ingress Annotations Used in Backend Tests

```yaml
metadata:
  annotations:
    forecastle.stakater.com/expose: "true"
    forecastle.stakater.com/appName: "Custom Name"
    forecastle.stakater.com/group: "my-group"
    forecastle.stakater.com/icon: "https://example.com/icon.png"
    forecastle.stakater.com/url: "https://custom-url.example.com"
```

---

## Troubleshooting

### Backend Tests Skip with "Kubernetes cluster not available"

This is expected behavior. Tests requiring a cluster gracefully skip. To run full tests:

```bash
# Create a Kind cluster
kind create cluster

# Verify connectivity
kubectl cluster-info

# Run tests
make test-e2e
```

### Frontend Tests Fail to Connect

Ensure Forecastle is running:

```bash
# Check if server is running
curl http://localhost:3000/healthz

# If not, start it
./forecastle --port 3000
```

### Playwright Browsers Not Installed

```bash
cd frontend
npx playwright install
```

### Tests Timeout Waiting for Cache

The backend uses a 20-second cache refresh interval. Tests wait up to 30 seconds for readiness. If tests timeout, check:

1. Forecastle logs for errors
2. Kubernetes connectivity (for discovery tests)
3. Config file availability
