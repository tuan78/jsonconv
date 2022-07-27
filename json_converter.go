package jsonconv

import (
	"fmt"
	"sort"
)

// A ToCsvOption converts a JSON Array to CSV data.
type ToCsvOption struct {
	FlattenOption *FlattenOption // Set it to apply JSON flattening
	BaseHeaders   CsvRow         // Base CSV headers used to add before dynamic headers
}

// ToCsv converts a JsonArray to CsvData with given op.
func ToCsv(arr JsonArray, op *ToCsvOption) CsvData {
	if len(arr) == 0 {
		return CsvData{}
	}

	// Flatten JSON.
	if op != nil && op.FlattenOption != nil {
		for _, obj := range arr {
			FlattenJsonObject(obj, op.FlattenOption)
		}
	}

	// Create CSV rows.
	var csvData CsvData
	var hs []string
	if op != nil && len(op.BaseHeaders) > 0 {
		hs = CreateCsvHeader(arr, op.BaseHeaders)
	} else {
		hs = CreateCsvHeader(arr, nil)
	}
	csvData = append(csvData, hs)
	for _, obj := range arr {
		row := make(CsvRow, 0)
		for _, h := range hs {
			if val, exist := obj[h]; exist {
				row = append(row, fmt.Sprintf("%v", val))
				continue
			}
			row = append(row, "")
		}
		csvData = append(csvData, row)
	}

	return csvData
}

// CreateCsvHeader creates CsvRow from arr and baseHs.
// A baseHs is base header that we want to put at the beginning of dynamic header,
// we can set baseHs to nil if we just want to have dynamic header only.
func CreateCsvHeader(arr JsonArray, baseHs CsvRow) CsvRow {
	hs := make(sort.StringSlice, 0)
	hss := make(map[string]struct{})

	// Get CSV header from json.
	for _, obj := range arr {
		for k := range obj {
			hss[k] = struct{}{}
		}
	}

	// Exclude base headers from detected headers, then sort filtered list.
	for _, h := range baseHs {
		delete(hss, h)
	}
	for h := range hss {
		hs = append(hs, h)
	}
	hs.Sort()

	// Insert BaseHeaders to the beginning of headers.
	if len(baseHs) > 0 {
		hs = append(baseHs, hs...)
	}
	return hs
}
