package main

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func setupVoyageStorage() {
	voyageCache = make(map[int]voyage)
}

func TestStoreVoyage(t *testing.T) {
	setupVoyageStorage()

	testVoyage := voyage{
		VoyageID:         1,
		StartEngineHours: 101,
	}
	voyageID, err := storeVoyage(context.Background(), testVoyage)
	assert.Equal(t, nil, err)
	assert.Equal(t, testVoyage.VoyageID, voyageID)

	mapEntry, ok := voyageCache[1]
	assert.Equal(t, true, ok)
	assert.Equal(t, testVoyage, mapEntry)
}

func TestRetrieveVoyage(t *testing.T) {
	setupVoyageStorage()
	voyageCache[1] = voyage{
		VoyageID:         1,
		StartEngineHours: 101,
	}

	vEntry, err := retrieveVoyage(context.Background(), 1)
	assert.Equal(t, nil, err)

	mapEntry, ok := voyageCache[1]
	assert.Equal(t, true, ok)
	mapEntry.RiskList = []risk{}
	assert.Equal(t, mapEntry, vEntry)
}

func TestRetrieveVoyageList(t *testing.T) {
	setupVoyageStorage()
	getTime := func(t *testing.T, str string) time.Time {
		tm, err := time.Parse(time.RFC3339, str)
		assert.Equal(t, nil, err)
		return tm
	}
	voyageList := []voyage{
		{
			VoyageID: 1,
			vessel: vessel{
				VesselID: 1,
			},
			StartEngineHours: 101,
			StartTime:        getTime(t, "2022-01-01T01:02:03Z"),
			RiskList:         []risk{},
		},
		{
			VoyageID: 2,
			vessel: vessel{
				VesselID: 1,
			},
			StartEngineHours: 103,
			StartTime:        getTime(t, "2022-01-01T03:02:03Z"),
			RiskList:         []risk{},
		},
	}
	for _, v := range voyageList {
		voyageCache[v.VoyageID] = v
	}

	vEntryList, err := retrieveVoyageList(context.Background(), voyageList[0].VesselID)
	assert.Equal(t, nil, err)
	assert.Equal(t, voyageList, vEntryList)
}

func TestStoreVoyageMultipleEntries(t *testing.T) {
	setupVoyageStorage()

	testVoyage := voyage{
		vessel: vessel{
			VesselID: 1,
		},
		StartEngineHours: 101,
	}
	voyageID, err := storeVoyage(context.Background(), testVoyage)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, voyageID)

	mapEntry, ok := voyageCache[1]
	assert.Equal(t, true, ok)
	testVoyage.VoyageID = 1
	assert.Equal(t, testVoyage, mapEntry)

	testVoyage = voyage{
		vessel: vessel{
			VesselID: 1,
		},
		StartEngineHours: 102,
	}
	voyageID, err = storeVoyage(context.Background(), testVoyage)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, voyageID)

	mapEntry, ok = voyageCache[2]
	assert.Equal(t, true, ok)
	testVoyage.VoyageID = 2
	assert.Equal(t, testVoyage, mapEntry)
}

func TestRetrieveVoyageMultipleEntries(t *testing.T) {
	setupVoyageStorage()
	voyageCache[1] = voyage{
		VoyageID:         1,
		StartEngineHours: 101,
	}
	voyageCache[2] = voyage{
		VoyageID:         2,
		StartEngineHours: 102,
	}
	voyageCache[8] = voyage{
		VoyageID:         8,
		StartEngineHours: 108,
		Skipper:          crew{Name: "Bob"},
	}

	vEntry, err := retrieveVoyage(context.Background(), 1)
	assert.Equal(t, nil, err)

	mapEntry, ok := voyageCache[1]
	assert.Equal(t, true, ok)
	mapEntry.RiskList = []risk{}
	assert.Equal(t, mapEntry, vEntry)

	vEntry, err = retrieveVoyage(context.Background(), 2)
	assert.Equal(t, nil, err)

	mapEntry, ok = voyageCache[2]
	assert.Equal(t, true, ok)
	mapEntry.RiskList = []risk{}
	assert.Equal(t, mapEntry, vEntry)

	vEntry, err = retrieveVoyage(context.Background(), 8)
	assert.Equal(t, nil, err)

	mapEntry, ok = voyageCache[8]
	assert.Equal(t, true, ok)
	mapEntry.RiskList = []risk{}
	assert.Equal(t, mapEntry, vEntry)
}

