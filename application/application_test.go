//go:build integration
// +build integration

package application

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	file, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	file.WriteString(`{
		"app_name": "TestApp",
		"app_mode": "debug",
		"port": "8080",
		"read_timeout": 30,
		"write_timeout": 30,
		"max_idle_time_out": 30,
		"server_timeout": 30,
		"db_config": {
			"host": "db",
			"port": 5432,
			"user": "transact_api",
			"dbname": "transact_api"
		}
	}`)
	file.Close()

	go Run(file.Name())

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080/api/transact/v1/healthz")
	if err != nil {
		t.Fatalf("Failed to send request to server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got '%s'", resp.Status)
	}

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}
	p.Signal(os.Interrupt)

	time.Sleep(1 * time.Second)
}
