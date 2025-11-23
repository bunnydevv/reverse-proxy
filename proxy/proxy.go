package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/bunnydevv/reverse-proxy/config"
)

type ReverseProxy struct {
	config       *config.Config
	server       *http.Server
	backends     []*Backend
	loadBalancer LoadBalancer
	healthCheck  *HealthChecker
	mu           sync.RWMutex
}

type Backend struct {
	URL          *url.URL
	Proxy        *httputil.ReverseProxy
	Alive        bool
	Weight       int
	Connections  int
	mu           sync.RWMutex
}

func New(cfg *config.Config) (*ReverseProxy, error) {
	if len(cfg.Backends) == 0 {
		return nil, fmt.Errorf("no backends configured")
	}

	rp := &ReverseProxy{
		config:   cfg,
		backends: make([]*Backend, 0, len(cfg.Backends)),
	}

	// Initialize backends
	for _, b := range cfg.Backends {
		backendURL, err := url.Parse(b.URL)
		if err != nil {
			return nil, fmt.Errorf("invalid backend URL %s: %w", b.URL, err)
		}

		weight := b.Weight
		if weight == 0 {
			weight = 1
		}

		backend := &Backend{
			URL:    backendURL,
			Proxy:  httputil.NewSingleHostReverseProxy(backendURL),
			Alive:  true,
			Weight: weight,
		}

		// Customize error handler
		backend.Proxy.ErrorHandler = rp.errorHandler

		rp.backends = append(rp.backends, backend)
	}

	// Initialize load balancer
	switch cfg.LoadBalancer.Algorithm {
	case "round-robin":
		rp.loadBalancer = NewRoundRobinBalancer(rp.backends)
	case "least-connections":
		rp.loadBalancer = NewLeastConnectionsBalancer(rp.backends)
	case "weighted":
		rp.loadBalancer = NewWeightedBalancer(rp.backends)
	default:
		rp.loadBalancer = NewRoundRobinBalancer(rp.backends)
	}

	// Initialize health checker
	if cfg.HealthCheck.Enabled {
		rp.healthCheck = NewHealthChecker(cfg, rp.backends)
	}

	// Create HTTP server
	rp.server = &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      rp,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return rp, nil
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get next backend
	backend := rp.loadBalancer.NextBackend()
	if backend == nil {
		http.Error(w, "No healthy backends available", http.StatusServiceUnavailable)
		log.Printf("No healthy backends available for request: %s %s", r.Method, r.URL.Path)
		return
	}

	// Track connection
	backend.mu.Lock()
	backend.Connections++
	backend.mu.Unlock()

	defer func() {
		backend.mu.Lock()
		backend.Connections--
		backend.mu.Unlock()
	}()

	// Log request
	log.Printf("Proxying request: %s %s -> %s", r.Method, r.URL.Path, backend.URL.String())

	// Proxy the request
	backend.Proxy.ServeHTTP(w, r)
}

func (rp *ReverseProxy) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Proxy error: %v", err)
	http.Error(w, "Bad Gateway", http.StatusBadGateway)
}

func (rp *ReverseProxy) Start() error {
	// Start health checker
	if rp.healthCheck != nil {
		rp.healthCheck.Start()
	}

	return rp.server.ListenAndServe()
}

func (rp *ReverseProxy) Shutdown() error {
	// Stop health checker
	if rp.healthCheck != nil {
		rp.healthCheck.Stop()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return rp.server.Shutdown(ctx)
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

func (b *Backend) GetConnections() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Connections
}