package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func request(method, uri, body string) (int, string, error) {
	httpRecorder := httptest.NewRecorder()
	var httpRequest *http.Request
	if req, err := http.NewRequest(http.MethodPost, HTTP_URI_BASE+uri, strings.NewReader(body)); err != nil {
		return 0, "", errors.Wrapf(err, "HTTP request tester helper")
	} else {
		httpRequest = req
	}
	httpHandler(httpRecorder, httpRequest)
	resp := httpRecorder.Result()
	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		return resp.StatusCode, "", errors.Wrapf(err, "HTTP request tester helper")
	} else {
		return resp.StatusCode, string(b), nil
	}
}

func equalJSON(t *testing.T, cmp map[string]interface{}, body string) bool {
	input := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &input); err != nil {
		t.Errorf("JSON parse failed: %v", err)
		return false
	}
	return assert.Equal(t, cmp, input)
}

func equalHTTPError(t *testing.T, expect httpError, code int, body string) bool {
	var herror httpError
	if err := json.Unmarshal([]byte(body), &herror); err != nil {
		t.Errorf("JSON parse failed: %v", err)
		return false
	}
	if !assert.Equal(t, expect.ID, herror.ID) {
		return false
	}
	if !assert.Equal(t, expect.statusCode, code) {
		return false
	}
	return true
}

func TestNewVoyageOK(t *testing.T) {
	code, body, err := request(http.MethodPost, "/voyage", "{\"vessel-id\":1}")
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, true, equalJSON(t, map[string]interface{}{
		"voyage-id": float64(1),
	}, body))
}

func TestNewVoyageBadID(t *testing.T) {
	code, body, err := request(http.MethodPost, "/voyage", "{\"vessel-id\":0}")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, equalHTTPError(t, INVALID_VESSEL_ID, code, body))
}
