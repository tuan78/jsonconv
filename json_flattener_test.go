package jsonconv

import (
	"testing"
)

func sampleJsonObject() JsonObject {
	return JsonObject{
		"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":      "Jon Doe",
		"score":     -100,
		"is active": false,
		"nested": JsonObject{
			"a": 1,
			"b": 2,
			"c": JsonObject{
				"d": JsonObject{
					"e": 3,
				},
			},
			"f": []int{4, 5, 6},
			"g": JsonObject{
				"h": "A",
				"i": true,
				"j": 1,
				"k": 1.5,
				"l": nil,
			},
		},
	}
}

func TestFlattenJsonObject_NilFlattenOption(t *testing.T) {
	// Prepare
	data := sampleJsonObject()

	// Process
	FlattenJsonObject(data, nil)

	// Check
	expected := JsonObject{
		"id":              "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":            "Jon Doe",
		"score":           -100,
		"is active":       false,
		"nested__a":       1,
		"nested__b":       2,
		"nested__c__d__e": 3,
		"nested__f[0]":    4,
		"nested__f[1]":    5,
		"nested__f[2]":    6,
		"nested__g__h":    "A",
		"nested__g__i":    true,
		"nested__g__j":    1,
		"nested__g__k":    1.5,
		"nested__g__l":    nil,
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", v, ev, k)
		}
	}
}

func TestFlattenJsonObject_UnlimitedLevel(t *testing.T) {
	// Prepare
	data := sampleJsonObject()

	// Process
	FlattenJsonObject(data, &FlattenOption{
		Level: FlattenLevelUnlimited,
		Gap:   "__",
	})

	// Check
	expected := JsonObject{
		"id":              "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":            "Jon Doe",
		"score":           -100,
		"is active":       false,
		"nested__a":       1,
		"nested__b":       2,
		"nested__c__d__e": 3,
		"nested__f[0]":    4,
		"nested__f[1]":    5,
		"nested__f[2]":    6,
		"nested__g__h":    "A",
		"nested__g__i":    true,
		"nested__g__j":    1,
		"nested__g__k":    1.5,
		"nested__g__l":    nil,
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", v, ev, k)
		}
	}
}

func TestFlattenJsonObject_NonNestedLevel(t *testing.T) {
	// Prepare
	data := sampleJsonObject()

	// Process
	FlattenJsonObject(data, &FlattenOption{
		Level: FlattenLevelNonNested,
		Gap:   "|",
	})

	// Check
	expected := JsonObject{
		"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":      "Jon Doe",
		"score":     -100,
		"is active": false,
		"nested": JsonObject{
			"a": 1,
			"b": 2,
			"c": JsonObject{
				"d": JsonObject{
					"e": 3,
				},
			},
			"f": []int{4, 5, 6},
			"g": JsonObject{
				"h": "A",
				"i": true,
				"j": 1,
				"k": 1.5,
				"l": nil,
			},
		},
	}

	// Check flattened values.
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if k != "nested" && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", v, ev, k)
		}
	}

	// Check nested object.
	nes := data["nested"].(JsonObject)
	enes := expected["nested"].(JsonObject)
	if nes["a"] != enes["a"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key a", nes["a"], enes["a"])
	}
	if nes["b"] != enes["b"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key b", nes["b"], enes["b"])
	}

	c := nes["c"].(JsonObject)
	d := c["d"].(JsonObject)
	ec := enes["c"].(JsonObject)
	ed := ec["d"].(JsonObject)
	if d["e"] != ed["e"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key e", d["e"], ed["e"])
	}

	f := nes["f"].([]int)
	ef := enes["f"].([]int)
	for idx := range ef {
		if f[idx] != ef[idx] {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key f", f[idx], ef[idx])
		}
	}

	g := nes["g"].(JsonObject)
	eg := enes["g"].(JsonObject)
	for k, v := range eg {
		if g[k] != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", g[k], v, k)
		}
	}
}