func TestStoreVoyageUpdateOK(t *testing.T) {
	setupVoyageStorage()
	getTime := func(t *testing.T, str string) time.Time {
		tm, err := time.Parse(time.RFC3339, str)
		assert.Equal(t, nil, err)
		return tm
	}

	testVoyage := voyage{
		vessel: vessel{
			VesselID: 1,
		},
		StartEngineHours: 101,
		StartTime:        getTime(t, "2022-01-01T01:11:12Z"),
	}
	voyageID, err := storeVoyage(context.Background(), testVoyage)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, voyageID)

	mapEntry, ok := voyageCache[1]
	assert.Equal(t, true, ok)
	testVoyage.VoyageID = 1
	assert.Equal(t, testVoyage, mapEntry)

	// Now update the existing entry and check that the structs are properly merged
	testVoyage = voyage{
		VoyageID: 1,
		Title:    "some title",
		Desc:     "some description",
	}
	voyageID, err = storeVoyage(context.Background(), testVoyage)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, voyageID)

	mapEntry, ok = voyageCache[1]
	assert.Equal(t, true, ok)
	testVoyage.VoyageID = 2
	assert.Equal(t, voyage{
		VoyageID: 1,
		vessel: vessel{
			VesselID: 1,
		},
		StartEngineHours: 101,
		StartTime:        getTime(t, "2022-01-01T01:11:12Z"),
		Title:            "some title",
		Desc:             "some description",
	}, mapEntry)
}

func setupAssistStorage() {
	assistCache = make(map[int]assist)
}

func TestStoreAssist(t *testing.T) {
	setupAssistStorage()

	testAssist := assist{
		AssistID: 1,
		Problem:  "Breakdown",
	}
	assistID, err := storeAssist(context.Background(), testAssist)
	assert.Equal(t, nil, err)
	assert.Equal(t, testAssist.AssistID, assistID)

	mapEntry, ok := assistCache[1]
	assert.Equal(t, true, ok)
	assert.Equal(t, testAssist, mapEntry)
}

func TestRetrieveAssist(t *testing.T) {
	setupAssistStorage()
	assistCache[1] = assist{
		AssistID: 1,
		Problem:  "Breakdown",
	}

	vEntry, err := retrieveAssist(context.Background(), 1)
	assert.Equal(t, nil, err)

	mapEntry, ok := assistCache[1]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)
}

func TestStoreAssistMultipleEntries(t *testing.T) {
	setupAssistStorage()

	testAssist := assist{
		VoyageID: 1,
		Problem:  "Breakdown",
	}
	assistID, err := storeAssist(context.Background(), testAssist)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, assistID)

	mapEntry, ok := assistCache[1]
	assert.Equal(t, true, ok)
	testAssist.AssistID = 1
	assert.Equal(t, testAssist, mapEntry)

	testAssist = assist{
		VoyageID: 1,
		Problem:  "Breakdown",
	}
	assistID, err = storeAssist(context.Background(), testAssist)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, assistID)

	mapEntry, ok = assistCache[2]
	assert.Equal(t, true, ok)
	testAssist.AssistID = 2
	assert.Equal(t, testAssist, mapEntry)
}

func TestRetrieveAssistMultipleEntries(t *testing.T) {
	setupAssistStorage()
	assistCache[1] = assist{
		AssistID: 1,
		Problem:  "Breakdown",
	}
	assistCache[2] = assist{
		AssistID: 2,
		Problem:  "Jump Start",
	}
	assistCache[8] = assist{
		AssistID: 8,
		Problem:  "Search",
		Client:   client{Name: "Bob"},
	}

	vEntry, err := retrieveAssist(context.Background(), 1)
	assert.Equal(t, nil, err)

	mapEntry, ok := assistCache[1]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)

	vEntry, err = retrieveAssist(context.Background(), 2)
	assert.Equal(t, nil, err)

	mapEntry, ok = assistCache[2]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)

	vEntry, err = retrieveAssist(context.Background(), 8)
	assert.Equal(t, nil, err)

	mapEntry, ok = assistCache[8]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)
}

