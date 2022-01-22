package main

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

var hardcodedCrewList []crew = []crew{
	{
		Name:      "Meat Loaf",
		Rank:      "Inshore Skipper",
		IsSkipper: true,
	}, {
		Name: "Mavis Beacon",
		Rank: "New Recruit",
	}, {
		Name: "Anne Frank",
		Rank: "Recruit",
	}, {
		Name: "Mick Dundee",
		Rank: "Crew",
	}, {
		Name:      "Fonzie",
		Rank:      "Coxswain",
		IsSkipper: true,
	}, {
		Name:      "Batman",
		Rank:      "Offshore Skipper",
		IsSkipper: true,
	},
}

func init() {
	// Add a logical crew ID to each entry in the hardcoded crew list. Numbers start from 1, not 0
	for i, _ := range hardcodedCrewList {
		hardcodedCrewList[i].CrewID = i + 1
	}
}

func listCrew(ctx context.Context, body string, query url.Values) ([]byte, error) {
	var listLimit int = 50
	if val := query.Get("limit"); val != "" {
		if lim, err := strconv.Atoi(val); err != nil {
			return []byte{}, INVALID_QSTRING.Errorf("invalid limit format (int): %v", val)
		} else {
			listLimit = lim
		}
	}
	if listLimit > len(hardcodedCrewList) {
		listLimit = len(hardcodedCrewList)
	}
	crewList := hardcodedCrewList[:listLimit]
	if b, err := json.Marshal(crewList); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("list crew: %v", err)
	} else {
		return b, nil
	}
}