func TestFlattenJsonObject_FirstLevel(t *testing.T) {
	// Prepare
	data := sampleJsonObject()

	// Process
	FlattenJsonObject(data, &FlattenOption{
		Level: 1,
		Gap:   "|",
	})

	// Check
	expected := JsonObject{
		"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":      "Jon Doe",
		"score":     -100,
		"is active": false,
		"nested|a":  1,
		"nested|b":  2,
		"nested|c": JsonObject{
			"d": JsonObject{
				"e": 3,
			},
		},
		"nested|f": []int{4, 5, 6},
		"nested|g": JsonObject{
			"h": "A",
			"i": true,
			"j": 1,
			"k": 1.5,
			"l": nil,
		},
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if k != "nested|c" && k != "nested|f" && k != "nested|g" && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", v, ev, k)
		}
	}

	// Check nested object.
	c := data["nested|c"].(JsonObject)
	d := c["d"].(JsonObject)
	ec := expected["nested|c"].(JsonObject)
	ed := ec["d"].(JsonObject)
	if d["e"] != ed["e"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key e", d["e"], ed["e"])
	}

	f := data["nested|f"].([]int)
	ef := expected["nested|f"].([]int)
	for idx := range ef {
		if f[idx] != ef[idx] {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key nested|f", f[idx], ef[idx])
		}
	}

	g := data["nested|g"].(JsonObject)
	eg := expected["nested|g"].(JsonObject)
	for k, v := range eg {
		if g[k] != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", g[k], v, k)
		}
	}
}

func TestFlattenJsonObject_Ignores_Map(t *testing.T) {
	// Prepare
	data := sampleJsonObject()

	// Process
	FlattenJsonObject(data, &FlattenOption{
		Level:   FlattenLevelUnlimited,
		Gap:     "|",
		SkipMap: true,
	})

	// Check
	expected := JsonObject{
		"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":      "Jon Doe",
		"score":     -100,
		"is active": false,
		"nested": JsonObject{
			"a": 1,
			"b": 2,
			"c": JsonObject{
				"d": JsonObject{
					"e": 3,
				},
			},
			"f": []int{4, 5, 6},
			"g": JsonObject{
				"h": "A",
				"i": true,
				"j": 1,
				"k": 1.5,
				"l": nil,
			},
		},
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if k != "nested" && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", v, ev, k)
		}
	}

	// Check nested object.
	nes := data["nested"].(JsonObject)
	enes := expected["nested"].(JsonObject)
	if nes["a"] != enes["a"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key a", nes["a"], enes["a"])
	}
	if nes["b"] != enes["b"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key b", nes["b"], enes["b"])
	}

	c := nes["c"].(JsonObject)
	d := c["d"].(JsonObject)
	ec := enes["c"].(JsonObject)
	ed := ec["d"].(JsonObject)
	if d["e"] != ed["e"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key e", d["e"], ed["e"])
	}

	f := nes["f"].([]int)
	ef := enes["f"].([]int)
	for idx := range ef {
		if f[idx] != ef[idx] {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key f", f[idx], ef[idx])
		}
	}

	g := nes["g"].(JsonObject)
	eg := enes["g"].(JsonObject)
	for k, v := range eg {
		if g[k] != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", g[k], v, k)
		}
	}
}

func TestFlattenJsonObject_Ignores_Array(t *testing.T) {
	// Prepare
	data := sampleJsonObject()

	// Process
	FlattenJsonObject(data, &FlattenOption{
		Level:     FlattenLevelUnlimited,
		Gap:       "|",
		SkipArray: true,
	})

	// Check
	expected := JsonObject{
		"id":           "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":         "Jon Doe",
		"score":        -100,
		"is active":    false,
		"nested|a":     1,
		"nested|b":     2,
		"nested|c|d|e": 3,
		"nested|f":     []int{4, 5, 6},
		"nested|g|h":   "A",
		"nested|g|i":   true,
		"nested|g|j":   1,
		"nested|g|k":   1.5,
		"nested|g|l":   nil,
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if k != "nested|f" && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", v, ev, k)
		}
	}

	// Check nested object.
	f := data["nested|f"].([]int)
	ef := expected["nested|f"].([]int)
	for idx := range ef {
		if f[idx] != ef[idx] {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key nested|f", f[idx], ef[idx])
		}
	}
}

