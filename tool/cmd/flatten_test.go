package cmd

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestProcessFlattenCmd_NoInputData(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	expMsg := "need to input either raw data, input file path or data from stdin"
	if err == nil || err.Error() != expMsg {
		t.Fatalf("It should throw an error with message: %s", expMsg)
	}
}

func TestProcessFlattenCmd_ReadFileError(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
		inputPath: "test.json",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()
	repo.openFileError = fmt.Errorf("mock open file error")

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	expMsg := repo.openFileError.Error()
	if err.Error() != expMsg {
		t.Fatalf("It should throw an error with message: %s", expMsg)
	}
}

func TestProcessFlattenCmd_CreateFileError(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
		raw: `
		{
			"id":        "b042ab5c-ca73-4460-b739-96410ea9d3a6",
			"user":      "Jon Doe",
			"score":     -100,
			"is active": false
		}`,
		outputPath: "test.json",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()
	repo.createFileError = fmt.Errorf("mock create file error")

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	expMsg := repo.createFileError.Error()
	if err.Error() != expMsg {
		t.Fatalf("It should throw an error with message: %s", expMsg)
	}
}

func TestProcessFlattenCmd_InvalidJsonObject(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
		raw: "{ \"invalid\": true",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	expMsg := "invalid JSON data"
	if err == nil || !strings.HasPrefix(err.Error(), expMsg) {
		t.Fatalf("It should throw an error starts with message: %s\ncurrent: %v", expMsg, err)
	}
}

func TestProcessFlattenCmd_InvalidJsonArray(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
		raw: "[ { \"invalid\": true }",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	expMsg := "invalid JSON data"
	if err == nil || !strings.HasPrefix(err.Error(), expMsg) {
		t.Fatalf("It should throw an error starts with message: %s\ncurrent: %v", expMsg, err)
	}
}

func TestProcessFlattenCmd_InvalidJsonArrayType(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
		raw: "[ { \"a\": [] }, 1, 2, 3 ]",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	expMsg := "unsupport type of JSON data"
	if err == nil || err.Error() != expMsg {
		t.Fatalf("It should throw an error with message: %s\ncurrent: %v", expMsg, err)
	}
}

func TestOutputJsonContent_InvalidJson_WhenOutputToConsole(t *testing.T) {
	// Prepare
	invalid := math.Inf(1)
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := outputJsonContent(logger, repo, invalid, "")

	// Check
	if err == nil {
		t.Fatalf("It should throw an error because of unsupported value")
	}
}

func TestOutputJsonContent_InvalidJson_WhenOutputToJsonFile(t *testing.T) {
	// Prepare
	invalid := math.Inf(1)
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := outputJsonContent(logger, repo, invalid, "test.json")

	// Check
	if err == nil {
		t.Fatalf("It should throw an error because of unsupported value")
	}
}

func TestProcessFlattenCmd_JsonObject(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
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
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process flatten cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := `{"id":"b042ab5c-ca73-4460-b739-96410ea9d3a6","is active":false,"nested__a":1,"nested__b":2,"nested__c__d__e":3,"nested__f[0]":4,"nested__f[1]":5,"nested__f[2]":6,"nested__g__h":"A","nested__g__i":true,"nested__g__j":1,"nested__g__k":1.5,"nested__g__l":"","score":-100,"user":"Jon Doe"}`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessFlattenCmd_JsonArray(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
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
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process flatten cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := `[{"id":"b042ab5c-ca73-4460-b739-96410ea9d3a6","is active":false,"nested__a":1,"nested__b":2,"nested__c__d__e":3,"nested__f[0]":4,"nested__f[1]":5,"nested__f[2]":6,"nested__g__h":"A","nested__g__i":true,"nested__g__j":1,"nested__g__k":1.5,"nested__g__l":"","score":-100,"user":"Jon Doe"},{"id":"ce06f5b1-5721-42c0-91e1-9f72a09c250a","is active":"true","score":"1.5","user":"Tuấn"},{"id":"4e01b638-44e5-4079-8043-baabbff21cc8","is active":"true","score":"100000000000000000000000","user":"高橋"}]`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessFlattenCmd_ReadFromJsonFile(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
		inputPath: "test.json",
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
	err := processFlattenCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process flatten cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := `{"id":"b042ab5c-ca73-4460-b739-96410ea9d3a6","is active":false,"nested__a":1,"nested__b":2,"nested__c__d__e":3,"nested__f[0]":4,"nested__f[1]":5,"nested__f[2]":6,"nested__g__h":"A","nested__g__i":true,"nested__g__j":1,"nested__g__k":1.5,"nested__g__l":"","score":-100,"user":"Jon Doe"}`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessFlattenCmd_ReadFromStdin(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{}
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
	err := processFlattenCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process flatten cmd, err: %v", err)
	}
	msg := strings.TrimSpace(logger.msg)
	expMsg := `{"id":"b042ab5c-ca73-4460-b739-96410ea9d3a6","is active":false,"nested__a":1,"nested__b":2,"nested__c__d__e":3,"nested__f[0]":4,"nested__f[1]":5,"nested__f[2]":6,"nested__g__h":"A","nested__g__i":true,"nested__g__j":1,"nested__g__k":1.5,"nested__g__l":"","score":-100,"user":"Jon Doe"}`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessFlattenCmd_ToJsonFile(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
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
		outputPath: "test.json",
	}
	logger := NewMockLogger()
	repo := NewMockRepository()

	// Process
	err := processFlattenCmd(logger, repo, in)

	// Check
	if err != nil {
		t.Fatalf("failed to process flatten cmd, err: %v", err)
	}
	msg := strings.TrimSpace(repo.writeBuffer.String())
	expMsg := `{"id":"b042ab5c-ca73-4460-b739-96410ea9d3a6","is active":false,"nested__a":1,"nested__b":2,"nested__c__d__e":3,"nested__f[0]":4,"nested__f[1]":5,"nested__f[2]":6,"nested__g__h":"A","nested__g__i":true,"nested__g__j":1,"nested__g__k":1.5,"nested__g__l":"","score":-100,"user":"Jon Doe"}`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}
