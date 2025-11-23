package proxy

import (
	"sync"
	"sync/atomic"
)

type LoadBalancer interface {
	NextBackend() *Backend
}

// Round Robin Load Balancer
type RoundRobinBalancer struct {
	backends []*Backend
	current  uint32
}

func NewRoundRobinBalancer(backends []*Backend) *RoundRobinBalancer {
	return &RoundRobinBalancer{
		backends: backends,
		current:  0,
	}
}

func (rb *RoundRobinBalancer) NextBackend() *Backend {
	n := len(rb.backends)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		idx := atomic.AddUint32(&rb.current, 1) % uint32(n)
		backend := rb.backends[idx]
		if backend.IsAlive() {
			return backend
		}
	}

	return nil
}

// Least Connections Load Balancer
type LeastConnectionsBalancer struct {
	backends []*Backend
	mu       sync.RWMutex
}

func NewLeastConnectionsBalancer(backends []*Backend) *LeastConnectionsBalancer {
	return &LeastConnectionsBalancer{
		backends: backends,
	}
}

func (lb *LeastConnectionsBalancer) NextBackend() *Backend {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	var selected *Backend
	minConnections := -1

	for _, backend := range lb.backends {
		if !backend.IsAlive() {
			continue
		}

		connections := backend.GetConnections()
		if minConnections == -1 || connections < minConnections {
			minConnections = connections
			selected = backend
		}
	}

	return selected
}

// Weighted Load Balancer
type WeightedBalancer struct {
	backends []*Backend
	current  uint32
}

func NewWeightedBalancer(backends []*Backend) *WeightedBalancer {
	return &WeightedBalancer{
		backends: backends,
		current:  0,
	}
}

func (wb *WeightedBalancer) NextBackend() *Backend {
	var expandedBackends []*Backend
	for _, backend := range wb.backends {
		for i := 0; i < backend.Weight; i++ {
			expandedBackends = append(expandedBackends, backend)
		}
	}

	n := len(expandedBackends)
	if n == 0 {
		return nil
	}

	for i := 0; i < n; i++ {
		idx := atomic.AddUint32(&wb.current, 1) % uint32(n)
		backend := expandedBackends[idx]
		if backend.IsAlive() {
			return backend
		}
	}

	return nil
}