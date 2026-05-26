package agent

import (
	"fmt"
	"os"
	"time"
)

// Config holds the configuration for the Docker agent.
type Config struct {
	// Host is the Docker daemon host to connect to.
	// Defaults to the value of DOCKER_HOST environment variable or unix:///var/run/docker.sock.
	Host string

	// TLSVerify enables TLS verification when connecting to the Docker daemon.
	TLSVerify bool

	// TLSCACert is the path to the CA certificate for TLS verification.
	TLSCACert string

	// TLSCert is the path to the client certificate for TLS authentication.
	TLSCert string

	// TLSKey is the path to the client key for TLS authentication.
	TLSKey string

	// PollInterval is the interval at which the agent polls the Docker daemon for events.
	PollInterval time.Duration

	// MaxRetries is the maximum number of retries when connecting to the Docker daemon.
	MaxRetries int

	// RetryDelay is the delay between retries when connecting to the Docker daemon.
	RetryDelay time.Duration

	// LogLevel sets the verbosity of agent logging (debug, info, warn, error).
	LogLevel string

	// Labels are additional key-value labels to attach to agent metadata.
	Labels map[string]string
}

// DefaultConfig returns a Config populated with sensible default values.
func DefaultConfig() *Config {
	host := os.Getenv("DOCKER_HOST")
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}

	return &Config{
		Host:         host,
		TLSVerify:    os.Getenv("DOCKER_TLS_VERIFY") == "1",
		TLSCACert:    os.Getenv("DOCKER_TLS_CACERT"),
		TLSCert:      os.Getenv("DOCKER_TLS_CERT"),
		TLSKey:       os.Getenv("DOCKER_TLS_KEY"),
		PollInterval: 30 * time.Second,
		MaxRetries:   5,
		RetryDelay:   2 * time.Second,
		LogLevel:     "info",
		Labels:       make(map[string]string),
	}
}

// Validate checks that the Config values are consistent and returns an error
// describing the first invalid field encountered.
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("config: host must not be empty")
	}

	if c.TLSVerify {
		if c.TLSCACert == "" {
			return fmt.Errorf("config: TLS CA cert path must be set when TLS verification is enabled")
		}
		if c.TLSCert == "" {
			return fmt.Errorf("config: TLS cert path must be set when TLS verification is enabled")
		}
		if c.TLSKey == "" {
			return fmt.Errorf("config: TLS key path must be set when TLS verification is enabled")
		}
	}

	if c.PollInterval <= 0 {
		return fmt.Errorf("config: poll interval must be a positive duration, got %s", c.PollInterval)
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("config: max retries must be non-negative, got %d", c.MaxRetries)
	}

	if c.RetryDelay < 0 {
		return fmt.Errorf("config: retry delay must be non-negative, got %s", c.RetryDelay)
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("config: invalid log level %q, must be one of debug, info, warn, error", c.LogLevel)
	}

	return nil
}
