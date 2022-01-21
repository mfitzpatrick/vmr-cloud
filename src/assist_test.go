package main

import (
	"encoding/json"
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
	equalHTTPError(t, INVALID_ASSIST_ID, code, body)
}

func testAssistStoreAndRetrieve(t *testing.T, expect, set map[string]interface{}) {
	code, body, err := request(http.MethodPost, "/assist", set)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, true, equalJSON(t, map[string]interface{}{
		"assist-id": expect["assist-id"],
	}, body))
	var returnedMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &returnedMap); err != nil {
		assert.Equal(t, nil, err)
	} else {
		assert.Equal(t, expect["assist-id"], returnedMap["assist-id"])
	}

	code, body, err = request(http.MethodGet, "/assist", map[string]interface{}{
		"assist-id": expect["assist-id"],
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	equalAssistMap(t, expect, body)
}

func TestGetAssistStoreAndRetrieve(t *testing.T) {
	setupAssistStorage()

	testAssistStoreAndRetrieve(t, map[string]interface{}{
		"assist-id": float64(1),
		"voyage-id": 4,
		"problem":   "Breakdown",
	}, map[string]interface{}{
		"voyage-id": 4,
		"problem":   "Breakdown",
	})

	testAssistStoreAndRetrieve(t, map[string]interface{}{
		"assist-id": float64(2),
		"voyage-id": 2,
		"problem":   "Jump Start",
	}, map[string]interface{}{
		"voyage-id": 2,
		"problem":   "Jump Start",
	})

	testAssistStoreAndRetrieve(t, map[string]interface{}{
		"assist-id": float64(2),
		"voyage-id": 2,
		"problem":   "Jump Start",
		"action":    "Jump Start",
	}, map[string]interface{}{
		"assist-id": float64(2),
		"action":    "Jump Start",
	})

	testAssistStoreAndRetrieve(t, map[string]interface{}{
		"assist-id": float64(3),
		"voyage-id": 2,
		"client": map[string]interface{}{
			"name":          "Persephone",
			"phone":         "123456",
			"member-number": 0,
		},
	}, map[string]interface{}{
		"voyage-id": 2,
		"client": map[string]interface{}{
			"name":          "Persephone",
			"phone":         "123456",
			"member-number": 0,
		},
	})

	testAssistStoreAndRetrieve(t, map[string]interface{}{
		"assist-id": float64(2),
		"voyage-id": 2,
		"problem":   "Jump Start",
		"action":    "Jump Start",
		"pickup": map[string]interface{}{
			"location": map[string]interface{}{
				"gps": map[string]interface{}{
					"lat":  -27.4748923,
					"long": 153.4839283,
				},
			},
			"time": "2022-01-01T01:01:18Z",
		},
		"destination": map[string]interface{}{
			"time": "0001-01-01T00:00:00Z",
		},
	}, map[string]interface{}{
		"assist-id": float64(2),
		"pickup": map[string]interface{}{
			"location": map[string]interface{}{
				"gps": map[string]interface{}{
					"lat":  -27.4748923,
					"long": 153.4839283,
				},
			},
			"time": "2022-01-01T01:01:18Z",
		},
	})
}
