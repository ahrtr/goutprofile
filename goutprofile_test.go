package goutprofile_test

import (
	"testing"

	gup "github.com/ahrtr/goutprofile"
)

func TestValidateFile1(t *testing.T) {
	validateFileImpl(t, "./testdata/demo1_utprofile.yml")
}

func TestValidateFile2(t *testing.T) {
	validateFileImpl(t, "./testdata/demo2_utprofile.yml")
}

func TestValidateDirWithoutRecursion(t *testing.T) {
	validateDirImpl(t, "./testdata", false, 2)
}

func TestValidateDirWithRecursion(t *testing.T) {
	validateDirImpl(t, "./testdata", true, 3)
}

func validateFileImpl(t *testing.T, fileName string) {
	if err := gup.ValidateFile(fileName); err != nil {
		t.Errorf("Failed to validate profile: %s, error: %v", fileName, err)
	}
}

func validateDirImpl(t *testing.T, dir string, recursive bool, expectedNum int) {
	if profiles, err := gup.ValidateDir(dir, recursive); err != nil {
		t.Errorf("Failed to validate dir: %s, recursive: %t, error: %v", dir, recursive, err)
	} else {
		if len(profiles) != expectedNum {
			t.Errorf("Expected validating %d profiles, but actually validated %d profiles, dir: %s, recursive: %t",
				expectedNum, len(profiles), dir, recursive)
		}
	}
}
