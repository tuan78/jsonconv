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
		"special1":  "&",
		"special2":  "<",
		"special3":  ">",
		"special4":  "\u0026",
		"special5":  "\u003c",
		"special6":  "\u003e",
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
			},
		},
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
		"special1":        "&",
		"special2":        "<",
		"special3":        ">",
		"special4":        "\u0026",
		"special5":        "\u003c",
		"special6":        "\u003e",
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
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", v, ev)
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
		"special1":  "&",
		"special2":  "<",
		"special3":  ">",
		"special4":  "\u0026",
		"special5":  "\u003c",
		"special6":  "\u003e",
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
			},
		},
	}

	// Check flattened values.
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if k != "nested" && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", v, ev)
		}
	}

	// Check nested object.
	nes := data["nested"].(JsonObject)
	enes := expected["nested"].(JsonObject)
	if nes["a"] != enes["a"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", nes["a"], enes["a"])
	}
	if nes["b"] != enes["b"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", nes["b"], enes["b"])
	}

	c := nes["c"].(JsonObject)
	d := c["d"].(JsonObject)
	ec := enes["c"].(JsonObject)
	ed := ec["d"].(JsonObject)
	if d["e"] != ed["e"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", d["e"], ed["e"])
	}

	f := nes["f"].([]int)
	ef := enes["f"].([]int)
	for idx := range ef {
		if f[idx] != ef[idx] {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", f[idx], ef[idx])
		}
	}

	g := nes["g"].(JsonObject)
	eg := enes["g"].(JsonObject)
	for k, v := range eg {
		if g[k] != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", g[k], v)
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
		"special1":  "&",
		"special2":  "<",
		"special3":  ">",
		"special4":  "\u0026",
		"special5":  "\u003c",
		"special6":  "\u003e",
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
		},
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if k != "nested|c" && k != "nested|f" && k != "nested|g" && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", v, ev)
		}
	}

	// Check nested object.
	c := data["nested|c"].(JsonObject)
	d := c["d"].(JsonObject)
	ec := expected["nested|c"].(JsonObject)
	ed := ec["d"].(JsonObject)
	if d["e"] != ed["e"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", d["e"], ed["e"])
	}

	f := data["nested|f"].([]int)
	ef := expected["nested|f"].([]int)
	for idx := range ef {
		if f[idx] != ef[idx] {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", f[idx], ef[idx])
		}
	}

	g := data["nested|g"].(JsonObject)
	eg := expected["nested|g"].(JsonObject)
	for k, v := range eg {
		if g[k] != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", g[k], v)
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
		"special1":  "&",
		"special2":  "<",
		"special3":  ">",
		"special4":  "\u0026",
		"special5":  "\u003c",
		"special6":  "\u003e",
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
			},
		},
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if k != "nested" && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", v, ev)
		}
	}

	// Check nested object.
	nes := data["nested"].(JsonObject)
	enes := expected["nested"].(JsonObject)
	if nes["a"] != enes["a"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", nes["a"], enes["a"])
	}
	if nes["b"] != enes["b"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", nes["b"], enes["b"])
	}

	c := nes["c"].(JsonObject)
	d := c["d"].(JsonObject)
	ec := enes["c"].(JsonObject)
	ed := ec["d"].(JsonObject)
	if d["e"] != ed["e"] {
		t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", d["e"], ed["e"])
	}

	f := nes["f"].([]int)
	ef := enes["f"].([]int)
	for idx := range ef {
		if f[idx] != ef[idx] {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", f[idx], ef[idx])
		}
	}

	g := nes["g"].(JsonObject)
	eg := enes["g"].(JsonObject)
	for k, v := range eg {
		if g[k] != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", g[k], v)
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
		"special1":     "&",
		"special2":     "<",
		"special3":     ">",
		"special4":     "\u0026",
		"special5":     "\u003c",
		"special6":     "\u003e",
		"nested|a":     1,
		"nested|b":     2,
		"nested|c|d|e": 3,
		"nested|f":     []int{4, 5, 6},
		"nested|g|h":   "A",
		"nested|g|i":   true,
		"nested|g|j":   1,
		"nested|g|k":   1.5,
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if k != "nested|f" && ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", v, ev)
		}
	}

	// Check nested object.
	f := data["nested|f"].([]int)
	ef := expected["nested|f"].([]int)
	for idx := range ef {
		if f[idx] != ef[idx] {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", f[idx], ef[idx])
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
		"special1":     "&",
		"special2":     "<",
		"special3":     ">",
		"special4":     "\u0026",
		"special5":     "\u003c",
		"special6":     "\u003e",
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
	}
	for k := range expected {
		ev := expected[k]
		v := data[k]
		if ev != v {
			t.Fatalf("flattened JSON object is incorrect, %v is not equal expected value %v", v, ev)
		}
	}
}
