package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRiskBadID(t *testing.T) {
	setupRiskStorage()

	code, body, err := request(http.MethodPost, "/risk", map[string]interface{}{
		"risk-id": 1,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, IMMUTABLE_RISK, code, body)
}

func TestRiskBadTime(t *testing.T) {
	setupRiskStorage()

	code, body, err := request(http.MethodPost, "/risk", map[string]interface{}{
		"voyage-id": 1,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, RISK_TIME_REQUIRED, code, body)
}

func TestRiskOK(t *testing.T) {
	setupRiskStorage()

	code, body, err := request(http.MethodPost, "/risk", map[string]interface{}{
		"voyage-id": 1,
		"time":      "2022-01-01T01:00:00.02Z",
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	equalJSON(t, map[string]interface{}{
		"risk-id": float64(1),
	}, body)

	code, body, err = request(http.MethodPost, "/risk", map[string]interface{}{
		"voyage-id":  1,
		"time":       "2022-01-01T01:01:00.02Z",
		"management": 1,
		"crew":       2,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	equalJSON(t, map[string]interface{}{
		"risk-id": float64(2),
	}, body)
}
