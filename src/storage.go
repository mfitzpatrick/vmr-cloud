package main

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"time"
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

// Recursive function to convert a structure into a map of string to interface. It uses the JSON struct
// tags to provide the map keys. Nested structs cause the function to recurse.
// If a nested struct is of type time.Time then this function will preserve the time object rather than
// creating a string representation (as the JSON module does).
func toInterfaceMap(theStruct interface{}) map[string]interface{} {
	theMap := map[string]interface{}{}
	structValue := reflect.ValueOf(theStruct)
	structType := reflect.TypeOf(theStruct)
	if structValue.Kind() != reflect.Struct {
		return map[string]interface{}{}
	}
	for _, field := range reflect.VisibleFields(structType) {
		fieldIDX := field.Index
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			jsonName := strings.SplitN(jsonTag, ",", 2)[0]
			if field.Type == reflect.TypeOf(time.Time{}) {
				theMap[jsonName] = structValue.FieldByIndex(fieldIDX).Interface()
			} else if structValue.FieldByIndex(fieldIDX).Kind() == reflect.Struct {
				theMap[jsonName] = toInterfaceMap(structValue.FieldByIndex(fieldIDX).Interface())
			} else {
				theMap[jsonName] = structValue.FieldByIndex(fieldIDX).Interface()
			}
		}
	}
	return theMap
}

// Convert the struct object to a JSON object, then re-parse it as a generic map to interface.
// Merge the interface maps.
// Convert the interface maps to JSON, then back to the struct objects and return.
func mergeVoyageStructs(vInto, vFrom voyage) (voyage, error) {
	vOut := voyage{}
	mapInto := toInterfaceMap(vInto)
	mapFrom := toInterfaceMap(vFrom)

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
	mapInto := toInterfaceMap(vInto)
	mapFrom := toInterfaceMap(vFrom)

	mergeMaps(mapInto, mapFrom) //NB: recursive (bounded by size of map)

	if bInto, err := json.Marshal(mapInto); err != nil {
		return assist{}, JSON_MARSHAL.Errorf("merge structs out marshal")
	} else if err := json.Unmarshal(bInto, &vOut); err != nil {
		return assist{}, JSON_MARSHAL.Errorf("merge structs out unmarshal")
	}

	return vOut, nil
}
