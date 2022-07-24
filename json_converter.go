package jsonconv

import (
	"fmt"
	"sort"
)

type ConvertInput struct {
	JsonArray    JsonArray
	FlattenLevel int // -1: unlimited, 0: no nested, [1...n]: n level
	BaseHeaders  CsvRow
}

func Convert(input *ConvertInput) (CsvData, error) {
	csvData := make(CsvData, 0)
	if len(input.JsonArray) == 0 {
		return nil, fmt.Errorf("empty JSON array")
	}

	// Flatten JSON object, so can display nested JSON values in CSV columns.
	for _, jsonObject := range input.JsonArray {
		FlattenJsonObject(jsonObject, input.FlattenLevel)
	}

	// Create CSV headers.
	headers := CreateCsvHeader(input.JsonArray, input.BaseHeaders)
	if len(headers) == 0 {
		return nil, fmt.Errorf("empty CSV headers")
	}

	// Create CSV rows.
	csvData = append(csvData, headers)
	for _, jsonObj := range input.JsonArray {
		row := make(CsvRow, 0)
		for _, header := range headers {
			if val, exist := jsonObj[header]; exist {
				row = append(row, fmt.Sprintf("%v", val))
				continue
			}
			row = append(row, "")
		}
		csvData = append(csvData, row)
	}

	return csvData, nil
}

func CreateCsvHeader(jsonArray JsonArray, baseHeaders CsvRow) CsvRow {
	headers := make(sort.StringSlice, 0)
	headerSet := make(map[string]struct{})

	// Get CSV header from json.
	for _, jsonObj := range jsonArray {
		for key := range jsonObj {
			headerSet[key] = struct{}{}
		}
	}

	// Exclude base headers from detected headers, then sort filtered list.
	for _, header := range baseHeaders {
		delete(headerSet, header)
	}
	for header := range headerSet {
		headers = append(headers, header)
	}
	headers.Sort()

	// Insert BaseHeaders to the beginning of headers.
	if len(baseHeaders) > 0 {
		headers = append(baseHeaders, headers...)
	}

	return headers
}
