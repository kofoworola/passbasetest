package fixer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	ApiKey string `required:"true"`
}

type Handler struct {
	apiKey string
}

func New(cfg Config) *Handler {
	return &Handler{cfg.ApiKey}
}

type responsePayload struct {
	Success bool `json:"success"`
	Error   struct {
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
	Rates map[string]float32 `json:"rates"`
}

// Base is removed due to api limit
const convertEndpoint = "http://data.fixer.io/api/latest?access_key=%s&symbols=%s"

func (h *Handler) Convert(from, to string, amount float32) (float32, error) {
	var converted float32

	parsedEndpoint := fmt.Sprintf(
		convertEndpoint,
		h.apiKey,
		//		from,
		to,
	)

	req, err := http.NewRequest(http.MethodGet, parsedEndpoint, nil)
	if err != nil {
		return converted, err
	}
	client := http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return converted, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusInternalServerError {
		return converted, fmt.Errorf("error with fixed api")
	}

	var payload responsePayload
	payload.Success = true
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return converted, fmt.Errorf("error reading body")
	}
	if !payload.Success {
		return converted, fmt.Errorf("%s:%s", payload.Error.Type, payload.Error.Info)
	}

	rate, ok := payload.Rates[to]
	if !ok {
		return converted, fmt.Errorf("invalid api response")
	}

	converted = rate * amount

	return converted, nil
}