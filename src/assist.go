package main

import (
	"context"
	"encoding/json"
)

type clientVessel struct {
	Name   string  `json:"name"`
	Rego   string  `json:"registration"`
	Type   string  `json:"type"`
	Colour string  `json:"colour"`
	Length float64 `json:"length-metres"`
	POB    int     `json:"pob"`
}

type client struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	MemberNo int    `json:"member-number"`

	Vessel clientVessel `json:"vessel"`
}

type coordinate struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type location struct {
	Name   string     `json:"name"`
	GPS    coordinate `json:"gps"`
	Depth  int        `json:"depth-metres"`
	Status string     `json:"status"`
}

type assist struct {
	VoyageID int `json:"voyage-id"` // For linking an assist to a voyage

	AssistID int `json:"assist-id"`

	Client client   `json:"client"`
	Pickup location `json:"pickup-location"`

	Problem string `json:"problem"`
	Action  string `json:"action"`

	Dest location `json:"destination"`
}

func postAssist(ctx context.Context, body string) ([]byte, error) {
	type request struct {
		assist
	}
	var req request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return []byte{}, JSON_UNMARSHAL.Errorf("post assist: %v with input %s", err, body)
	}
	if req.AssistID == 0 && req.VoyageID == 0 {
		return []byte{}, INVALID_VOYAGE_ID.Errorf("post assist - set voyage ID to create a new assist entry")
	}
	if assistID, err := storeAssist(ctx, req.assist); err != nil {
		return []byte{}, STORAGE_FAIL.Errorf("post assist: %v", err)
	} else if b, err := json.Marshal(struct {
		AssistID int `json:"assist-id"`
	}{
		AssistID: assistID,
	}); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("post assist: %v", err)
	} else {
		return b, nil
	}
}

func getAssist(ctx context.Context, body string) ([]byte, error) {
	type request struct {
		AssistID int `json:"assist-id"`
	}
	var req request
	if err := json.Unmarshal([]byte(body), &req); err != nil {
		return []byte{}, JSON_UNMARSHAL.Errorf("get assist: %v with input %s", err, body)
	}
	if req.AssistID == 0 {
		return []byte{}, INVALID_ASSIST_ID.Errorf("get assist")
	}
	if foundItem, err := retrieveAssist(ctx, req.AssistID); err != nil {
		return []byte{}, RETRIEVAL_FAIL.Errorf("get assist: %v", err)
	} else if b, err := json.Marshal(foundItem); err != nil {
		return []byte{}, JSON_MARSHAL.Errorf("get assist: %v", err)
	} else {
		return b, nil
	}
}
