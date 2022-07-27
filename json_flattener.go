package jsonconv

import (
	"fmt"
	"reflect"
)

const (
	FlattenLevelUnlimited = -1 // Set it to FlattenOption.Level for unlimited flattening
	FlattenLevelNonNested = 0  // Set it to FlattenOption.Level for non-nested flattening (equivalent to non-flattening)
)

// A FlattenOption is for JSON object flattening.
type FlattenOption struct {
	Level     int    // Level of flattening, it can be FlattenLevelUnlimited, FlattenLevelNonNested or an int value in [1..n]
	Gap       string // A gap between nested JSON and its parent JSON. It will be used when merging nested JSON's key with parent JSON's key
	SkipMap   bool   // Skip Map type (typically JSON Object type) from flattening process
	SkipArray bool   // Skip Array type (JSON array, string array, int array, float array, etc.) from flattening process
}

// FlattenJsonObject flattens obj with given op.
func FlattenJsonObject(obj JsonObject, op *FlattenOption) {
	kset := make(map[string]struct{})
	ks := make([]string, 0)
	for k := range obj {
		ks = append(ks, k)
	}
	for _, k := range ks {
		curLvl := 0
		val := reflect.ValueOf(obj[k])
		extractJsonObject(k, &val, obj, kset, op, curLvl)
	}
	for k := range kset {
		delete(obj, k)
	}
}

// extractJsonObject processes obj extraction with k, refval pairs and given op.
// When extracting map, slice and array, a new key will be stores in kset and curLvl will be increased.
func extractJsonObject(k string, refval *reflect.Value, obj JsonObject, kset map[string]struct{}, op *FlattenOption, curLvl int) {
	more := op.Level == FlattenLevelUnlimited || op.Level > curLvl
	for refval.Kind() == reflect.Interface {
		*refval = refval.Elem()
	}
	switch refval.Kind() {
	case reflect.Map:
		if !more || op.SkipMap {
			obj[k] = refval.Interface()
			return
		}
		kset[k] = struct{}{}
		ks := refval.MapKeys()
		for _, nk := range ks {
			newK := fmt.Sprintf("%s%s%s", k, op.Gap, nk.String())
			nv := refval.MapIndex(nk)
			extractJsonObject(newK, &nv, obj, kset, op, curLvl+1)
		}
	case reflect.Slice, reflect.Array:
		if !more || op.SkipArray {
			obj[k] = refval.Interface()
			return
		}
		kset[k] = struct{}{}
		len := refval.Len()
		for i := 0; i < len; i++ {
			nv := refval.Index(i)
			newK := fmt.Sprintf("%s[%v]", k, i)
			extractJsonObject(newK, &nv, obj, kset, op, curLvl+1)
		}
	default:
		obj[k] = refval.Interface()
	}
}