func TestFlattenJsonObject_Gap(t *testing.T) {
	// Prepare
	data := sampleJsonObject()

	// Process
	FlattenJsonObject(data, &FlattenOption{
		Level: FlattenLevelUnlimited,
		Gap:   "|",
	})

	// Check
	expected := JsonObject{
		"id":           "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":         "Jon Doe",
		"score":        -100,
		"is active":    false,
		"nested|a":     1,
		"nested|b":     2,
		"nested|c|d|e": 3,
		"nested|f[0]":  4,
		"nested|f[1]":  5,
		"nested|f[2]":  6,
		"nested|g|h":   "A",
		"nested|g|i":   true,
		"nested|g|j":   1,
		"nested|g|k":   1.5,
		"nested|g|l":   nil,
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", v, ev, k)
		}
	}
}

func TestFlattenJsonObject_Special(t *testing.T) {
	type specialS struct{}
	type specialI interface{}

	// Prepare
	spI := new(specialI)
	data := JsonObject{
		"special1": "&",
		"special2": "<",
		"special3": ">",
		"special4": "\u0026",
		"special5": "\u003c",
		"special6": "\u003e",
		"a":        1,
		"b":        2,
		"c": JsonObject{
			"d": JsonObject{
				"e": 3,
			},
		},
		"f": []int{4, 5, 6},
		"g": JsonObject{
			"h": "A",
			"i": true,
			"j": 1,
			"k": 1.5,
			"l": nil,
			"m": func() {},
			"n": specialS{},
			"o": spI,
		},
		"nested": JsonObject{
			"special1": "&",
			"special2": "<",
			"special3": ">",
			"special4": "\u0026",
			"special5": "\u003c",
			"special6": "\u003e",
			"a":        1,
			"b":        2,
			"c": JsonObject{
				"d": JsonObject{
					"e": 3,
				},
			},
			"f": []int{4, 5, 6},
			"g": JsonObject{
				"h": "A",
				"i": true,
				"j": 1,
				"k": 1.5,
				"l": nil,
				"m": func() {},
				"n": specialS{},
				"o": spI,
			},
		},
	}

	// Process
	FlattenJsonObject(data, &FlattenOption{
		Level: FlattenLevelUnlimited,
		Gap:   "|",
	})

	// Check
	expected := JsonObject{
		"special1":        "&",
		"special2":        "<",
		"special3":        ">",
		"special4":        "\u0026",
		"special5":        "\u003c",
		"special6":        "\u003e",
		"a":               1,
		"b":               2,
		"c|d|e":           3,
		"f[0]":            4,
		"f[1]":            5,
		"f[2]":            6,
		"g|h":             "A",
		"g|i":             true,
		"g|j":             1,
		"g|k":             1.5,
		"g|l":             nil,
		"g|m":             func() {},
		"g|n":             specialS{},
		"g|o":             spI,
		"nested|special1": "&",
		"nested|special2": "<",
		"nested|special3": ">",
		"nested|special4": "\u0026",
		"nested|special5": "\u003c",
		"nested|special6": "\u003e",
		"nested|a":        1,
		"nested|b":        2,
		"nested|c|d|e":    3,
		"nested|f[0]":     4,
		"nested|f[1]":     5,
		"nested|f[2]":     6,
		"nested|g|h":      "A",
		"nested|g|i":      true,
		"nested|g|j":      1,
		"nested|g|k":      1.5,
		"nested|g|l":      nil,
		"nested|g|m":      func() {},
		"nested|g|n":      specialS{},
		"nested|g|o":      spI,
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		skipFunc := k != "g|m" && k != "nested|g|m"
		if skipFunc && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v for key %s", v, ev, k)
		}
	}
}
