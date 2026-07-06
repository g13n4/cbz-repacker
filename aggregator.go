package main

import (
	"archive/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type Aggregator struct {
	chapters     []*Chapter
	outputFolder string

	CBZSize int
}

func NewAggregator(outputFolder string, cbzSize int) *Aggregator {
	return &Aggregator{
		outputFolder: outputFolder,
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

	fileName := fmt.Sprintf("%04d-%04d.cbz", startNumber, endNumber)
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
	for _, chapter := range *chapterSlice {
		startingPageNumber, err = chapter.CopyChapter(zipWriter, startingPageNumber)
		if err != nil {
			return err
		}
		log.Printf("successfully created a repacked chapter %s in %s", chapter.name, fileName)
	}

	return nil
}
