package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

const DUMMY_SERVICE_ACCOUNT = ""


func TestMissingRequiredInputExits(t *testing.T) {
	// testing os.Exit called from dependency package requires running
	// the command in a sub-process. Doing this inside the test allows
	// execution to be controlled and results to be captured
	if os.Getenv("CAUSE_EXIT") == "1" {
		main()
		return
	}

	type errorTestCases struct {
		description    string
		requiredEnvVar string
		expectedError  string
	}

	for _, scenario := range []errorTestCases{
		{
			description: "missing filename",
			requiredEnvVar: "INPUT_FILENAME",
			expectedError: "::error::missing input 'filename'",
		},
		{
			description: "missing folderId",
			requiredEnvVar: "INPUT_FOLDERID",
			expectedError: "::error::missing input 'folderId'",
		},
		{
			description: "missing credentials",
			requiredEnvVar: "INPUT_CREDENTIALS",
			expectedError: "::error::missing input 'credentials'",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			// GIVEN
			cmd := exec.Command(os.Args[0], "-test.run=^TestMissingRequiredInputExits$")
			envVariables := []string{
				// special toggle var for this test
				"CAUSE_EXIT=1",
				// all required variables
				"INPUT_FILENAME=foo",
				"INPUT_FOLDERID=foo",
				"INPUT_CREDENTIALS=foo",
			}

			// Find and remove the required var
			for i, v := range envVariables {
				if strings.Contains(v, scenario.requiredEnvVar) {
					envVariables = append(envVariables[:i], envVariables[i+1:]...)
					break
				}
			}

			cmd.Env = append(os.Environ(), envVariables...)

			// WHEN
			stdout, err := cmd.Output()
			state, ok := err.(*exec.ExitError)

			// THEN
			if !ok {
				t.Fatalf("process ran with err %v, want exit status 1", err)
				return
			}
			if state.Success() {
				t.Fatalf("process ran with err %v, want exit status 1", err)
			    return
			}
			result := strings.TrimSpace(string(stdout))
			if result != scenario.expectedError {
				t.Fatalf("unexpected result %v", result)
			}
		})
	}
}
