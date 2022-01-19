package main

import (
	"context"
	"encoding/json"
	"time"
)

type vessel struct {
	VesselID int `json:"vessel-id"`
}

type crew struct {
	CrewID int    `json:"crew-id"`
	Name   string `json:"name"`
	Rank   string `json:"rank"`
}

type wind struct {
	Speed int `json:"speed-knots"`
	Dir   int `json:"direction-degrees"`
}

type tide struct {
	Height float64   `json:"height-metres"`
	Time   time.Time `json:"time"`
}

type weather struct {
	Wind wind `json:"wind"`
	Tide tide `json:"seaway-tide"`
}

type voyage struct {
	vessel

	VoyageID int     `json:"voyage-id"`
	RiskList []risk  `json:"risk-history"`
	Weather  weather `json:"weather"`

	StartTime        time.Time `json"start-time"`
	StartEngineHours int       `json"start-hours"`
	EndTime          time.Time `json"end-time"`
	EndEngineHours   int       `json"end-hours"`

	Skipper  crew   `json:"skipper"`
	CrewList []crew `json:"crew"`

	Title string `json:"title"`
	Desc  string `json:"description"`
}

func postVoyage(ctx context.Context, body string) ([]byte, error) {
	type request struct {
		voyage
	}
	var req request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return []byte{}, JSON_UNMARSHAL.Errorf("post voyage: %v with input %s", err, body)
	}
	if req.VesselID == 0 && req.VoyageID == 0 {
		return []byte{}, INVALID_VOYAGE_ID.Errorf("post voyage - set vessel ID to create new voyage entry")
	}
	if voyageID, err := storeVoyage(ctx, req.voyage); err != nil {
		return []byte{}, STORAGE_FAIL.Errorf("post voyage: %v", err)
	} else if b, err := json.Marshal(struct {
		VoyageID int `json:"voyage-id"`
	}{
		VoyageID: voyageID,
	}); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("post voyage: %v", err)
	} else {
		return b, nil
	}
}

func getVoyage(ctx context.Context, body string) ([]byte, error) {
	type request struct {
		VoyageID int `json:"voyage-id"`
	}
	var req request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return []byte{}, JSON_UNMARSHAL.Errorf("get voyage: %v with input %s", err, body)
	}
	if req.VoyageID == 0 {
		return []byte{}, INVALID_VOYAGE_ID.Errorf("get voyage")
	}
	if foundItem, err := retrieveVoyage(ctx, req.VoyageID); err != nil {
		return []byte{}, RETRIEVAL_FAIL.Errorf("get voyage: %v", err)
	} else if b, err := json.Marshal(foundItem); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("get voyage: %v", err)
	} else {
		return b, nil
	}
}
