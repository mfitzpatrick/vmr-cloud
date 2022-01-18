package main

import (
	"context"
	"encoding/json"
)

func newVoyage(ctx context.Context, body string) ([]byte, error) {
	type request struct {
		VesselID int `json:"vessel-id"`
	}
	var req request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return []byte{}, JSON_UNMARSHAL.Errorf("new voyage: %v with input %s", err, body)
	}
	if req.VesselID == 0 {
		return []byte{}, INVALID_VESSEL_ID.Errorf("new voyage")
	}
	if b, err := json.Marshal(struct {
		VoyageID int `json:"voyage-id"`
	}{
		VoyageID: 1,
	}); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("new voyage: %v", err)
	} else {
		return b, nil
	}
}

func getVoyage(ctx context.Context, body string) ([]byte, error) {
	return []byte{}, nil
}
