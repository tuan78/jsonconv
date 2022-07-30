package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tuan78/jsonconv"
)

func TestProcessCsvCmd_NoInputData(t *testing.T) {
	// Prepare
	in := &csvCmdInput{}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	expMsg := "need to input either raw data, input file path or data from stdin"
	if err == nil || err.Error() != expMsg {
		t.Fatalf("It should throw an error with message: %s", expMsg)
	}
}

func TestProcessCsvCmd_ReadFileError(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		inputPath: "test.json",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()
	repo.openFileError = fmt.Errorf("mock open file error")

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	expMsg := repo.openFileError.Error()
	if err.Error() != expMsg {
		t.Fatalf("It should throw an error with message: %s", expMsg)
	}
}

func TestProcessCsvCmd_CreateFileError(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: `
		{
			"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
			"user":      "Jon Doe",
			"score":     -100,
			"is active": false
		}`,
		outputPath: "test.csv",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()
	repo.createFileError = fmt.Errorf("mock create file error")

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	expMsg := repo.createFileError.Error()
	if err.Error() != expMsg {
		t.Fatalf("It should throw an error with message: %s", expMsg)
	}
}

func TestProcessCsvCmd_InvalidJsonObject(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "{ \"invalid\": true",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	expMsg := "invalid JSON data"
	if err == nil || !strings.HasPrefix(err.Error(), expMsg) {
		t.Fatalf("It should throw an error starts with message: %s\ncurrent: %v", expMsg, err)
	}
}

