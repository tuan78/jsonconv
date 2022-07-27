package jsonconv

// A CsvRow (equivalent to list of string)
type CsvRow = []string

// CsvData contains list of CsvRow.
type CsvData = []CsvRow

// A JsonObject used to hold JSON object data.
type JsonObject = map[string]interface{}

// A JsonArray contains list of JsonObject.
type JsonArray = []JsonObject
