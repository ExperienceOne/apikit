package identifier_test

import (
	"testing"

	"github.com/ExperienceOne/apikit/generator/identifier"
)

func TestMakeOperationsID(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		ExpectedOutput string
	}{
		{
			name:           "Valid operation id",
			input:          "CreateSession",
			ExpectedOutput: "CreateSession",
		},
		{
			name:           "Valid operation id",
			input:          "DeleteViewsSet",
			ExpectedOutput: "DeleteViewsSet",
		},
		{
			name:           "Valid operation id",
			input:          "Delete#Views#Set",
			ExpectedOutput: "DeleteViewsSet",
		},
		{
			name:           "Valid operation id",
			input:          "deleteViewsSet",
			ExpectedOutput: "DeleteViewsSet",
		},
		{
			name:           "Valid operation id",
			input:          "Delete_Views_Set",
			ExpectedOutput: "Delete_Views_Set",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualOutput, err := identifier.ValidateAndCleanOperationsID(testCase.input)
			if err != nil {
				t.Error(err)
				return
			}
			if testCase.ExpectedOutput != actualOutput {
				t.Errorf(`unexpected bool value (Actual: "%v", Expected: "%v")`, actualOutput, testCase.ExpectedOutput)
			}
		})
	}
}

func TestMakeOperationsIDFail(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		errorMessage string
	}{
		{
			name:         "Operation ID with space",
			input:        "Create Session",
			errorMessage: "operations ID isn't a camel case string (operation ID: Create Session)",
		},
		{
			name:         "Operation ID with undscore first",
			input:        "_CreateSession",
			errorMessage: "operations ID cannot start with _ (operation ID: _CreateSession)",
		},
		{
			name:         "Operation ID with dash",
			input:        "Create-Session",
			errorMessage: "operations ID isn't a camel case string (operation ID: Create-Session)",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := identifier.ValidateAndCleanOperationsID(testCase.input)
			if err == nil {
				return
			}

			if err.Error() != testCase.errorMessage {
				t.Errorf(`unexpected error message (Auctal: "%s", Expected: "%s")`, err.Error(), testCase.errorMessage)
				return
			}
		})
	}
}

func TestMakeIdentifier(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid identifier",
			input:    "Create session",
			expected: "CreateSession",
		},
		{
			name:     "Valid identifier",
			input:    "Create-session",
			expected: "CreateSession",
		},
		{
			name:     "Valid identifier",
			input:    "Delete#Views#Set",
			expected: "DeleteViewsSet",
		},
		{
			name:     "Valid identifier",
			input:    "_sort",
			expected: "Sort",
		},
		{
			name:     "Valid identifier",
			input:    "sort_sort",
			expected: "sortSort",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := identifier.MakeIdentifier(testCase.input)
			if testCase.expected != got {
				t.Errorf(`unexpected bool value (Actual: "%v", Expected: "%v")`, got, testCase.expected)
			}
		})
	}
}
