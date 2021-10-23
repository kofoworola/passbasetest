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
	input := float32(300.0)
	out, err := handler.Convert("EUR", "EUR", input)
	if err != nil {
		t.Fatal(err)
	}

	if out != input {
		t.Fatalf("expected %f got %f", input, out)
	}
}

// to test the error output in case of error
func TestInvalidApiKey(t *testing.T) {
	handler := New(Config{
		ApiKey: "dummy-key",
	})
	input := float32(300.0)
	_, err := handler.Convert("USD", "USD", input)
	if err == nil || !strings.Contains(err.Error(), "You have not supplied a valid API Access Key") {
		t.Fatalf("expected invalid api key error got %v", err)
	}
}
