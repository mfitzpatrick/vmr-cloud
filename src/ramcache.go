package main

import (
	"context"
	"sort"

	"github.com/pkg/errors"
)

func init() {
	ramcacheVoyageInit()
	ramcacheAssistInit()
	ramcacheRiskInit()
}

// VOYAGE RAMCACHE ITEMS
type voyageSerialiseResponse struct {
	val voyage
	err error
}

type voyageSerialiseCMD struct {
	set  bool
	val  voyage
	done chan<- voyageSerialiseResponse
}

var voyageCache map[int]voyage
var voyageCacheChan chan voyageSerialiseCMD

func ramcacheVoyageInit() {
	voyageCache = make(map[int]voyage)
	voyageCacheChan = make(chan voyageSerialiseCMD)
	newVoyageEntryID := func() int {
		// Return an integer which is 1 more than the maximum index currently in the map
		max := 0
		for num, _ := range voyageCache {
			if num > max {
				max = num
			}
		}
		return max + 1
	}
	set := func(cmd voyageSerialiseCMD) (voyage, error) {
		if cmd.val.VoyageID == 0 && cmd.val.VesselID == 0 {
			return voyage{}, INVALID_VESSEL_ID.Errorf("voyage storage map set")
		} else if cmd.val.VoyageID == 0 && cmd.val.VesselID != 0 {
			// A new record should be created. Get the ID
			cmd.val.VoyageID = newVoyageEntryID()
		}
		if cached, ok := voyageCache[cmd.val.VoyageID]; ok {
			if merged, err := mergeVoyageStructs(cached, cmd.val); err != nil {
				return voyage{}, STORAGE_FAIL.Errorf("voyage storage map merge: %v", err)
			} else {
				voyageCache[cmd.val.VoyageID] = merged
			}
		} else {
			voyageCache[cmd.val.VoyageID] = cmd.val
		}
		return voyage{VoyageID: cmd.val.VoyageID}, nil
	}
	get := func(cmd voyageSerialiseCMD) (voyage, error) {
		if cmd.val.VoyageID == 0 {
			return voyage{}, INVALID_VOYAGE_ID.Errorf("voyage storage map get")
		}
		if val, ok := voyageCache[cmd.val.VoyageID]; !ok {
			return voyage{}, VOYAGE_NOT_FOUND.Errorf("voyage storage map get")
		} else {
			return val, nil
		}
	}
	go func() {
		for {
			select {
			case cmd := <-voyageCacheChan:
				resp := voyageSerialiseResponse{}
				if cmd.set {
					resp.val, resp.err = set(cmd)
				} else {
					resp.val, resp.err = get(cmd)
				}
				cmd.done <- resp
			}
		}
	}()
}

func storeVoyage(ctx context.Context, v voyage) (int, error) {
	done := make(chan voyageSerialiseResponse)
	cmd := voyageSerialiseCMD{
		set:  true,
		val:  v,
		done: done,
	}
	select { //Send to map
	case voyageCacheChan <- cmd:
	case <-ctx.Done():
		return 0, errors.Wrapf(ctx.Err(), "store voyage")
	}
	select { //retrieve answer
	case resp := <-done:
		return resp.val.VoyageID, errors.Wrapf(resp.err, "store voyage")
	case <-ctx.Done():
		return 0, errors.Wrapf(ctx.Err(), "store voyage")
	}
	return 0, STORAGE_FAIL.Errorf("store voyage")
}

func retrieveVoyage(ctx context.Context, voyageID int) (voyage, error) {
	done := make(chan voyageSerialiseResponse)
	cmd := voyageSerialiseCMD{
		val: voyage{
			VoyageID: voyageID,
		},
		done: done,
	}
	select { //Send to map
	case voyageCacheChan <- cmd:
	case <-ctx.Done():
		return voyage{}, errors.Wrapf(ctx.Err(), "retrieve voyage")
	}
	select { //retrieve answer
	case resp := <-done:
		if resp.err != nil {
			return voyage{}, errors.Wrapf(resp.err, "retrieve voyage")
		} else if riskList, err := retrieveRiskForVoyage(ctx, resp.val.VoyageID); err != nil {
			return voyage{}, RETRIEVAL_FAIL.Errorf("retrieve voyage risk list")
		} else {
			resp.val.RiskList = riskList
		}
		return resp.val, nil
	case <-ctx.Done():
		return voyage{}, errors.Wrapf(ctx.Err(), "retrieve voyage")
	}
	return voyage{}, RETRIEVAL_FAIL.Errorf("retrieve voyage")
}

// RAMCACHE ASSIST ITEMS
type assistSerialiseResponse struct {
	val assist
	err error
}

type assistSerialiseCMD struct {
	set  bool
	val  assist
	done chan<- assistSerialiseResponse
}

var assistCache map[int]assist
var assistCacheChan chan assistSerialiseCMD

