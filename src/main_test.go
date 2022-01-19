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

func request(method, uri string, body map[string]interface{}) (int, string, error) {
	httpRecorder := httptest.NewRecorder()
	var httpRequest *http.Request
	var bodyString string
	if b, err := json.Marshal(body); err != nil {
		return 0, "", errors.Wrapf(err, "HTTP request tester helper")
	} else {
		bodyString = string(b)
	}
	if req, err := http.NewRequest(method, HTTP_URI_BASE+uri, strings.NewReader(bodyString)); err != nil {
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

func equalVoyageMap(t *testing.T, expect map[string]interface{}, body string) bool {
	expectVoyage := voyage{}
	actualVoyage := voyage{}
	if expectString, err := json.Marshal(expect); err != nil {
		t.Errorf("JSON marshal: %v", err)
		return false
	} else if err := json.Unmarshal(expectString, &expectVoyage); err != nil {
		t.Errorf("JSON unmarshal: %v", err)
		return false
	} else if err := json.Unmarshal([]byte(body), &actualVoyage); err != nil {
		t.Errorf("JSON unmarshal: %v", err)
		return false
	} else {
		return assert.Equal(t, expectVoyage, actualVoyage)
	}
}

func equalAssistMap(t *testing.T, expect map[string]interface{}, body string) bool {
	expectAssist := assist{}
	actualAssist := assist{}
	if expectString, err := json.Marshal(expect); err != nil {
		t.Errorf("JSON marshal: %v", err)
		return false
	} else if err := json.Unmarshal(expectString, &expectAssist); err != nil {
		t.Errorf("JSON unmarshal: %v", err)
		return false
	} else if err := json.Unmarshal([]byte(body), &actualAssist); err != nil {
		t.Errorf("JSON unmarshal: %v", err)
		return false
	} else {
		return assert.Equal(t, expectAssist, actualAssist)
	}
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
