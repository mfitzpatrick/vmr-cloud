package main

import (
	"context"
	"encoding/json"
	"reflect"
)

type storer interface {
	storeVoyage(context.Context, voyage) error
	retrieveVoyage(context.Context, int) (voyage, error)
}

// Recursive function to merge mapFrom into mapInto. When this function returns, the mapInto object will
// have all its existing fields, as well as the nonzero fields from mapFrom.
// NB: if a mapFrom field is a zero-value, this code will not erase the entry in mapInto.
func mergeMaps(mapInto, mapFrom map[string]interface{}) {
	for k, _ := range mapFrom {
		if valFrom, ok := mapFrom[k]; ok && valFrom != nil && !reflect.ValueOf(valFrom).IsZero() {
			if nestedMapFrom, ok := valFrom.(map[string]interface{}); ok {
				if nestedMapInto, ok := mapInto[k].(map[string]interface{}); ok {
					mergeMaps(nestedMapInto, nestedMapFrom)
				}
			} else {
				mapInto[k] = valFrom
			}
		}
	}
}

// Convert the struct object to a JSON object, then re-parse it as a generic map to interface.
// Merge the interface maps.
// Convert the interface maps to JSON, then back to the struct objects and return.
func mergeVoyageStructs(vInto, vFrom voyage) (voyage, error) {
	vOut := voyage{}
	mapInto := make(map[string]interface{})
	mapFrom := make(map[string]interface{})
	if bInto, err := json.Marshal(vInto); err != nil {
		return voyage{}, JSON_MARSHAL.Errorf("merge structs into")
	} else if bFrom, err := json.Marshal(vFrom); err != nil {
		return voyage{}, JSON_MARSHAL.Errorf("merge structs from")
	} else if err := json.Unmarshal(bInto, &mapInto); err != nil {
		return voyage{}, JSON_MARSHAL.Errorf("merge structs into unmarshal")
	} else if err := json.Unmarshal(bFrom, &mapFrom); err != nil {
		return voyage{}, JSON_MARSHAL.Errorf("merge structs from unmarshal")
	}

	mergeMaps(mapInto, mapFrom) //NB: recursive (bounded by size of map)

	if bInto, err := json.Marshal(mapInto); err != nil {
		return voyage{}, JSON_MARSHAL.Errorf("merge structs out marshal")
	} else if err := json.Unmarshal(bInto, &vOut); err != nil {
		return voyage{}, JSON_MARSHAL.Errorf("merge structs out unmarshal")
	}

	return vOut, nil
}

// Convert the struct object to a JSON object, then re-parse it as a generic map to interface.
// Merge the interface maps.
// Convert the interface maps to JSON, then back to the struct objects and return.
func mergeAssistStructs(vInto, vFrom assist) (assist, error) {
	vOut := assist{}
	mapInto := make(map[string]interface{})
	mapFrom := make(map[string]interface{})
	if bInto, err := json.Marshal(vInto); err != nil {
		return assist{}, JSON_MARSHAL.Errorf("merge structs into")
	} else if bFrom, err := json.Marshal(vFrom); err != nil {
		return assist{}, JSON_MARSHAL.Errorf("merge structs from")
	} else if err := json.Unmarshal(bInto, &mapInto); err != nil {
		return assist{}, JSON_MARSHAL.Errorf("merge structs into unmarshal")
	} else if err := json.Unmarshal(bFrom, &mapFrom); err != nil {
		return assist{}, JSON_MARSHAL.Errorf("merge structs from unmarshal")
	}

	mergeMaps(mapInto, mapFrom) //NB: recursive (bounded by size of map)

	if bInto, err := json.Marshal(mapInto); err != nil {
		return assist{}, JSON_MARSHAL.Errorf("merge structs out marshal")
	} else if err := json.Unmarshal(bInto, &vOut); err != nil {
		return assist{}, JSON_MARSHAL.Errorf("merge structs out unmarshal")
	}

	return vOut, nil
}
