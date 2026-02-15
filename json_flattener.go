package jsonconv

import (
	"fmt"
	"reflect"
)

const (
	// FlattenLevelUnlimited can be set to FlattenOption.Level for unlimited flattening.
	FlattenLevelUnlimited = -1

	// FlattenLevelNonNested can be set to FlattenOption.Level for non-nested flattening
	// (equivalent to non-flattening).
	FlattenLevelNonNested = 0
)

const (
	// DefaultFlattenLevel can be set to FlattenOption.Level for default level flattening
	// (equivalent to FlattenLevelUnlimited).
	DefaultFlattenLevel = FlattenLevelUnlimited

	// DefaultFlattenGap can be set to FlattenOption.Gap for default gap flattening.
	DefaultFlattenGap = "__"
)

// A FlattenOption is for JSON object flattening.
type FlattenOption struct {
	// Level of flattening, it can be FlattenLevelUnlimited,
	// FlattenLevelNonNested or an int value in [1..n]
	Level int

	// A gap between nested JSON and its parent JSON.
	// It will be used when merging nested JSON's key with parent JSON's key
	Gap string

	// Skip Map type (typically JSON Object type) from flattening process
	SkipMap bool

	// Skip Array type (JSON array, string array, int array, float array, etc.)
	// from flattening process
	SkipArray bool
}

// DefaultFlattenOption provides default settings for flattening operations.
var DefaultFlattenOption = &FlattenOption{
	Level: DefaultFlattenLevel,
	Gap:   DefaultFlattenGap,
}

// Flatten flattens obj with given opt. If opt is nil,
// it will use opt value from DefaultFlattenOption instead.
func Flatten(obj map[string]any, opt *FlattenOption) {
	if opt == nil {
		opt = DefaultFlattenOption
	}

	kset := make(map[string]struct{})
	ks := make([]string, 0)
	for k := range obj {
		ks = append(ks, k)
	}
	for _, k := range ks {
		curLvl := 0
		val := reflect.ValueOf(obj[k])
		extract(k, &val, obj, kset, opt, curLvl)
	}
	for k := range kset {
		delete(obj, k)
	}
}

// extract processes obj extraction with k, refval pairs and given opt.
// When extracting map, slice and array, a new key will be stored in kset and curLvl will be increased.
func extract(k string, refval *reflect.Value, obj map[string]any, kset map[string]struct{}, opt *FlattenOption, curLvl int) {
	more := opt.Level == FlattenLevelUnlimited || opt.Level > curLvl
	for refval.Kind() == reflect.Interface {
		*refval = refval.Elem()
	}
	switch refval.Kind() {
	case reflect.Map:
		if !more || opt.SkipMap {
			obj[k] = refval.Interface()
			return
		}
		kset[k] = struct{}{}
		ks := refval.MapKeys()
		for _, nk := range ks {
			newK := fmt.Sprintf("%s%s%s", k, opt.Gap, nk.String())
			nv := refval.MapIndex(nk)
			extract(newK, &nv, obj, kset, opt, curLvl+1)
		}
	case reflect.Slice, reflect.Array:
		if !more || opt.SkipArray {
			obj[k] = refval.Interface()
			return
		}
		kset[k] = struct{}{}
		length := refval.Len()
		for i := 0; i < length; i++ {
			nv := refval.Index(i)
			newK := fmt.Sprintf("%s[%v]", k, i)
			extract(newK, &nv, obj, kset, opt, curLvl+1)
		}
	case reflect.Invalid:
		obj[k] = nil
	default:
		obj[k] = refval.Interface()
	}
}
