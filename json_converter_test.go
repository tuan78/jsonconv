package jsonconv

import (
	"strings"
	"testing"
)

func TestToCsv_EmptyArray(t *testing.T) {
	// Prepare
	data := JsonArray{}

	// Process
	csvData := ToCsv(data, nil)

	// Check
	if len(csvData) != 0 {
		t.Fatalf("It should be safe to put empty JSON array as param")
	}
}

func TestToCsv_NonFlatten(t *testing.T) {
	// Prepare
	data := JsonArray{
		{
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
		},
		{
			"id":        "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
			"user":      "Tuấn",
			"score":     1.5,
			"is active": true,
			"nested": JsonObject{
				"a": 1,
				"b": 2,
			},
		},
		{
			"id":        "4e01b638-44e5-4079-8043-baabbff21cc8",
			"user":      "高橋",
			"score":     100000000000000000,
			"is active": true,
		},
	}

	// Process
	csvData := ToCsv(data, nil)

	// Check
	r1 := strings.Join(csvData[0], ",")
	r2 := strings.Join(csvData[1], ",")
	r3 := strings.Join(csvData[2], ",")
	r4 := strings.Join(csvData[3], ",")
	exp1 := "id,is active,nested,score,special1,special2,special3,special4,special5,special6,user"
	exp2 := "b042ab5c-ca73-4460-b739-96410ea9d3a6,false,,-100,&,<,>,&,<,>,Jon Doe"
	exp3 := "ce06f5b1-5721-42c0-91e1-9f72a09c250a,true,map[a:1 b:2],1.5,,,,,,,Tuấn"
	exp4 := "4e01b638-44e5-4079-8043-baabbff21cc8,true,,100000000000000000,,,,,,,高橋"
	if r1 != exp1 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r1, exp1)
	}
	if r2 != exp2 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r2, exp2)
	}
	if r3 != exp3 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r3, exp3)
	}
	if r4 != exp4 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r4, exp4)
	}
}

func TestToCsv_WithBaseHeader(t *testing.T) {
	// Prepare
	data := JsonArray{
		{
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
		},
		{
			"id":        "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
			"user":      "Tuấn",
			"score":     1.5,
			"is active": true,
			"nested": JsonObject{
				"a": 1,
				"b": 2,
			},
		},
		{
			"id":        "4e01b638-44e5-4079-8043-baabbff21cc8",
			"user":      "高橋",
			"score":     100000000000000000,
			"is active": true,
		},
	}

	// Process
	csvData := ToCsv(data, &ToCsvOption{
		BaseHeaders: []string{"x", "y", "z", "3", "2", "1"},
	})

	// Check
	r1 := strings.Join(csvData[0], ",")
	r2 := strings.Join(csvData[1], ",")
	r3 := strings.Join(csvData[2], ",")
	r4 := strings.Join(csvData[3], ",")
	exp1 := "x,y,z,3,2,1,id,is active,nested,score,special1,special2,special3,special4,special5,special6,user"
	exp2 := ",,,,,,b042ab5c-ca73-4460-b739-96410ea9d3a6,false,,-100,&,<,>,&,<,>,Jon Doe"
	exp3 := ",,,,,,ce06f5b1-5721-42c0-91e1-9f72a09c250a,true,map[a:1 b:2],1.5,,,,,,,Tuấn"
	exp4 := ",,,,,,4e01b638-44e5-4079-8043-baabbff21cc8,true,,100000000000000000,,,,,,,高橋"
	if r1 != exp1 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r1, exp1)
	}
	if r2 != exp2 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r2, exp2)
	}
	if r3 != exp3 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r3, exp3)
	}
	if r4 != exp4 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r4, exp4)
	}
}

func TestToCsv_WithFlattening(t *testing.T) {
	// Prepare
	data := JsonArray{
		{
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
		},
		{
			"id":        "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
			"user":      "Tuấn",
			"score":     1.5,
			"is active": true,
			"nested": JsonObject{
				"a": 1,
				"b": 2,
			},
		},
		{
			"id":        "4e01b638-44e5-4079-8043-baabbff21cc8",
			"user":      "高橋",
			"score":     100000000000000000,
			"is active": true,
		},
	}

	// Process
	csvData := ToCsv(data, &ToCsvOption{
		FlattenOption: &FlattenOption{
			Level: FlattenLevelUnlimited,
			Gap:   "_",
		},
	})

	// Check
	r1 := strings.Join(csvData[0], ",")
	r2 := strings.Join(csvData[1], ",")
	r3 := strings.Join(csvData[2], ",")
	r4 := strings.Join(csvData[3], ",")
	exp1 := "id,is active,nested_a,nested_b,score,special1,special2,special3,special4,special5,special6,user"
	exp2 := "b042ab5c-ca73-4460-b739-96410ea9d3a6,false,,,-100,&,<,>,&,<,>,Jon Doe"
	exp3 := "ce06f5b1-5721-42c0-91e1-9f72a09c250a,true,1,2,1.5,,,,,,,Tuấn"
	exp4 := "4e01b638-44e5-4079-8043-baabbff21cc8,true,,,100000000000000000,,,,,,,高橋"
	if r1 != exp1 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r1, exp1)
	}
	if r2 != exp2 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r2, exp2)
	}
	if r3 != exp3 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r3, exp3)
	}
	if r4 != exp4 {
		t.Fatalf("created headers are incorrect, %s is not equal expected %s", r4, exp4)
	}
}

func TestCreateCsvHeader(t *testing.T) {
	// Prepare
	data := JsonArray{
		{
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
		},
		{
			"id":        "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
			"user":      "Tuấn",
			"score":     1.5,
			"is active": true,
			"nested": JsonObject{
				"a": 1,
				"b": 2,
			},
		},
		{
			"id":        "4e01b638-44e5-4079-8043-baabbff21cc8",
			"user":      "高橋",
			"score":     100000000000000000,
			"is active": true,
		},
	}

	// Process
	headers := CreateCsvHeader(data, []string{"x", "y", "z", "3", "2", "1"})

	// Check
	if len(headers) != 17 {
		t.Fatalf("created headers have wrong length %v", len(headers))
	}
	s := strings.Join(headers, ",")
	expected := `x,y,z,3,2,1,id,is active,nested,score,special1,special2,special3,special4,special5,special6,user`
	if s != expected {
		t.Fatalf("created headers are incorrect")
	}
}
