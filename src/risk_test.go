package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRiskBadID(t *testing.T) {
	code, body, err := request(http.MethodPost, "/risk", map[string]interface{}{
		"risk-id": 1,
	})
	assert.Equal(t, nil, err)
	equalHTTPError(t, INVALID_RISK_ID, code, body)
}
