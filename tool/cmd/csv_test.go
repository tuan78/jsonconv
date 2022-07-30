package cmd

import (
	"strings"
	"testing"
)

func TestProcessCsvCmd_NoInputData(t *testing.T) {
	// Prepare
	in := &csvCmdInput{}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)

	// Check
	msg := "need to input either raw data, input file path or data from stdin"
	if err == nil || err.Error() != msg {
		t.Fatalf("It should throw an error with message: %s", msg)
	}
}

func TestProcessCsvCmd_InvalidJsonObject(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "{ \"invalid\": true",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)

	// Check
	msg := "invalid JSON data"
	if err == nil || !strings.HasPrefix(err.Error(), msg) {
		t.Fatalf("It should throw an error starts with message: %s\ncurrent: %v", msg, err)
	}
}

func TestProcessCsvCmd_InvalidJsonArray(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[ { \"invalid\": true }",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)

	// Check
	msg := "invalid JSON data"
	if err == nil || !strings.HasPrefix(err.Error(), msg) {
		t.Fatalf("It should throw an error starts with message: %s\ncurrent: %v", msg, err)
	}
}

func TestProcessCsvCmd_InvalidJsonArrayType(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[ { \"a\": [] }, 1, 2, 3 ]",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)

	// Check
	msg := "unsupport type of JSON data"
	if err == nil || err.Error() != msg {
		t.Fatalf("It should throw an error with message: %s\ncurrent: %v", msg, err)
	}
}

func TestProcessCsvCmd_EmptyCsvData_WhenEmptyJsonObject(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "{}",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}

	// Check
	msg := logger.msg
	expMsg := `
`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_EmptyCsvData_WhenEmptyJsonArray(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[]",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}

	// Check
	msg := logger.msg
	expMsg := `
`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_EmptyCsvData_WhenEmptyJsonObjectInJsonArray(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[{}]",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}

	// Check
	msg := logger.msg
	expMsg := `
`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_JsonObject(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "{ \"a\": true }",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}

	// Check
	msg := logger.msg
	expMsg := `a
true

`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessCsvCmd_JsonArray(t *testing.T) {
	// Prepare
	in := &csvCmdInput{
		raw: "[{ \"a\": true }]",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processCsvCmd(logger, repo, in)
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}

	// Check
	msg := logger.msg
	expMsg := `a
true

`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestOutputCsvContent(t *testing.T) {

}
