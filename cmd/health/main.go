package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arctoslabs-io/chain-health/internal/config"
	"github.com/arctoslabs-io/chain-health/internal/monitor"
	"github.com/arctoslabs-io/chain-health/internal/server"
)

var version = "dev"

func main() {
	configPath := flag.String("config", "config.yaml", "path to configuration file")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("chain-health %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mon, err := monitor.New(cfg)
	if err != nil {
		log.Fatalf("failed to create monitor: %v", err)
	}

	if cfg.Monitoring.Prometheus.Enabled {
		go server.StartMetrics(cfg.Monitoring.Prometheus.Port)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go mon.Run()
	log.Printf("chain-health %s started, monitoring %d chains", version, len(cfg.Chains))

	<-stop
	log.Println("shutting down")
	mon.Stop()
}
