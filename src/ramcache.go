package main

import (
	"context"
	"github.com/pkg/errors"
)

func init() {
	ramcacheVoyageInit()
	ramcacheAssistInit()
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
		return resp.val, errors.Wrapf(resp.err, "retrieve voyage")
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
