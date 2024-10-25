package tf

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteResourceFromFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	origFile := filepath.Join(cwd, "test_data", "main.tf")

	testDir, err := os.MkdirTemp("", "testDeleteResourceFromFile")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	srcBytes, err := os.ReadFile(origFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	testFile := filepath.Join(testDir, "main.tf")
	if err = os.WriteFile(testFile, srcBytes, 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	if err := DeleteResourceFromFile(testFile, "google_secret_manager_secret_iam_member", "rube_dev_access"); err != nil {
		t.Fatalf("DeleteResourceFromFile() = %v, wanted nil", err)
	}

	actual, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	t.Logf("Actual file content:\n%s", string(actual))

}
