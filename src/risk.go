package main

import (
	"context"
	"encoding/json"
	"time"
)

type risk struct {
	RiskID   int       `json:"risk-id"`
	VoyageID int       `json:"voyage-id"`
	Time     time.Time `json:"time"`
	Type     int       `json:"type"`
	Mgmt     int       `json:"management"`
	Crew     int       `json:"crew"`
	Equip    int       `json:"equipment"`
	Env      int       `json:"environment"`
}

func newRisk(ctx context.Context, body string) ([]byte, error) {
	type request struct {
		risk
	}
	var req request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return []byte{}, JSON_UNMARSHAL.Errorf("new risk: %v with input %s", err, body)
	}
	if req.RiskID != 0 {
		return []byte{}, IMMUTABLE_RISK.Errorf("existing risk entries cannot be changed")
	}
	if req.VoyageID == 0 {
		return []byte{}, INVALID_VOYAGE_ID.Errorf("new risk")
	}
	if req.Time.IsZero() {
		return []byte{}, RISK_TIME_REQUIRED.Errorf("new risk time cannot be 0")
	}
	if riskID, err := storeRisk(ctx, req.risk); err != nil {
		return []byte{}, STORAGE_FAIL.Errorf("new risk")
	} else if b, err := json.Marshal(struct {
		RiskID int `json:"risk-id"`
	}{
		RiskID: riskID,
	}); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("new risk: %v", err)
	} else {
		return b, nil
	}
}