func ramcacheAssistInit() {
	assistCache = make(map[int]assist)
	assistCacheChan = make(chan assistSerialiseCMD)
	newAssistEntryID := func() int {
		// Return an integer which is 1 more than the maximum index currently in the map
		max := 0
		for num, _ := range assistCache {
			if num > max {
				max = num
			}
		}
		return max + 1
	}
	set := func(cmd assistSerialiseCMD) (assist, error) {
		if cmd.val.AssistID == 0 && cmd.val.VoyageID == 0 {
			return assist{}, INVALID_VOYAGE_ID.Errorf("assist storage map set")
		} else if cmd.val.AssistID == 0 && cmd.val.VoyageID != 0 {
			// A new record should be created. Get the ID
			cmd.val.AssistID = newAssistEntryID()
		}
		if cached, ok := assistCache[cmd.val.AssistID]; ok {
			if merged, err := mergeAssistStructs(cached, cmd.val); err != nil {
				return assist{}, STORAGE_FAIL.Errorf("assist storage map merge: %v", err)
			} else {
				assistCache[cmd.val.AssistID] = merged
			}
		} else {
			assistCache[cmd.val.AssistID] = cmd.val
		}
		return assist{AssistID: cmd.val.AssistID}, nil
	}
	get := func(cmd assistSerialiseCMD) (assist, error) {
		if cmd.val.AssistID == 0 {
			return assist{}, INVALID_ASSIST_ID.Errorf("assist storage map get")
		}
		if val, ok := assistCache[cmd.val.AssistID]; !ok {
			return assist{}, ASSIST_NOT_FOUND.Errorf("assist storage map get")
		} else {
			return val, nil
		}
	}
	go func() {
		for {
			select {
			case cmd := <-assistCacheChan:
				resp := assistSerialiseResponse{}
				if cmd.set {
					resp.val, resp.err = set(cmd)
				} else {
					resp.val, resp.err = get(cmd)
				}
				cmd.done <- resp
			}
		}
	}()
}

func storeAssist(ctx context.Context, v assist) (int, error) {
	done := make(chan assistSerialiseResponse)
	cmd := assistSerialiseCMD{
		set:  true,
		val:  v,
		done: done,
	}
	select { //Send to map
	case assistCacheChan <- cmd:
	case <-ctx.Done():
		return 0, errors.Wrapf(ctx.Err(), "store assist")
	}
	select { //retrieve answer
	case resp := <-done:
		return resp.val.AssistID, errors.Wrapf(resp.err, "store assist")
	case <-ctx.Done():
		return 0, errors.Wrapf(ctx.Err(), "store assist")
	}
	return 0, STORAGE_FAIL.Errorf("store assist")
}

func retrieveAssist(ctx context.Context, assistID int) (assist, error) {
	done := make(chan assistSerialiseResponse)
	cmd := assistSerialiseCMD{
		val: assist{
			AssistID: assistID,
		},
		done: done,
	}
	select { //Send to map
	case assistCacheChan <- cmd:
	case <-ctx.Done():
		return assist{}, errors.Wrapf(ctx.Err(), "retrieve assist")
	}
	select { //retrieve answer
	case resp := <-done:
		return resp.val, errors.Wrapf(resp.err, "retrieve assist")
	case <-ctx.Done():
		return assist{}, errors.Wrapf(ctx.Err(), "retrieve assist")
	}
	return assist{}, RETRIEVAL_FAIL.Errorf("retrieve assist")
}

// RISK RAMCACHE ITEMS
type riskSerialiseResponse struct {
	val []risk
	err error
}

type riskSerialiseCMD struct {
	set    bool
	getAll bool
	val    risk
	done   chan<- riskSerialiseResponse
}

var riskCache map[int]risk
var riskCacheChan chan riskSerialiseCMD

