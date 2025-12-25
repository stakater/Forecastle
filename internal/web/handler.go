package web

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/stakater/Forecastle/v1/pkg/config"
	"github.com/stakater/Forecastle/v1/pkg/forecastle"
	"github.com/stakater/Forecastle/v1/pkg/forecastle/crdapps"
	"github.com/stakater/Forecastle/v1/pkg/forecastle/customapps"
	"github.com/stakater/Forecastle/v1/pkg/forecastle/ingressapps"
	"github.com/stakater/Forecastle/v1/pkg/kube"
	"github.com/stakater/Forecastle/v1/pkg/kube/util"
	"github.com/stakater/Forecastle/v1/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var logger = log.New()

// AppsResponse is the response structure for the /api/apps endpoint
type AppsResponse struct {
	Apps      []forecastle.App `json:"apps"`
	CachedAt  time.Time        `json:"cachedAt"`
	ExpiresAt time.Time        `json:"expiresAt"`
}

// Handler handles HTTP requests with background caching
type Handler struct {
	clients       *kube.Clients
	configFunc    func() (*config.Config, error)
	cacheInterval time.Duration

	// Cached apps data
	appsCache     []forecastle.App
	appsCacheMu   sync.RWMutex
	appsCacheTime time.Time

	// Cached config
	configCache   *config.Config
	configCacheMu sync.RWMutex
}

// NewHandler creates a new Handler instance
func NewHandler(clients *kube.Clients, configFunc func() (*config.Config, error), cacheInterval time.Duration) *Handler {
	return &Handler{
		clients:       clients,
		configFunc:    configFunc,
		cacheInterval: cacheInterval,
	}
}

// StartBackgroundCache starts the background cache refresh goroutine
func (h *Handler) StartBackgroundCache(ctx context.Context) {
	// Initial load
	h.refreshCache(ctx)

	// Periodic refresh
	ticker := time.NewTicker(h.cacheInterval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				h.refreshCache(ctx)
			}
		}
	}()

	logger.Info("Background cache started with interval: ", h.cacheInterval)
}

func (h *Handler) refreshCache(ctx context.Context) {
	// Refresh config
	cfg, err := h.configFunc()
	if err != nil {
		logger.Error("Failed to refresh config cache: ", err)
	} else {
		h.configCacheMu.Lock()
		h.configCache = cfg
		h.configCacheMu.Unlock()
	}

	// Refresh apps
	apps, err := h.discoverApps(cfg)
	if err != nil {
		logger.Error("Failed to refresh apps cache: ", err)
		return
	}

	h.appsCacheMu.Lock()
	h.appsCache = apps
	h.appsCacheTime = time.Now()
	h.appsCacheMu.Unlock()

	logger.Info("Cache refreshed with ", len(apps), " apps")
}

func (h *Handler) discoverApps(cfg *config.Config) ([]forecastle.App, error) {
	if cfg == nil {
		var err error
		cfg, err = h.configFunc()
		if err != nil {
			return nil, err
		}
	}

	namespaces, err := util.PopulateNamespaceList(h.clients.KubernetesClient, cfg.NamespaceSelector)
	if err != nil {
		return nil, err
	}

	var namespacesString string
	if len(namespaces) == 1 && namespaces[0] == metav1.NamespaceAll {
		namespacesString = "* (All Namespaces)"
	} else {
		namespacesString = strings.Join(namespaces, ",")
	}
	logger.Info("Looking for forecastle apps in namespaces: " + namespacesString)

	var allApps []forecastle.App

	// Discover from Ingress resources
	ingressAppsList := ingressapps.NewList(h.clients.KubernetesClient, *cfg)
	ingressApps, err := ingressAppsList.Populate(namespaces...).Get()
	if err != nil {
		logger.Error("Error discovering ingress apps: ", err)
	} else {
		allApps = append(allApps, ingressApps...)
	}

	// Discover from custom apps config
	customAppsList := customapps.NewList(*cfg)
	customApps, err := customAppsList.Populate().Get()
	if err != nil {
		logger.Error("Error discovering custom apps: ", err)
	} else {
		allApps = append(allApps, customApps...)
	}

	// Discover from CRD if enabled
	if cfg.CRDEnabled {
		crdAppsList := crdapps.NewList(*h.clients, *cfg)
		crdApps, err := crdAppsList.Populate(namespaces...).Get()
		if err != nil {
			logger.Error("Error discovering CRD apps: ", err)
		} else {
			allApps = append(allApps, crdApps...)
		}
	}

	if allApps == nil {
		allApps = []forecastle.App{}
	}

	return allApps, nil
}

// AppsHandler handles GET /api/apps
func (h *Handler) AppsHandler(w http.ResponseWriter, r *http.Request) {
	h.appsCacheMu.RLock()
	apps := h.appsCache
	cacheTime := h.appsCacheTime
	h.appsCacheMu.RUnlock()

	// Ensure apps is never nil
	if apps == nil {
		apps = []forecastle.App{}
	}

	response := AppsResponse{
		Apps:      apps,
		CachedAt:  cacheTime,
		ExpiresAt: cacheTime.Add(h.cacheInterval),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response.Apps); err != nil {
		logger.Error("Error encoding apps response: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ConfigHandler handles GET /api/config
func (h *Handler) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	h.configCacheMu.RLock()
	cfg := h.configCache
	h.configCacheMu.RUnlock()

	if cfg == nil {
		// Fallback to direct fetch if cache not ready
		var err error
		cfg, err = h.configFunc()
		if err != nil {
			logger.Error("Failed to get config: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cfg); err != nil {
		logger.Error("Error encoding config response: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// HealthzHandler handles GET /healthz (liveness probe)
func (h *Handler) HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// ReadyzHandler handles GET /readyz (readiness probe)
func (h *Handler) ReadyzHandler(w http.ResponseWriter, r *http.Request) {
	h.appsCacheMu.RLock()
	hasCache := !h.appsCacheTime.IsZero()
	h.appsCacheMu.RUnlock()

	if hasCache {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("cache not ready"))
	}
}
