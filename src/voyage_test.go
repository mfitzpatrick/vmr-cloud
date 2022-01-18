package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVoyageOK(t *testing.T) {
	code, body, err := request(http.MethodPost, "/voyage", map[string]interface{}{
		"vessel-id": 1,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, true, equalJSON(t, map[string]interface{}{
		"voyage-id": float64(1), //NB: integers must be float64 because that is the default for the JSON lib
	}, body))
}

func TestNewVoyageBadID(t *testing.T) {
	code, body, err := request(http.MethodPost, "/voyage", map[string]interface{}{
		"vessel-id": 0,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, equalHTTPError(t, INVALID_VESSEL_ID, code, body))
}

func TestNewVoyageBadMethod(t *testing.T) {
	code, body, err := request(http.MethodPut, "/voyage", map[string]interface{}{
		"vessel-id": 0,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, ENDPOINT_NOT_FOUND, code, body)
}

func TestNewVoyageBadURI(t *testing.T) {
	code, body, err := request(http.MethodGet, "/bogus-path", map[string]interface{}{
		"vessel-id": 0,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, ENDPOINT_NOT_FOUND, code, body)
}

func TestGetVoyageBadID(t *testing.T) {
	code, body, err := request(http.MethodGet, "/voyage", map[string]interface{}{
		"voyage-id": 0,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, INVALID_VOYAGE_ID, code, body)
}
