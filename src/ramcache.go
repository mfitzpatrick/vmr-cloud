package main

import (
	"context"
	"github.com/pkg/errors"
)

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

func init() {
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
