package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/Jagreen1970/battleship/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	cfg := &app.Config{
		Server: app.ServerConfig{
			Port:    8081,
			Timeout: 5 * time.Second,
		},
	}

	srv := New(cfg.Server)
	assert.NotNil(t, srv)

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Make a test request
	resp, err := http.Get("http://localhost:8081")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	resp.Body.Close()

	// Shutdown the server
	err = srv.Shutdown()
	assert.NoError(t, err)

	// Check if server stopped with expected error
	err = <-errCh
	assert.ErrorIs(t, err, http.ErrServerClosed)
}
