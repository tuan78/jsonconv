package cmd

import "testing"

func TestProcessFlattenCmd_JsonObject(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
		raw: "{ \"a\": true }",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processFlattenCmd(logger, repo, in)
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}

	// Check
	msg := logger.msg
	expMsg := `{"a":true}

`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestProcessFlattenCmd_JsonArray(t *testing.T) {
	// Prepare
	in := &flattenCmdInput{
		raw: "[{ \"a\": true }]",
	}

	// Process
	logger := NewMockCmdLogger()
	repo := NewMockRepository()
	err := processFlattenCmd(logger, repo, in)
	if err != nil {
		t.Fatalf("failed to process CSV cmd, err: %v", err)
	}

	// Check
	msg := logger.msg
	expMsg := `[{"a":true}]

`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestOutputJsonContent(t *testing.T) {

}
