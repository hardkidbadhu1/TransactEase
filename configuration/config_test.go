package configuration

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	file, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	os.Setenv("MYSQL_PASSWORD", "testpassword")
	file.WriteString(`{
		"app_name": "TestApp",
		"port": "8080",
		"read_timeout": 30,
		"write_timeout": 30,
		"max_idle_time_out": 30,
		"server_timeout": 30,
		"db_config": {
		"host": "localhost",
		"port": 3306,
		"user": "testuser",
		"dbname": "testdb"
		}
	}`)
	file.Close()

	cfg, err := Parse(file.Name())
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if cfg.GetAppName() != "TestApp" {
		t.Errorf("Expected AppName to be 'TestApp', got '%s'", cfg.GetAppName())
	}
	if cfg.GetPort() != "8080" {
		t.Errorf("Expected Port to be 8080, got '%s'", cfg.GetPort())
	}

	if cfg.GetDBHost() != "localhost" {
		t.Errorf("Expected DBHost to be 'localhost', got '%s'", cfg.GetDBHost())
	}
	if cfg.GetDBPort() != 3306 {
		t.Errorf("Expected DBPort to be 3306, got '%d'", cfg.GetDBPort())
	}

	if cfg.GetWriteTimeout() != 30 {
		t.Errorf("Expected WriteTimeout to be 30, got '%d'", cfg.GetWriteTimeout())
	}

	if cfg.GetMaxIdleTimeOut() != 30 {
		t.Errorf("Expected MaxIdleTimeOut to be 30, got '%d'", cfg.GetMaxIdleTimeOut())
	}

	if cfg.GetServerTimeOut() != 30 {
		t.Errorf("Expected ServerTimeOut to be 30, got '%d'", cfg.GetServerTimeOut())
	}

	if cfg.GetDBUser() != "testuser" {
		t.Errorf("Expected DBUser to be 'testuser', got '%s'", cfg.GetDBUser())
	}

	if cfg.GetDBPassword() != "testpassword" {
		t.Errorf("Expected DBPassword to be 'testpassword', got '%s'", cfg.GetDBPassword())
	}

	if cfg.GetReadTimeout() != 30 {
		t.Errorf("Expected ReadTimeout to be 30, got '%d'", cfg.GetReadTimeout())
	}
}

func TestParseFileError(t *testing.T) {
	_, err := Parse("non_existent_file.json")
	// Check if the error was logged
	if err == nil {
		t.Errorf("Expected an error to be logged for non-existent file, got '%v'", err)
	}

	t.Log("err", err)

}

func TestParseUnmarshalError(t *testing.T) {
	file, err := os.CreateTemp("", "config")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	file.WriteString(`{
		"app_name": "TestApp",
		"port": "invalid_port", // This should be an integer
	}`)
	file.Close()

	_, err = Parse(file.Name())
	if err == nil {
		t.Errorf("Expected an error to be returned for invalid JSON, got '%v'", err)
	}
}
