package proxy

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/bunnydevv/reverse-proxy/config"
)

type HealthChecker struct {
	config   *config.Config
	backends []*Backend
	client   *http.Client
	stop     chan struct{}
}

func NewHealthChecker(cfg *config.Config, backends []*Backend) *HealthChecker {
	return &HealthChecker{
		config:   cfg,
		backends: backends,
		client: &http.Client{
			Timeout: cfg.HealthCheck.Timeout,
		},
		stop: make(chan struct{}),
	}
}

func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.config.HealthCheck.Interval)
	go func() {
		// Do initial health check
		hc.checkAll()

		for {
			select {
			case <-ticker.C:
				hc.checkAll()
			case <-hc.stop:
					ticker.Stop()
					return
			}
		}
	}()
}

func (hc *HealthChecker) Stop() {
	close(hc.stop)
}

func (hc *HealthChecker) checkAll() {
	for _, backend := range hc.backends {
		go hc.check(backend)
	}
}

func (hc *HealthChecker) check(backend *Backend) {
	url := backend.URL.String() + hc.config.HealthCheck.Path
	ctx, cancel := context.WithTimeout(context.Background(), hc.config.HealthCheck.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Health check failed for %s: %v", backend.URL.String(), err)
		backend.SetAlive(false)
		return
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		log.Printf("Health check failed for %s: %v", backend.URL.String(), err)
		backend.SetAlive(false)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		if !backend.IsAlive() {
			log.Printf("Backend %s is now healthy", backend.URL.String())
		}
		backend.SetAlive(true)
	} else {
		log.Printf("Health check failed for %s: status code %d", backend.URL.String(), resp.StatusCode)
		backend.SetAlive(false)
	}
}