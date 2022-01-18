package main

import (
	"context"
	"encoding/json"
)

type risk struct {
	RiskID int `json:"risk-id"`
	Type   int `json:"type"`
	Mgmt   int `json:"management"`
	Crew   int `json:"crew"`
	Equip  int `json:"equipment"`
	Env    int `json:"environment"`
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
		return []byte{}, INVALID_RISK_ID.Errorf("new risk")
	}
	if b, err := json.Marshal(struct {
		RiskID int `json:"risk-id"`
	}{
		RiskID: 1,
	}); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("new risk: %v", err)
	} else {
		return b, nil
	}
}
