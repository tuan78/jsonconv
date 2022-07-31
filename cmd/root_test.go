package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmd_CsvCmd(t *testing.T) {
	// Prepare
	raw := `
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
	outBuf := &bytes.Buffer{}
	rootCmd := NewRootCmd()
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"csv", "-d", raw})

	// Process
	err := rootCmd.Execute()

	// Check
	if err != nil {
		t.Fatalf("failed to execute csv cmd, err: %v", err)
	}
	msg := strings.TrimSpace(outBuf.String())
	expMsg := `id,is active,nested__a,nested__b,nested__c__d__e,nested__f[0],nested__f[1],nested__f[2],nested__g__h,nested__g__i,nested__g__j,nested__g__k,nested__g__l,score,user
b042ab5c-ca73-4460-b739-96410ea9d3a6,false,1,2,3,4,5,6,A,true,1,1.5,,-100,Jon Doe`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}

func TestRootCmd_FlattenCmd(t *testing.T) {
	// Prepare
	raw := `
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
	outBuf := &bytes.Buffer{}
	rootCmd := NewRootCmd()
	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(outBuf)
	rootCmd.SetArgs([]string{"flatten", "-d", raw})

	// Process
	err := rootCmd.Execute()

	// Check
	if err != nil {
		t.Fatalf("failed to execute flatten cmd, err: %v", err)
	}
	msg := strings.TrimSpace(outBuf.String())
	expMsg := `{"id":"b042ab5c-ca73-4460-b739-96410ea9d3a6","is active":false,"nested__a":1,"nested__b":2,"nested__c__d__e":3,"nested__f[0]":4,"nested__f[1]":5,"nested__f[2]":6,"nested__g__h":"A","nested__g__i":true,"nested__g__j":1,"nested__g__k":1.5,"nested__g__l":"","score":-100,"user":"Jon Doe"}`
	if msg != expMsg {
		t.Fatalf("It should show message: %s\ncurrent: %s", expMsg, msg)
	}
}
