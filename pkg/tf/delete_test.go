package tf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDeleteResourceFromFile(t *testing.T) {
	type testCase struct {
		testName string
		source   string
		resource string
		name     string
		expected string
	}

	cases := []testCase{
		{
			testName: "basic",
			source:   "main.tf",
			resource: "google_secret_manager_secret_iam_member",
			name:     "rube_dev_access",
			expected: "main.expected.tf",
		},
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	testDir, err := os.MkdirTemp("", "testDeleteResourceFromFile")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	for _, c := range cases {
		t.Run(c.testName, func(t *testing.T) {
			origFile := filepath.Join(cwd, "test_data", c.source)
			srcBytes, err := os.ReadFile(origFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			testFile := filepath.Join(testDir, c.source)
			if err = os.WriteFile(testFile, srcBytes, 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			if err := DeleteResourceFromFile(testFile, c.resource, c.name); err != nil {
				t.Fatalf("DeleteResourceFromFile failed; %+v", err)
			}

			actual, err := os.ReadFile(testFile)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}

			eFile := filepath.Join(cwd, "test_data", c.expected)
			expected, err := os.ReadFile(eFile)
			if err != nil {
				t.Fatalf("Failed to read expected file: %v", err)
			}
			if d := cmp.Diff(string(expected), string(actual)); d != "" {
				t.Errorf("Mismatch (-expected +actual):\n%s", d)
			}
		})
	}
}
