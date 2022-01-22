package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCrew(t *testing.T) {
	crewInterfaceMap := []map[string]interface{}{}
	for _, crewEntry := range hardcodedCrewList {
		crewInterfaceMap = append(crewInterfaceMap, toInterfaceMap(crewEntry))
	}
	code, body, err := request(http.MethodGet, "/crew/list", map[string]interface{}{})
	assert.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, code)
	equalCrewMapList(t, crewInterfaceMap, body)
}