func TestProcessCsvCmd_InvalidJsonArray(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[ { \"invalid\": true }",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	expMsg := "invalid JSON data"
	if err == nil || !strings.HasPrefix(err.Error(), expMsg) {
		t.Fatalf("It should throw an error starts with message: %s\ncurrent: %v", expMsg, err)
	}
}

func TestProcessCsvCmd_InvalidJsonArrayType(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[ { \"a\": [] }, 1, 2, 3 ]",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	expMsg := "unsupport type of JSON data"
	if err == nil || err.Error() != expMsg {
		t.Fatalf("It should throw an error with message: %s\ncurrent: %v", expMsg, err)
	}
}

func TestProcessCsvCmd_InvalidDelimiter_WhenOutputToConsole(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: `
		{
			"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
			"user":      "Jon Doe",
			"score":     -100,
			"is active": false
		}`,
		baseHs:     []string{"z", "y", "x"},
		delim:      "\n",
		flattenOpt: jsonconv.DefaultFlattenOption,
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	msg := "csv: invalid field or comment delimiter"
	if err == nil || !strings.HasPrefix(err.Error(), msg) {
		t.Fatalf("It should throw an error starts with message: %s\ncurrent: %v", msg, err)
	}
}

func TestProcessCsvCmd_InvalidDelimiter_WhenOutputToCsvFile(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: `
		{
			"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
			"user":      "Jon Doe",
			"score":     -100,
			"is active": false
		}`,
		outputPath: "test.csv",
		baseHs:     []string{"z", "y", "x"},
		delim:      "\n",
		flattenOpt: jsonconv.DefaultFlattenOption,
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	msg := "csv: invalid field or comment delimiter"
	if err == nil || !strings.HasPrefix(err.Error(), msg) {
		t.Fatalf("It should throw an error starts with message: %s\ncurrent: %v", msg, err)
	}
}

func TestProcessCsvCmd_EmptyCsvData_WhenEmptyJsonObject(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "{}",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := ""
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_EmptyCsvData_WhenEmptyJsonArray(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[]",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := ""
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_EmptyCsvData_WhenEmptyJsonObjectInJsonArray(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[{}]",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := ""
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_JsonObject_ToConsole(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: `
		{
			"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
			"user":      "Jon Doe",
			"score":     -100,
			"is active": false,
			"nested": {
				"a": 1,
				"b": 2,
				"c": {
					"d": {
						"e": 3
					}
				},
				"f": [4, 5, 6],
				"g": {
					"h": "A",
					"i": true,
					"j": 1,
					"k": 1.5,
					"l": ""
				}
			}
		}`,
		baseHs:     []string{"z", "y", "x"},
		delim:      "|",
		flattenOpt: jsonconv.DefaultFlattenOption,
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := `z|y|x|id|is active|nested__a|nested__b|nested__c__d__e|nested__f[0]|nested__f[1]|nested__f[2]|nested__g__h|nested__g__i|nested__g__j|nested__g__k|nested__g__l|score|user
|||b042ab5c-ca73-4460-b739-96410ea9d3a6|false|1|2|3|4|5|6|A|true|1|1.5||-100|Jon Doe`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_JsonArray_ToConsole(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: `
		[
			{
				"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
				"user":      "Jon Doe",
				"score":     -100,
				"is active": false,
				"nested": {
					"a": 1,
					"b": 2,
					"c": {
						"d": {
							"e": 3
						}
					},
					"f": [4, 5, 6],
					"g": {
						"h": "A",
						"i": true,
						"j": 1,
						"k": 1.5,
						"l": ""
					}
				}
			},
			{
				"id": "ce06f5b1-5721-42c0-91e1-9f72a09c250a",
				"user": "Tuấn",
				"score": "1.5",
				"is active": "true"
			},
			{
				"id": "4e01b638-44e5-4079-8043-baabbff21cc8",
				"user": "高橋",
				"score": "100000000000000000000000",
				"is active": "true"
			}
		]`,
		baseHs:     []string{"z", "y", "x"},
		delim:      "|",
		flattenOpt: jsonconv.DefaultFlattenOption,
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := `z|y|x|id|is active|nested__a|nested__b|nested__c__d__e|nested__f[0]|nested__f[1]|nested__f[2]|nested__g__h|nested__g__i|nested__g__j|nested__g__k|nested__g__l|score|user
|||b042ab5c-ca73-4460-b739-96410ea9d3a6|false|1|2|3|4|5|6|A|true|1|1.5||-100|Jon Doe
|||ce06f5b1-5721-42c0-91e1-9f72a09c250a|true||||||||||||1.5|Tuấn
|||4e01b638-44e5-4079-8043-baabbff21cc8|true||||||||||||100000000000000000000000|高橋`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_ReadFromJsonFile(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		inputPath:  "test.json",
		baseHs:     []string{"z", "y", "x"},
		delim:      "|",
		flattenOpt: jsonconv.DefaultFlattenOption,
	}
	logger := NewMockLogger()
	repo := NewMockRepository()
	repo.readContent = `
	{
		"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":      "Jon Doe",
		"score":     -100,
		"is active": false,
		"nested": {
			"a": 1,
			"b": 2,
			"c": {
				"d": {
					"e": 3
				}
			},
			"f": [4, 5, 6],
			"g": {
				"h": "A",
				"i": true,
				"j": 1,
				"k": 1.5,
				"l": ""
			}
		}
	}`

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := `z|y|x|id|is active|nested__a|nested__b|nested__c__d__e|nested__f[0]|nested__f[1]|nested__f[2]|nested__g__h|nested__g__i|nested__g__j|nested__g__k|nested__g__l|score|user
|||b042ab5c-ca73-4460-b739-96410ea9d3a6|false|1|2|3|4|5|6|A|true|1|1.5||-100|Jon Doe`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_ReadFromStdin(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		baseHs:     []string{"z", "y", "x"},
		delim:      "|",
		flattenOpt: jsonconv.DefaultFlattenOption,
	}
	logger := NewMockLogger()
	repo := NewMockRepository()
	repo.isStdinEmpty = false // fake stdin data
	repo.readContent = `
	{
		"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
		"user":      "Jon Doe",
		"score":     -100,
		"is active": false,
		"nested": {
			"a": 1,
			"b": 2,
			"c": {
				"d": {
					"e": 3
				}
			},
			"f": [4, 5, 6],
			"g": {
				"h": "A",
				"i": true,
				"j": 1,
				"k": 1.5,
				"l": ""
			}
		}
	}`

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := `z|y|x|id|is active|nested__a|nested__b|nested__c__d__e|nested__f[0]|nested__f[1]|nested__f[2]|nested__g__h|nested__g__i|nested__g__j|nested__g__k|nested__g__l|score|user
|||b042ab5c-ca73-4460-b739-96410ea9d3a6|false|1|2|3|4|5|6|A|true|1|1.5||-100|Jon Doe`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_ToCsvFile(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: `
			{
				"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
				"user":      "Jon Doe",
				"score":     -100,
				"is active": false,
				"nested": {
					"a": 1,
					"b": 2,
					"c": {
						"d": {
							"e": 3
						}
					},
					"f": [4, 5, 6],
					"g": {
						"h": "A",
						"i": true,
						"j": 1,
						"k": 1.5,
						"l": ""
					}
				}
			}`,
		outputPath: "test.csv",
		baseHs:     []string{"z", "y", "x"},
		delim:      "|",
		flattenOpt: jsonconv.DefaultFlattenOption,
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processCsvCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}
	msg := strings.TrimSpace(repo.writeBuffer.String())
	expMsg := `z|y|x|id|is active|nested__a|nested__b|nested__c__d__e|nested__f[0]|nested__f[1]|nested__f[2]|nested__g__h|nested__g__i|nested__g__j|nested__g__k|nested__g__l|score|user
|||b042ab5c-ca73-4460-b739-96410ea9d3a6|false|1|2|3|4|5|6|A|true|1|1.5||-100|Jon Doe`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}
