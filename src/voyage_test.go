package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVoyageOK(t *testing.T) {
	setupVoyageStorage()

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
	assert.Equal(t, true, equalHTTPError(t, INVALID_VOYAGE_ID, code, body))
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

func testVoyageStoreAndRetrieve(t *testing.T, expect, set map[string]interface{}) {
	code, body, err := request(http.MethodPost, "/voyage", set)
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, true, equalJSON(t, map[string]interface{}{
		"voyage-id": expect["voyage-id"],
	}, body))
	var returnedMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &returnedMap); err != nil {
		assert.Equal(t, nil, err)
	} else {
		assert.Equal(t, expect["voyage-id"], returnedMap["voyage-id"])
	}

	code, body, err = request(http.MethodGet, "/voyage", map[string]interface{}{
		"voyage-id": expect["voyage-id"],
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	equalVoyageMap(t, expect, body)
}

func TestGetVoyageStoreAndRetrieve(t *testing.T) {
	setupVoyageStorage()

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":   float64(1),
		"vessel-id":   4,
		"start-hours": 101,
		"title":       "Breakdown Coomera",
	}, map[string]interface{}{
		"vessel-id":   4,
		"start-hours": 101,
		"title":       "Breakdown Coomera",
	})

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":   float64(2),
		"vessel-id":   2,
		"start-hours": 102,
		"title":       "Jump Start Tipplers",
	}, map[string]interface{}{
		"vessel-id":   2,
		"start-hours": 102,
		"title":       "Jump Start Tipplers",
	})

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":   float64(2),
		"vessel-id":   2,
		"start-hours": 102,
		"title":       "Jump Start Tipplers",
		"description": "Jump start jet ski on beach at tipplers. No incidents",
	}, map[string]interface{}{
		"voyage-id":   float64(2),
		"description": "Jump start jet ski on beach at tipplers. No incidents",
	})

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":   float64(3),
		"vessel-id":   2,
		"start-hours": 103,
		"title":       "Search Offshore",
		"skipper": map[string]interface{}{
			"name": "Gerry Hatrick",
			"rank": "Offshore Skipper",
		},
	}, map[string]interface{}{
		"vessel-id":   2,
		"start-hours": 103,
		"title":       "Search Offshore",
		"skipper": map[string]interface{}{
			"name": "Gerry Hatrick",
			"rank": "Offshore Skipper",
		},
	})

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":   float64(2),
		"vessel-id":   2,
		"start-hours": 102,
		"title":       "Jump Start Tipplers",
		"description": "Jump start jet ski on beach at tipplers. No incidents",
		"weather": map[string]interface{}{
			"wind": map[string]interface{}{
				"speed-knots":       11,
				"direction-degrees": 050,
			},
		},
	}, map[string]interface{}{
		"voyage-id": float64(2),
		"weather": map[string]interface{}{
			"wind": map[string]interface{}{
				"speed-knots":       11,
				"direction-degrees": 050,
			},
		},
	})
}
