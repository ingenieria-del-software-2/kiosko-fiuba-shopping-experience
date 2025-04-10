package server

import "time"

// Server constants
const (
	// ShutdownTimeout is the maximum wait time for server shutdown
	ShutdownTimeout = 15 * time.Second

	// DefaultServerPort is the default port for the server
	DefaultServerPort = 8001

	// ReadTimeout is the maximum duration for reading the entire request
	ReadTimeout = 10 * time.Second

	// WriteTimeout is the maximum duration before timing out writes of the response
	WriteTimeout = 10 * time.Second

	// IdleTimeout is the maximum amount of time to wait for the next request
	IdleTimeout = 120 * time.Second
)
