package web

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServeIndexWithBasePath_InjectsScript(t *testing.T) {
	indexHTML := []byte(`<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`)

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), basePathContextKey, "/forecastle")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	serveIndexWithBasePath(rec, req, indexHTML)

	body := rec.Body.String()

	if !strings.Contains(body, `window.__BASE_PATH__ = "/forecastle"`) {
		t.Errorf("expected basePath script injection, got: %s", body)
	}

	if !strings.Contains(body, `<script>window.__BASE_PATH__ = "/forecastle";</script></head>`) {
		t.Errorf("expected script before </head>, got: %s", body)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		t.Errorf("expected Content-Type text/html; charset=utf-8, got %s", contentType)
	}
}

func TestServeIndexWithBasePath_EmptyBasePath(t *testing.T) {
	indexHTML := []byte(`<!DOCTYPE html><html><head></head><body></body></html>`)

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), basePathContextKey, "")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	serveIndexWithBasePath(rec, req, indexHTML)

	body := rec.Body.String()

	if !strings.Contains(body, `window.__BASE_PATH__ = ""`) {
		t.Errorf("expected empty basePath script injection, got: %s", body)
	}
}

func TestServeIndexWithBasePath_SpecialCharacters(t *testing.T) {
	indexHTML := []byte(`<!DOCTYPE html><html><head></head><body></body></html>`)

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), basePathContextKey, `/path"with'quotes`)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	serveIndexWithBasePath(rec, req, indexHTML)

	body := rec.Body.String()

	// %q escapes quotes properly
	if !strings.Contains(body, `window.__BASE_PATH__ = "/path\"with'quotes"`) {
		t.Errorf("expected properly escaped basePath, got: %s", body)
	}
}

func TestServeIndexWithBasePath_NoHeadTag(t *testing.T) {
	indexHTML := []byte(`<!DOCTYPE html><html><body>No head tag</body></html>`)

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), basePathContextKey, "/test")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	serveIndexWithBasePath(rec, req, indexHTML)

	// Should still return the original HTML (no replacement happened)
	body := rec.Body.String()
	if body != string(indexHTML) {
		t.Errorf("expected original HTML when no </head> tag, got: %s", body)
	}
}

func TestServeIndexWithBasePath_NestedPath(t *testing.T) {
	indexHTML := []byte(`<!DOCTYPE html><html><head></head><body></body></html>`)

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), basePathContextKey, "/apps/monitoring/forecastle")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	serveIndexWithBasePath(rec, req, indexHTML)

	body := rec.Body.String()

	if !strings.Contains(body, `window.__BASE_PATH__ = "/apps/monitoring/forecastle"`) {
		t.Errorf("expected nested basePath, got: %s", body)
	}
}

func TestReadIndexHTML_Success(t *testing.T) {
	fs := GetFrontendFS()

	content, err := readIndexHTML(fs)
	if err != nil {
		t.Fatalf("expected no error reading index.html, got: %v", err)
	}

	if len(content) == 0 {
		t.Error("expected non-empty index.html content")
	}

	if !strings.Contains(string(content), "<!DOCTYPE html>") && !strings.Contains(string(content), "<!doctype html>") {
		t.Error("expected HTML doctype in index.html")
	}

	if !strings.Contains(string(content), "</head>") {
		t.Error("expected </head> tag in index.html for script injection")
	}
}

func TestIntegration_BasePathInjectesIntoRealIndex(t *testing.T) {
	fs := GetFrontendFS()

	indexHTML, err := readIndexHTML(fs)
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	ctx := context.WithValue(req.Context(), basePathContextKey, "/forecastle")
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()
	serveIndexWithBasePath(rec, req, indexHTML)

	body := rec.Body.String()

	if !strings.Contains(body, `window.__BASE_PATH__ = "/forecastle"`) {
		t.Error("expected basePath injection in real index.html")
	}

	if !strings.Contains(body, "<div id=\"root\">") {
		t.Error("expected React root element to be preserved")
	}
}
