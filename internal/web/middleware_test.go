package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasePathMiddleware_NoHeaderNoConfig(t *testing.T) {
	var capturedPath string
	var capturedBasePath string

	handler := BasePathMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedBasePath = GetBasePath(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/apps", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedPath != "/api/apps" {
		t.Errorf("expected path /api/apps, got %s", capturedPath)
	}
	if capturedBasePath != "" {
		t.Errorf("expected empty basePath, got %s", capturedBasePath)
	}
}

func TestBasePathMiddleware_HeaderPresent(t *testing.T) {
	var capturedPath string
	var capturedBasePath string

	handler := BasePathMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedBasePath = GetBasePath(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/forecastle/api/apps", nil)
	req.Header.Set("X-Forwarded-Prefix", "/forecastle")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedPath != "/api/apps" {
		t.Errorf("expected path /api/apps, got %s", capturedPath)
	}
	if capturedBasePath != "/forecastle" {
		t.Errorf("expected basePath /forecastle, got %s", capturedBasePath)
	}
}

func TestBasePathMiddleware_HeaderWithTrailingSlash(t *testing.T) {
	var capturedBasePath string

	handler := BasePathMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBasePath = GetBasePath(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/forecastle/api/apps", nil)
	req.Header.Set("X-Forwarded-Prefix", "/forecastle/")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedBasePath != "/forecastle" {
		t.Errorf("expected basePath /forecastle (no trailing slash), got %s", capturedBasePath)
	}
}

func TestBasePathMiddleware_ConfiguredBasePath(t *testing.T) {
	var capturedPath string
	var capturedBasePath string

	handler := BasePathMiddleware("/dashboard")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedBasePath = GetBasePath(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/dashboard/api/apps", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedPath != "/api/apps" {
		t.Errorf("expected path /api/apps, got %s", capturedPath)
	}
	if capturedBasePath != "/dashboard" {
		t.Errorf("expected basePath /dashboard, got %s", capturedBasePath)
	}
}

func TestBasePathMiddleware_HeaderTakesPriority(t *testing.T) {
	var capturedBasePath string

	handler := BasePathMiddleware("/configured")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedBasePath = GetBasePath(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/from-header/api/apps", nil)
	req.Header.Set("X-Forwarded-Prefix", "/from-header")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedBasePath != "/from-header" {
		t.Errorf("expected basePath /from-header (from header), got %s", capturedBasePath)
	}
}

func TestBasePathMiddleware_PathNotMatchingPrefix(t *testing.T) {
	var capturedPath string

	handler := BasePathMiddleware("/forecastle")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/other/path", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedPath != "/other/path" {
		t.Errorf("expected path /other/path (unchanged), got %s", capturedPath)
	}
}

func TestBasePathMiddleware_RootPath(t *testing.T) {
	var capturedPath string

	handler := BasePathMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/forecastle", nil)
	req.Header.Set("X-Forwarded-Prefix", "/forecastle")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedPath != "/" {
		t.Errorf("expected path /, got %s", capturedPath)
	}
}

func TestBasePathMiddleware_NestedPath(t *testing.T) {
	var capturedPath string
	var capturedBasePath string

	handler := BasePathMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedBasePath = GetBasePath(r)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/apps/forecastle/api/config", nil)
	req.Header.Set("X-Forwarded-Prefix", "/apps/forecastle")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedPath != "/api/config" {
		t.Errorf("expected path /api/config, got %s", capturedPath)
	}
	if capturedBasePath != "/apps/forecastle" {
		t.Errorf("expected basePath /apps/forecastle, got %s", capturedBasePath)
	}
}

func TestBasePathMiddleware_StaticAssets(t *testing.T) {
	var capturedPath string

	handler := BasePathMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/forecastle/static/js/main.js", nil)
	req.Header.Set("X-Forwarded-Prefix", "/forecastle")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedPath != "/static/js/main.js" {
		t.Errorf("expected path /static/js/main.js, got %s", capturedPath)
	}
}

func TestGetBasePath_NoContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/apps", nil)

	basePath := GetBasePath(req)

	if basePath != "" {
		t.Errorf("expected empty basePath, got %s", basePath)
	}
}

func TestBasePathMiddleware_EmptyPath(t *testing.T) {
	var capturedPath string

	handler := BasePathMiddleware("/forecastle")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/forecastle", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedPath != "/" {
		t.Errorf("expected path /, got %s", capturedPath)
	}
}

func TestBasePathMiddleware_PreservesMethod(t *testing.T) {
	var capturedMethod string

	handler := BasePathMiddleware("/api")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/apps", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedMethod != "POST" {
		t.Errorf("expected method POST, got %s", capturedMethod)
	}
}

func TestBasePathMiddleware_PreservesHeaders(t *testing.T) {
	var capturedHeader string

	handler := BasePathMiddleware("")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeader = r.Header.Get("X-Custom-Header")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/apps", nil)
	req.Header.Set("X-Custom-Header", "test-value")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if capturedHeader != "test-value" {
		t.Errorf("expected header test-value, got %s", capturedHeader)
	}
}
