package jsonconv

import (
	"reflect"

	"github.com/tuan78/jsonconv/utils"
)

func FlattenJsonObject(jsonObj JsonObject, flattenLevel int) {
	nestedLevel := 0
	needLoop := true
	for needLoop {
		needLoop = false // To exit the loop.

		// Find and extract nested json object.
		extractedKeys := make([]string, 0)
		flattenedKV := make(JsonObject)
		for key, val := range jsonObj {
			switch nested := val.(type) {
			case JsonObject:
				for nkey, nval := range nested {
					// Append json object's key with nested json object's key.
					csvKey := key + "__" + nkey
					flattenedKV[csvKey] = nval

					// Check if nested value is kind of map.
					if reflect.ValueOf(nval).Kind() == reflect.Map {
						if flattenLevel == -1 || (flattenLevel > nestedLevel) {
							needLoop = true // Need one more loop.
						}
					}
				}
				// Store extracted nested json object's key for later uses.
				extractedKeys = append(extractedKeys, key)
			}
		}

		// Update json object with its flattened children.
		utils.CopyMap(flattenedKV, jsonObj)
		for _, key := range extractedKeys {
			delete(jsonObj, key)
		}

		nestedLevel++
	}
}