func TestStoreAssistUpdateOK(t *testing.T) {
	setupAssistStorage()

	testAssist := assist{
		VoyageID: 1,
		Problem:  "Search",
	}
	assistID, err := storeAssist(context.Background(), testAssist)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, assistID)

	mapEntry, ok := assistCache[1]
	assert.Equal(t, true, ok)
	testAssist.AssistID = 1
	assert.Equal(t, testAssist, mapEntry)

	// Now update the existing entry and check that the structs are properly merged
	testAssist = assist{
		AssistID: 1,
		Action:   "some action",
	}
	assistID, err = storeAssist(context.Background(), testAssist)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, assistID)

	mapEntry, ok = assistCache[1]
	assert.Equal(t, true, ok)
	testAssist.AssistID = 2
	assert.Equal(t, assist{
		AssistID: 1,
		VoyageID: 1,
		Problem:  "Search",
		Action:   "some action",
	}, mapEntry)
}

func setupRiskStorage() {
	riskCache = make(map[int]risk)
}

func TestStoreRisk(t *testing.T) {
	setupRiskStorage()

	testRisk := risk{
		VoyageID: 1,
	}
	riskID, err := storeRisk(context.Background(), testRisk)
	assert.Equal(t, nil, err)
	testRisk.RiskID = 1
	assert.Equal(t, testRisk.RiskID, riskID)

	mapEntry, ok := riskCache[1]
	assert.Equal(t, true, ok)
	assert.Equal(t, testRisk, mapEntry)
}

func TestUpdateRiskFails(t *testing.T) {
	setupRiskStorage()

	_, err := storeRisk(context.Background(), risk{RiskID: 1})
	var herror httpError
	ok := errors.As(err, &herror)
	assert.Equal(t, true, ok)
	assert.Equal(t, IMMUTABLE_RISK.ID, herror.ID)
}

func TestRetrieveRisk(t *testing.T) {
	setupRiskStorage()
	riskCache[1] = risk{
		RiskID:   1,
		VoyageID: 1,
	}

	vEntry, err := retrieveRisk(context.Background(), 1)
	assert.Equal(t, nil, err)

	mapEntry, ok := riskCache[1]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)
}

func TestStoreRiskMultipleEntries(t *testing.T) {
	setupRiskStorage()

	testRisk := risk{
		VoyageID: 1,
	}
	riskID, err := storeRisk(context.Background(), testRisk)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, riskID)

	mapEntry, ok := riskCache[1]
	assert.Equal(t, true, ok)
	testRisk.RiskID = 1
	assert.Equal(t, testRisk, mapEntry)

	testRisk = risk{
		VoyageID: 1,
	}
	riskID, err = storeRisk(context.Background(), testRisk)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, riskID)

	mapEntry, ok = riskCache[2]
	assert.Equal(t, true, ok)
	testRisk.RiskID = 2
	assert.Equal(t, testRisk, mapEntry)
}

func TestRetrieveRiskMultipleEntries(t *testing.T) {
	setupRiskStorage()
	riskCache[1] = risk{
		RiskID:   1,
		VoyageID: 1,
	}
	riskCache[2] = risk{
		RiskID:   2,
		VoyageID: 1,
	}
	riskCache[8] = risk{
		RiskID:   8,
		VoyageID: 1,
	}

	vEntry, err := retrieveRisk(context.Background(), 1)
	assert.Equal(t, nil, err)

	mapEntry, ok := riskCache[1]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)

	vEntry, err = retrieveRisk(context.Background(), 2)
	assert.Equal(t, nil, err)

	mapEntry, ok = riskCache[2]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)

	vEntry, err = retrieveRisk(context.Background(), 8)
	assert.Equal(t, nil, err)

	mapEntry, ok = riskCache[8]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)
}

func TestRetrieveRiskForVoyage(t *testing.T) {
	setupRiskStorage()
	riskList := []risk{{
		RiskID:   1,
		VoyageID: 1,
		Mgmt:     1,
	}, {
		RiskID:   2,
		VoyageID: 1,
		Mgmt:     3,
	}, {
		RiskID:   3,
		VoyageID: 1,
		Mgmt:     3,
	}}
	for i, v := range riskList {
		riskCache[i+1] = v
	}

	entryList, err := retrieveRiskForVoyage(context.Background(), 1)
	assert.Equal(t, nil, err)
	sort.Slice(riskList, func(i, j int) bool {
		if !riskList[i].Time.IsZero() && !riskList[j].Time.IsZero() {
			return riskList[i].Time.After(riskList[j].Time)
		} else {
			return riskList[i].RiskID > riskList[j].RiskID
		}
	})
	assert.Equal(t, riskList, entryList)

	entryList, err = retrieveRiskForVoyage(context.Background(), 2)
	assert.Equal(t, nil, err)
	assert.Equal(t, []risk{}, entryList)
}
