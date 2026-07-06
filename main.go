package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
)

func main() {
	var workingFolder, outputFolder, matchPattern *string
	var defaultFolder, defaultOutputFolder, defaultMatchPattern string
	var repackedCBZSize *int

	defaultMatchPattern = `[Cc]hapter[_\-\s\t]*(\d+)`
	defaultFolder = "."

	ignoreCheck := flag.Bool("ignore_order_check", false, "ignore missing or repeating cbz file error")
	ignoreNotMatched := flag.Bool("ignore_not_matched", false, "ignore not matched cbz files")
	matchPattern = flag.String("pattern", defaultMatchPattern, "change default pattern is to be used to identify chapter number")
	workingFolder = flag.String("input", defaultFolder, "change default folder to where cbz files are located")
	outputFolder = flag.String("output", defaultOutputFolder, "change default folder to where repacked cbz files will be created")
	repackedCBZSize = flag.Int("size", 7, "change amount of chapters in one repacked cbz")

	flag.Parse()

	setFolderAbsolute, err := filepath.Abs(*workingFolder)
	if err != nil {
		panic(fmt.Sprintf("can't find absolute path for directory %s:\n%+v", setFolderAbsolute, err))
	}

	r, err := regexp.Compile(*matchPattern)
	if err != nil {
		panic(fmt.Sprintf("can't compile regex: %s\n%+v", *matchPattern, err))
	}

	cbzChan := FindCBZ(setFolderAbsolute, r)

	absoluteOutputFolder, err := filepath.Abs(*outputFolder)
	if err != nil {
		panic(fmt.Sprintf("can't find absolute path for directory %s:\n%+v", setFolderAbsolute, err))
	}

	agg := NewAggregator(absoluteOutputFolder, *repackedCBZSize)

	notMatchedCounter := 0
	for chapter := range cbzChan {
		if chapter.hasMatchedName {
			agg.Add(chapter)
		} else if !(*ignoreNotMatched) {
			log.Printf("%s is not matched (%s)", chapter.name, chapter.absolutePath)
			notMatchedCounter++
		}
	}

	agg.SortChapters()
	err = agg.CheckOrder()
	if err != nil && !*ignoreCheck {
		panic(err)
	}

	err = agg.Repack()
	if err != nil {
		panic(err)
	}
	return
}
