package main

import (
	"context"
	"testing"

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
	assert.Equal(t, mapEntry, vEntry)
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
	assert.Equal(t, mapEntry, vEntry)

	vEntry, err = retrieveVoyage(context.Background(), 2)
	assert.Equal(t, nil, err)

	mapEntry, ok = voyageCache[2]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)

	vEntry, err = retrieveVoyage(context.Background(), 8)
	assert.Equal(t, nil, err)

	mapEntry, ok = voyageCache[8]
	assert.Equal(t, true, ok)
	assert.Equal(t, mapEntry, vEntry)
}

func TestStoreVoyageUpdateOK(t *testing.T) {
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
