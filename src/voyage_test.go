package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"testing"
	"time"

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
	setupRiskStorage()

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":    float64(1),
		"vessel-id":    4,
		"start-hours":  101,
		"title":        "Breakdown Coomera",
		"risk-history": []risk{},
	}, map[string]interface{}{
		"vessel-id":   4,
		"start-hours": 101,
		"title":       "Breakdown Coomera",
	})

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":    float64(2),
		"vessel-id":    2,
		"start-hours":  102,
		"title":        "Jump Start Tipplers",
		"risk-history": []risk{},
	}, map[string]interface{}{
		"vessel-id":   2,
		"start-hours": 102,
		"title":       "Jump Start Tipplers",
	})

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":    float64(2),
		"vessel-id":    2,
		"start-hours":  102,
		"title":        "Jump Start Tipplers",
		"description":  "Jump start jet ski on beach at tipplers. No incidents",
		"risk-history": []risk{},
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
		"risk-history": []risk{},
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
		"risk-history": []risk{},
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

func TestVoyageAndRisk(t *testing.T) {
	setupVoyageStorage()
	setupRiskStorage()
	staticTime := func(t *testing.T, timeString string) time.Time {
		tm, err := time.Parse(time.RFC3339, timeString)
		assert.Equal(t, nil, err)
		return tm
	}
	riskList := []risk{{
		VoyageID: 1,
		Mgmt:     1,
		Crew:     1,
		Time:     staticTime(t, "2022-01-01T01:00:00.00Z"),
	}, {
		VoyageID: 1,
		Mgmt:     2,
		Crew:     2,
		Time:     staticTime(t, "2022-01-01T02:00:00.00Z"),
	}, {
		VoyageID: 1,
		Mgmt:     3,
		Crew:     3,
		Time:     staticTime(t, "2022-01-01T03:00:00.00Z"),
	}}

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":    float64(1),
		"vessel-id":    4,
		"start-hours":  101,
		"title":        "Breakdown Coomera",
		"risk-history": []risk{},
	}, map[string]interface{}{
		"vessel-id":   4,
		"start-hours": 101,
		"title":       "Breakdown Coomera",
	})

	// Add risk entries for voyage
	for i, v := range riskList {
		riskJSON, err := json.Marshal(v)
		assert.Equal(t, nil, err)
		var riskMap map[string]interface{}
		err = json.Unmarshal(riskJSON, &riskMap)
		assert.Equal(t, nil, err)
		code, body, err := request(http.MethodPost, "/risk", riskMap)
		assert.Equal(t, nil, err)
		assert.Equal(t, http.StatusOK, code)
		equalJSON(t, map[string]interface{}{
			"risk-id": float64(i + 1),
		}, body)
		riskList[i].RiskID = i + 1
	}

	// Retrieve voyage info and check risk entries
	code, body, err := request(http.MethodGet, "/voyage", map[string]interface{}{
		"voyage-id": float64(1),
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	retrievedVoyage := voyage{}
	err = json.Unmarshal([]byte(body), &retrievedVoyage)
	assert.Equal(t, nil, err)
	sort.Slice(riskList, func(i, j int) bool {
		if !riskList[i].Time.IsZero() && !riskList[j].Time.IsZero() {
			return riskList[i].Time.After(riskList[j].Time)
		} else {
			return riskList[i].RiskID > riskList[j].RiskID
		}
	})
	assert.Equal(t, riskList, retrievedVoyage.RiskList)
}

func TestVoyageUpdateDoesntEraseTime(t *testing.T) {
	setupVoyageStorage()
	setupRiskStorage()

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":    float64(1),
		"vessel-id":    2,
		"start-hours":  101,
		"start-time":   "2022-01-03T15:13:12Z",
		"end-time":     "0001-01-01T00:00:00Z",
		"title":        "Breakdown Coomera",
		"risk-history": []risk{},
		"weather": map[string]interface{}{
			"seaway-tide": map[string]interface{}{
				"height-metres": 1.12,
				"time":          "2022-01-01T11:12:13Z",
			},
		},
	}, map[string]interface{}{
		"vessel-id":   2,
		"start-hours": 101,
		"start-time":  "2022-01-03T15:13:12Z",
		"title":       "Breakdown Coomera",
		"weather": map[string]interface{}{
			"seaway-tide": map[string]interface{}{
				"height-metres": 1.12,
				"time":          "2022-01-01T11:12:13Z",
			},
		},
	})

	testVoyageStoreAndRetrieve(t, map[string]interface{}{
		"voyage-id":    float64(1),
		"vessel-id":    2,
		"start-hours":  101,
		"start-time":   "2022-01-03T15:13:12Z",
		"end-time":     "0001-01-01T00:00:00Z",
		"title":        "Breakdown Coomera",
		"desc":         "Breakdown of 14' cruiser at coomera waters",
		"risk-history": []risk{},
		"weather": map[string]interface{}{
			"seaway-tide": map[string]interface{}{
				"height-metres": 1.12,
				"time":          "2022-01-01T11:12:13Z",
			},
		},
	}, map[string]interface{}{
		"voyage-id": float64(1),
		"desc":      "Breakdown of 14' cruiser at coomera waters",
	})
}

func TestVoyageList(t *testing.T) {
	setupVoyageStorage()
	setupRiskStorage()
	expectList := []map[string]interface{}{
		{
			"voyage-id":    float64(1),
			"vessel-id":    2,
			"start-hours":  101,
			"start-time":   "2022-01-03T15:13:12Z",
			"end-time":     "0001-01-01T00:00:00Z",
			"title":        "Breakdown Coomera",
			"risk-history": []risk{},
			"weather": map[string]interface{}{
				"seaway-tide": map[string]interface{}{
					"height-metres": 1.12,
					"time":          "2022-01-03T11:12:13Z",
				},
			},
		}, {
			"voyage-id":    float64(2),
			"vessel-id":    2,
			"start-hours":  101,
			"start-time":   "2022-01-04T05:03:02Z",
			"end-time":     "0001-01-01T00:00:00Z",
			"title":        "Breakdown Coomera Waters",
			"risk-history": []risk{},
			"weather": map[string]interface{}{
				"seaway-tide": map[string]interface{}{
					"height-metres": 2.12,
					"time":          "2022-01-04T11:12:13Z",
				},
			},
		},
	}

	testVoyageStoreAndRetrieve(t, expectList[0], map[string]interface{}{
		"vessel-id":   2,
		"start-hours": 101,
		"start-time":  "2022-01-03T15:13:12Z",
		"title":       "Breakdown Coomera",
		"weather": map[string]interface{}{
			"seaway-tide": map[string]interface{}{
				"height-metres": 1.12,
				"time":          "2022-01-03T11:12:13Z",
			},
		},
	})

	testVoyageStoreAndRetrieve(t, expectList[1], map[string]interface{}{
		"vessel-id":   2,
		"start-hours": 101,
		"start-time":  "2022-01-04T05:03:02Z",
		"title":       "Breakdown Coomera Waters",
		"weather": map[string]interface{}{
			"seaway-tide": map[string]interface{}{
				"height-metres": 2.12,
				"time":          "2022-01-04T11:12:13Z",
			},
		},
	})

	// List all voyages
	code, body, err := request(http.MethodGet, "/voyage/list", map[string]interface{}{
		"vessel-id": expectList[0]["vessel-id"],
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	equalVoyageMapList(t, expectList, body)
}
