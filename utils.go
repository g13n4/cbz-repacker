package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func createEmptyCBZ(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}()

	zipWriter := zip.NewWriter(file)
	defer func() {
		err := zipWriter.Close()
		if err != nil {
			log.Printf("Failed to close zip writer: %v", err)
		}
	}()

	return nil
}

func formatIntSlice(prefix string, slice *[]int) string {
	var output string

	comma := ""
	for _, v := range *slice {
		output += fmt.Sprintf("%s%s", comma, strconv.Itoa(v))
		comma = ", "
	}
	return fmt.Sprintf("%s: %s", prefix, output)
}

func extractFileNameInfo(file *zip.File) (string, string) {
	ext := filepath.Ext(file.Name)
	return strings.TrimRight(file.Name, "."+ext), ext
}
