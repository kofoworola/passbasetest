package fixer

import (
	"strings"
	"testing"

	"github.com/kelseyhightower/envconfig"
)

func TestFixerIntegration(t *testing.T) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		t.Fatal(err)
	}

	handler := New(config)
	rate, err := handler.Convert("EUR", "EUR")
	if err != nil {
		t.Fatal(err)
	}

	if rate != 1 {
		t.Fatalf("expected 1 got %f", rate)
	}
}

// to test the error output in case of error
func TestInvalidApiKey(t *testing.T) {
	handler := New(Config{
		ApiKey: "dummy-key",
	})
	_, err := handler.Convert("USD", "USD")
	if err == nil || !strings.Contains(err.Error(), "You have not supplied a valid API Access Key") {
		t.Fatalf("expected invalid api key error got %v", err)
	}
}