func ramcacheRiskInit() {
	riskCache = make(map[int]risk)
	riskCacheChan = make(chan riskSerialiseCMD)
	newRiskEntryID := func() int {
		// Return an integer which is 1 more than the maximum index currently in the map
		max := 0
		for num, _ := range riskCache {
			if num > max {
				max = num
			}
		}
		return max + 1
	}
	findAllForVoyageID := func(voyageID int) []risk {
		riskOut := []risk{}
		for _, v := range riskCache {
			if v.VoyageID == voyageID {
				riskOut = append(riskOut, v)
			}
		}
		// Return the slice sorted by time descending
		sort.Slice(riskOut, func(i, j int) bool {
			if !riskOut[i].Time.IsZero() && !riskOut[j].Time.IsZero() {
				return riskOut[i].Time.After(riskOut[j].Time)
			} else {
				return riskOut[i].RiskID > riskOut[j].RiskID
			}
		})
		return riskOut
	}
	set := func(cmd riskSerialiseCMD) ([]risk, error) {
		if cmd.val.RiskID != 0 {
			return []risk{}, IMMUTABLE_RISK.Errorf("risk storage map set")
		} else if cmd.val.VoyageID == 0 {
			return []risk{}, INVALID_VOYAGE_ID.Errorf("risk storage map set")
		} else {
			// A new record should be created. Get the ID
			cmd.val.RiskID = newRiskEntryID()
		}
		riskCache[cmd.val.RiskID] = cmd.val
		return []risk{{RiskID: cmd.val.RiskID}}, nil
	}
	getAll := func(cmd riskSerialiseCMD) ([]risk, error) {
		if cmd.val.VoyageID == 0 {
			return []risk{}, INVALID_VOYAGE_ID.Errorf("risk storage map getAll voyage must be nonzero")
		} else if cmd.val.RiskID != 0 {
			return []risk{}, INVALID_RISK_ID.Errorf("risk storage map getAll risk cannot be set")
		}
		// The request is to return every item tagged with the voyage ID.
		return findAllForVoyageID(cmd.val.VoyageID), nil
	}
	get := func(cmd riskSerialiseCMD) ([]risk, error) {
		if cmd.val.VoyageID != 0 {
			return []risk{}, INVALID_VOYAGE_ID.Errorf("risk storage map get voyage cannot be set")
		} else if cmd.val.RiskID == 0 {
			return []risk{}, INVALID_RISK_ID.Errorf("risk storage map get risk must be nonzero")
		}
		// The request is to return a single item with a given risk ID.
		if val, ok := riskCache[cmd.val.RiskID]; !ok {
			return []risk{}, RISK_NOT_FOUND.Errorf("risk storage map get no risk for ID %d",
				cmd.val.RiskID)
		} else {
			return []risk{val}, nil
		}
	}
	go func() {
		for {
			select {
			case cmd := <-riskCacheChan:
				resp := riskSerialiseResponse{}
				if cmd.set {
					resp.val, resp.err = set(cmd)
				} else if cmd.getAll {
					resp.val, resp.err = getAll(cmd)
				} else {
					resp.val, resp.err = get(cmd)
				}
				cmd.done <- resp
			}
		}
	}()
}

func storeRisk(ctx context.Context, v risk) (int, error) {
	done := make(chan riskSerialiseResponse)
	cmd := riskSerialiseCMD{
		set:  true,
		val:  v,
		done: done,
	}
	select { //Send to map
	case riskCacheChan <- cmd:
	case <-ctx.Done():
		return 0, errors.Wrapf(ctx.Err(), "store risk")
	}
	select { //retrieve answer
	case resp := <-done:
		if resp.err != nil {
			return 0, errors.Wrapf(resp.err, "store risk")
		} else if len(resp.val) != 1 {
			return 0, STORAGE_FAIL.Errorf("store risk returned too many entries: %d", len(resp.val))
		}
		return resp.val[0].RiskID, nil
	case <-ctx.Done():
		return 0, errors.Wrapf(ctx.Err(), "store risk")
	}
	return 0, STORAGE_FAIL.Errorf("store risk")
}

func retrieveRisk(ctx context.Context, riskID int) (risk, error) {
	done := make(chan riskSerialiseResponse)
	cmd := riskSerialiseCMD{
		val: risk{
			RiskID: riskID,
		},
		done: done,
	}
	select { //Send to map
	case riskCacheChan <- cmd:
	case <-ctx.Done():
		return risk{}, errors.Wrapf(ctx.Err(), "retrieve risk")
	}
	select { //retrieve answer
	case resp := <-done:
		if resp.err != nil {
			return risk{}, errors.Wrapf(resp.err, "retrieve risk")
		} else if len(resp.val) != 1 {
			return risk{}, RETRIEVAL_FAIL.Errorf("retrieve risk too many entries: %d", len(resp.val))
		}
		return resp.val[0], nil
	case <-ctx.Done():
		return risk{}, errors.Wrapf(ctx.Err(), "retrieve risk")
	}
	return risk{}, RETRIEVAL_FAIL.Errorf("retrieve risk")
}

func retrieveRiskForVoyage(ctx context.Context, voyageID int) ([]risk, error) {
	done := make(chan riskSerialiseResponse)
	cmd := riskSerialiseCMD{
		getAll: true,
		val: risk{
			VoyageID: voyageID,
		},
		done: done,
	}
	select { //Send to map
	case riskCacheChan <- cmd:
	case <-ctx.Done():
		return []risk{}, errors.Wrapf(ctx.Err(), "retrieve risk")
	}
	select { //retrieve answer
	case resp := <-done:
		return resp.val, errors.Wrapf(resp.err, "retrieve risk")
	case <-ctx.Done():
		return []risk{}, errors.Wrapf(ctx.Err(), "retrieve risk")
	}
	return []risk{}, RETRIEVAL_FAIL.Errorf("retrieve risk")
}
