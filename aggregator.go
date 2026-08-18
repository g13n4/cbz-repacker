package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Aggregator struct {
	chapters     []*Chapter
	outputFolder string
	prefix       string

	CBZSize int
}

func NewAggregator(outputFolder, prefix string, cbzSize int) *Aggregator {
	return &Aggregator{
		outputFolder: outputFolder,
		prefix:       prefix,
		CBZSize:      cbzSize,

		chapters: make([]*Chapter, 0),
	}
}

func (a *Aggregator) Add(chapter *Chapter) {
	a.chapters = append(a.chapters, chapter)
}

func (a *Aggregator) SortChapters() {
	sort.Slice(a.chapters, func(i, j int) bool {
		return a.chapters[i].number < a.chapters[j].number
	})
}

func (a *Aggregator) CheckOrder() error {
	chapterMap := make(map[int]int)

	minNumber := a.chapters[0].number
	maxNumber := a.chapters[len(a.chapters)-1].number
	expectedTotalNumber := (maxNumber - minNumber) + 1
	if expectedTotalNumber != len(a.chapters) {
		expectedNumber := minNumber
		missingChapters := make([]int, 0)
		for _, chapter := range a.chapters {
			if chapter.number != expectedNumber {
				missingChapters = append(missingChapters, expectedNumber)
			}
			expectedNumber = chapter.number + 1
			chapterMap[chapter.number]++
		}

		copyChapters := make([]int, 0)
		for k, v := range chapterMap {
			if v > 1 {
				copyChapters = append(copyChapters, k)
			}
		}

		errorText := ""

		if len(missingChapters) > 0 {
			errorText += formatIntSlice("\nFound missing chapters: ", &missingChapters)
		}

		if len(copyChapters) > 0 {
			errorText += formatIntSlice("\nFound repeating chapters: ", &copyChapters)
		}

		if errorText != "" {
			return fmt.Errorf("found missing chapters!:%s", errorText)
		}

		return nil
	}

	return nil
}

func (a *Aggregator) Repack() error {
	for len(a.chapters) > 0 {
		newOffset := min(a.CBZSize, len(a.chapters))
		currentSlice := a.chapters[:newOffset]
		err := a.repackChapters(&currentSlice)
		if err != nil {
			return err
		}
		a.chapters = a.chapters[newOffset:]
	}
	return nil
}

func (a *Aggregator) repackChapters(chapterSlice *[]*Chapter) error {
	startNumber := (*chapterSlice)[0].number
	endNumber := (*chapterSlice)[len(*chapterSlice)-1].number

	fileName := fmt.Sprintf("%s%04d-%04d.cbz", a.prefix, startNumber, endNumber)
	absFilePath := filepath.Join(a.outputFolder, fileName)

	repackedFile, err := os.Create(absFilePath)
	if err != nil {
		return err
	}
	defer func(repackedFile *os.File) {
		err := repackedFile.Close()
		if err != nil {
			panic(err)
		}
	}(repackedFile)

	zipWriter := zip.NewWriter(repackedFile)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			panic(err)
		}
	}(zipWriter)

	startingPageNumber := 1
	successfullyRepacked := make([]string, 0, len(*chapterSlice))
	for idx, chapter := range *chapterSlice {
		startingPageNumber, err = chapter.CopyChapter(zipWriter, startingPageNumber)
		if err != nil {
			return err
		}
		successfullyRepacked[idx] = chapter.name
	}
	log.Printf("successfully created a repacked cbz %s\n%s", fileName, strings.Join(successfullyRepacked, "\n"))

	return nil
}
