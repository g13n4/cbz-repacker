package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func FindCBZ(path string, reg *regexp.Regexp) chan *Chapter {
	wg := sync.WaitGroup{}
	cbzChan := make(chan *Chapter, 10)

	wg.Go(func() {
		parseFolder(path, cbzChan, reg, &wg)
	})

	go func() {
		wg.Wait()
		close(cbzChan)
	}()

	return cbzChan
}

func parseFolder(path string, cbzChan chan *Chapter, reg *regexp.Regexp, wg *sync.WaitGroup) {
	items, err := os.ReadDir(path)
	if err != nil {
		log.Printf("can't parse directory %s:\n%+v", path, err)
	}

	for _, item := range items {
		absolutePath := filepath.Join(path, item.Name())
		if item.IsDir() {
			wg.Go(func() {
				parseFolder(absolutePath, cbzChan, reg, wg)
			})

		}

		if strings.HasSuffix(item.Name(), ".cbz") {
			newCBZChapter := NewChapter(item, absolutePath, reg)
			cbzChan <- newCBZChapter
		}
	}
}
