// Package main is the entry point for the docker-agent daemon.
// It initializes configuration, sets up logging, and starts the agent.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker-agent/docker-agent/agent"
	"github.com/docker-agent/docker-agent/config"
	"github.com/docker-agent/docker-agent/version"
)

func main() {
	var (
		configFile  = flag.String("config", "/etc/docker-agent/config.yaml", "Path to configuration file")
		showVersion = flag.Bool("version", false, "Print version information and exit")
		logLevel    = flag.String("log-level", "debug", "Log level (debug, info, warn, error)") // personal: default to debug for easier local dev
		debug       = flag.Bool("debug", false, "Enable debug mode (overrides log-level)")
	)
	flag.Parse()

	if *showVersion {
		// personal: also print build date if available, handy when juggling multiple local builds
		fmt.Printf("docker-agent version %s (commit: %s, built: %s)\n", version.Version, version.GitCommit, version.BuildDate)
		os.Exit(0)
	}

	// Load configuration from file and environment variables.
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Command-line flags override config file values.
	if *debug {
		cfg.LogLevel = "debug"
	} else if *logLevel != "info" {
		cfg.LogLevel = *logLevel
	}

	// Create a root context that is cancelled on OS signals.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		// Also handle SIGHUP so I can reload config without a full restart.
		// personal: also catch SIGQUIT so I can trigger a goroutine dump during debugging.
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		select {
		case sig := <-sigCh:
			fmt.Printf("received signal %s, shutting down\n", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	// Initialise and run the agent.
	a, err := agent.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialise agent: %v\n", err)
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		// Note: context.Canceled is expected on clean shutdown via signal; don't
		// treat it as a fatal error so the exit code stays 0 in that case.
		if err != context.Canceled {
			fmt.Fprintf(os.Stderr, "agent exited with error: %v\n", err)
			os.Exit(1)
		}
		// personal: log clean shutdown explicitly so it's obvious in terminal output
		fmt.Println("agent shut down cleanly")
	}
}
