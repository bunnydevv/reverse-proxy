package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Server       ServerConfig       `yaml:"server"`
	Backends     []Backend          `yaml:"backends"`
	LoadBalancer LoadBalancerConfig `yaml:"load_balancer"`
	HealthCheck  HealthCheckConfig  `yaml:"health_check"`
	Logging      LoggingConfig      `yaml:"logging"`
	TLS          *TLSConfig         `yaml:"tls,omitempty"`
	Limits       LimitsConfig       `yaml:"limits"`
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// Backend represents a backend server configuration
type Backend struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
}

// LoadBalancerConfig contains load balancing algorithm configuration
type LoadBalancerConfig struct {
	Algorithm string `yaml:"algorithm"` // round-robin, least-connections, weighted
}

// HealthCheckConfig contains health check configuration
type HealthCheckConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Path     string        `yaml:"path"`
}

// LoggingConfig contains logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`  // debug, info, warn, error
	Format string `yaml:"format"` // json, text
}

// TLSConfig contains TLS/HTTPS configuration
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// LimitsConfig contains connection and request limits
type LimitsConfig struct {
	MaxConnections    int           `yaml:"max_connections"`
	MaxIdleConns      int           `yaml:"max_idle_conns"`
	MaxConnsPerHost   int           `yaml:"max_conns_per_host"`
	RequestTimeout    time.Duration `yaml:"request_timeout"`
	MaxRequestBodySize int64        `yaml:"max_request_body_size"`
}

// Load reads and parses the configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	setDefaults(&cfg)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

func setDefaults(cfg *Config) {
	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 10 * time.Second
	}
	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 10 * time.Second
	}
	if cfg.Server.IdleTimeout == 0 {
		cfg.Server.IdleTimeout = 120 * time.Second
	}
	if cfg.LoadBalancer.Algorithm == "" {
		cfg.LoadBalancer.Algorithm = "round-robin"
	}
	if cfg.HealthCheck.Interval == 0 {
		cfg.HealthCheck.Interval = 10 * time.Second
	}
	if cfg.HealthCheck.Timeout == 0 {
		cfg.HealthCheck.Timeout = 5 * time.Second
	}
	if cfg.HealthCheck.Path == "" {
		cfg.HealthCheck.Path = "/health"
	}
	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}
	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "text"
	}
	if cfg.Limits.MaxConnections == 0 {
		cfg.Limits.MaxConnections = 10000
	}
	if cfg.Limits.MaxIdleConns == 0 {
		cfg.Limits.MaxIdleConns = 100
	}
	if cfg.Limits.MaxConnsPerHost == 0 {
		cfg.Limits.MaxConnsPerHost = 100
	}
	if cfg.Limits.RequestTimeout == 0 {
		cfg.Limits.RequestTimeout = 30 * time.Second
	}
	if cfg.Limits.MaxRequestBodySize == 0 {
		cfg.Limits.MaxRequestBodySize = 10 * 1024 * 1024 // 10MB
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate server address
	if c.Server.Address == "" {
		return fmt.Errorf("server address is required")
	}

	// Validate backends
	if len(c.Backends) == 0 {
		return fmt.Errorf("at least one backend is required")
	}

	for i, backend := range c.Backends {
		if backend.URL == "" {
			return fmt.Errorf("backend %d: URL is required", i)
		}
		
		// Validate URL format
		_, err := url.Parse(backend.URL)
		if err != nil {
			return fmt.Errorf("backend %d: invalid URL %s: %w", i, backend.URL, err)
		}

		// Validate weight
		if backend.Weight < 0 {
			return fmt.Errorf("backend %d: weight must be non-negative", i)
		}
	}

	// Validate load balancer algorithm
	validAlgorithms := map[string]bool{
		"round-robin":      true,
		"least-connections": true,
		"weighted":         true,
	}
	if !validAlgorithms[c.LoadBalancer.Algorithm] {
		return fmt.Errorf("invalid load balancer algorithm: %s (must be one of: round-robin, least-connections, weighted)", c.LoadBalancer.Algorithm)
	}

	// Validate timeouts
	if c.Server.ReadTimeout < 0 {
		return fmt.Errorf("server read_timeout must be non-negative")
	}
	if c.Server.WriteTimeout < 0 {
		return fmt.Errorf("server write_timeout must be non-negative")
	}
	if c.Server.IdleTimeout < 0 {
		return fmt.Errorf("server idle_timeout must be non-negative")
	}
	if c.HealthCheck.Enabled && c.HealthCheck.Interval < 0 {
		return fmt.Errorf("health_check interval must be non-negative")
	}
	if c.HealthCheck.Enabled && c.HealthCheck.Timeout < 0 {
		return fmt.Errorf("health_check timeout must be non-negative")
	}

	// Validate logging
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[strings.ToLower(c.Logging.Level)] {
		return fmt.Errorf("invalid logging level: %s (must be one of: debug, info, warn, error)", c.Logging.Level)
	}

	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	if !validFormats[strings.ToLower(c.Logging.Format)] {
		return fmt.Errorf("invalid logging format: %s (must be one of: json, text)", c.Logging.Format)
	}

	// Validate TLS configuration
	if c.TLS != nil && c.TLS.Enabled {
		if c.TLS.CertFile == "" {
			return fmt.Errorf("TLS cert_file is required when TLS is enabled")
		}
		if c.TLS.KeyFile == "" {
			return fmt.Errorf("TLS key_file is required when TLS is enabled")
		}
	}

	// Validate limits
	if c.Limits.MaxConnections < 0 {
		return fmt.Errorf("max_connections must be non-negative")
	}
	if c.Limits.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns must be non-negative")
	}
	if c.Limits.MaxConnsPerHost < 0 {
		return fmt.Errorf("max_conns_per_host must be non-negative")
	}
	if c.Limits.RequestTimeout < 0 {
		return fmt.Errorf("request_timeout must be non-negative")
	}
	if c.Limits.MaxRequestBodySize < 0 {
		return fmt.Errorf("max_request_body_size must be non-negative")
	}

	return nil
}