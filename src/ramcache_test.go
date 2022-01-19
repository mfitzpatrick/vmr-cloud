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
