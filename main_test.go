package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

var requiredVariables = []string{
	"INPUT_FILENAME=foo",
	"INPUT_FOLDERID=foo",
	"INPUT_CREDENTIALS=foo",
}

func TestRequiredInputExits(t *testing.T) {
	// testing os.Exit called from dependency package requires running
	// the command in a sub-process. Doing this inside the test allows
	// execution to be controlled and results to be captured
	if os.Getenv("CAUSE_EXIT") == "1" {
		main()
		return
	}

	type requiredVarCases struct {
		description    string
		requiredEnvVar string
		expectedError  string
	}

	for _, scenario := range []requiredVarCases{
		{
			description:    "missing filename",
			requiredEnvVar: "INPUT_FILENAME",
			expectedError:  "::error::missing input 'filename'",
		},
		{
			description:    "missing folderId",
			requiredEnvVar: "INPUT_FOLDERID",
			expectedError:  "::error::missing input 'folderId'",
		},
		{
			description:    "missing credentials",
			requiredEnvVar: "INPUT_CREDENTIALS",
			expectedError:  "::error::missing input 'credentials'",
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			// GIVEN
			cmd := exec.Command(os.Args[0], "-test.run=^TestRequiredInputExits$")
			// set special toggle var for this test
			envVariables := []string{
				"CAUSE_EXIT=1",
			}
			envVariables = append(envVariables, requiredVariables...)

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
				t.Fatalf("process ran with err %v, want exit status 1, output %v", err, string(stdout))
				return
			}
			if state.Success() {
				t.Fatalf("process ran with err %v, want exit status 1, output %v", err, string(stdout))
				return
			}
			result := strings.TrimSpace(string(stdout))
			if result != scenario.expectedError {
				t.Fatalf("unexpected result %v", result)
			}
		})
	}
}

func TestIncorrectInputExits(t *testing.T) {
	// testing os.Exit called from dependency package requires running
	// the command in a sub-process. Doing this inside the test allows
	// execution to be controlled and results to be captured
	if os.Getenv("CAUSE_EXIT") == "1" {
		main()
		return
	}

	// GIVEN
	cmd := exec.Command(os.Args[0], "-test.run=^TestIncorrectInputExits$")
	cmd.Env = append(os.Environ(), "CAUSE_EXIT=1", "INPUT_ENCODED=foobar")
	cmd.Env = append(cmd.Env, requiredVariables...)

	// WHEN
	stdout, err := cmd.Output()
	state, ok := err.(*exec.ExitError)

	// THEN
	if !ok {
		t.Fatalf("process ran with err %v, want exit status 1, output %v", err, string(stdout))
		return
	}
	if state.Success() {
		t.Fatalf("process ran with err %v, want exit status 1, output %v", err, string(stdout))
		return
	}
	result := strings.Split(strings.TrimSpace(string(stdout)), "\n")
	containsDesiredMessage := false
	for _, output := range result {
		if output == "::error::incorrect input 'encoded' reason: encoded needs to be either empty, `false` or `true`." {
			containsDesiredMessage = true
			break
		}
	}
	if !containsDesiredMessage {
		t.Fatalf("unexpected result %v", result)
	}
}
