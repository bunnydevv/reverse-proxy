package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bunnydevv/reverse-proxy/config"
	"github.com/bunnydevv/reverse-proxy/proxy"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start the reverse proxy
	rp, err := proxy.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create reverse proxy: %v", err)
	}

	// Start the proxy server
	go func() {
		log.Printf("Starting reverse proxy on %s", cfg.Server.Address)
		if err := rp.Start(); err != nil {
			log.Fatalf("Failed to start reverse proxy: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down reverse proxy...")
	if err := rp.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown reverse proxy: %v", err)
	}
	log.Println("Reverse proxy stopped")
}