// Package agent provides the core Docker agent functionality,
// including container management, event handling, and communication
// with the Docker daemon.
package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// Config holds the configuration for the Docker agent.
type Config struct {
	// DockerHost is the address of the Docker daemon socket or TCP endpoint.
	DockerHost string

	// PollInterval is how often the agent polls for container state changes.
	PollInterval time.Duration

	// Labels are key-value pairs used to filter which containers this agent manages.
	Labels map[string]string

	// LogLevel sets the verbosity of agent logging.
	LogLevel string
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DockerHost:   client.DefaultDockerHost,
		PollInterval: 30 * time.Second,
		Labels:       map[string]string{},
		LogLevel:     "info",
	}
}

// Agent manages communication with the Docker daemon and orchestrates
// container lifecycle operations on behalf of the host system.
type Agent struct {
	cfg    Config
	client *client.Client
	log    *logrus.Logger

	mu      sync.Mutex
	running bool
	stopCh  chan struct{}
}

// New creates a new Agent with the provided configuration.
// It initialises a Docker client but does not yet connect to the daemon.
func New(cfg Config) (*Agent, error) {
	log := logrus.New()

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", cfg.LogLevel, err)
	}
	log.SetLevel(level)

	dockerClient, err := client.NewClientWithOpts(
		client.WithHost(cfg.DockerHost),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating docker client: %w", err)
	}

	return &Agent{
		cfg:    cfg,
		client: dockerClient,
		log:    log,
		stopCh: make(chan struct{}),
	}, nil
}

// Start begins the agent's main event loop. It blocks until the provided
// context is cancelled or Stop is called.
func (a *Agent) Start(ctx context.Context) error {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		return fmt.Errorf("agent is already running")
	}
	a.running = true
	a.mu.Unlock()

	a.log.WithField("host", a.cfg.DockerHost).Info("docker agent starting")

	// Verify connectivity before entering the loop.
	if _, err := a.client.Ping(ctx); err != nil {
		return fmt.Errorf("connecting to docker daemon: %w", err)
	}
	a.log.Info("connected to docker daemon")

	ticker := time.NewTicker(a.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.log.Info("context cancelled, agent shutting down")
			return ctx.Err()
		case <-a.stopCh:
			a.log.Info("stop signal received, agent shutting down")
			return nil
		case <-ticker.C:
			if err := a.reconcile(ctx); err != nil {
				a.log.WithError(err).Warn("reconciliation error")
			}
		}
	}
}

// Stop signals the agent to cease operation gracefully.
func (a *Agent) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.running {
		close(a.stopCh)
		a.running = false
	}
}

// reconcile inspects the current container state and performs any
// corrective actions required to reach the desired state.
func (a *Agent) reconcile(ctx context.Context) error {
	containers, err := a.client.ContainerList(ctx, containertypes.ListOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("listing containers: %w", err)
	}

	a.log.WithField("count", len(containers)).Debug("reconcile: container snapshot")
	return nil
}
