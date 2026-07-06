package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func createTestFolder(t *testing.T, folderPath string, fileCount int) {
	for i := 1; i <= fileCount; i++ {
		fileName := fmt.Sprintf("comic_book_%02d.cbz", i)
		filePath := filepath.Join(folderPath, fileName)

		err := createEmptyCBZ(filePath)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", fileName, err)
		}
	}

}

func TestCreateCBZFiles(t *testing.T) {
	tmpDirOuter, err := os.MkdirTemp(".", "cbz_test_outer")
	if err != nil {
		t.Fatalf("Failed to create first temp dir: %v", err)
	}
	defer func() {
		err = os.RemoveAll(tmpDirOuter)
		if err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	tmpDirInner, err := os.MkdirTemp(tmpDirOuter, "cbz_test_inner")
	if err != nil {
		t.Fatalf("Failed to create first temp dir: %v", err)
	}
	defer func() {
		err = os.RemoveAll(tmpDirInner)
		if err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	outerSize := rand.Intn(100)
	innerSize := rand.Intn(100)

	createTestFolder(t, tmpDirOuter, outerSize)
	t.Logf("Create test folder %s with %v files", tmpDirOuter, outerSize)
	createTestFolder(t, tmpDirInner, innerSize)
	t.Logf("Create test folder %s with %v files", tmpDirInner, innerSize)

	match := regexp.MustCompile(`comic_book_(\d+)`)
	chanCBZ := FindCBZ(tmpDirOuter, match)

	count := 0
	for _ = range chanCBZ {
		count++
	}
	totalSize := innerSize + outerSize
	if count != totalSize {
		t.Fatalf("Found %v chapters, expected %v", count, totalSize)
	}
}
