package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAssistOK(t *testing.T) {
	code, body, err := request(http.MethodPost, "/assist", map[string]interface{}{
		"voyage-id": 1,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, true, equalJSON(t, map[string]interface{}{
		"assist-id": float64(1), //NB: integers must be float64 because that is the default for the JSON lib
	}, body))
}

func TestNewAssistBadID(t *testing.T) {
	code, body, err := request(http.MethodPost, "/assist", map[string]interface{}{
		"voyage-id": 0,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, true, equalHTTPError(t, INVALID_VOYAGE_ID, code, body))
}

func TestNewAssistBadMethod(t *testing.T) {
	code, body, err := request(http.MethodPut, "/assist", map[string]interface{}{
		"voyage-id": 0,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, ENDPOINT_NOT_FOUND, code, body)
}

func TestNewAssistBadURI(t *testing.T) {
	code, body, err := request(http.MethodGet, "/bogus-path", map[string]interface{}{
		"voyage-id": 0,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, ENDPOINT_NOT_FOUND, code, body)
}

func TestGetAssistBadID(t *testing.T) {
	code, body, err := request(http.MethodGet, "/assist", map[string]interface{}{
		"assist-id": 0,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, INVALID_VOYAGE_ID, code, body)
}
