package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Page struct {
	number int
	file   *zip.File
}

func NewPage(file *zip.File) *Page {
	name, _ := extractFileNameInfo(file)
	value, err := strconv.Atoi(name)
	if err != nil {
		panic(fmt.Errorf("unexpected file name. all files inside cbz file should have a name that only consists of integers: %v", err))
	}

	return &Page{number: value, file: file}
}

type Chapter struct {
	name         string
	absolutePath string
	number       int

	hasMatchedName bool

	pages []*Page
}

func NewChapter(entry os.DirEntry, absolutePath string, reg *regexp.Regexp) *Chapter {
	hasMatchedName := reg.MatchString(entry.Name())

	chapter := Chapter{
		name:         entry.Name(),
		absolutePath: absolutePath,
		pages:        make([]*Page, 0, 22),
	}

	if hasMatchedName {
		valList := reg.FindStringSubmatch(chapter.name)
		if len(valList) > 1 {
			chapterNumber := valList[1]
			number, err := strconv.Atoi(chapterNumber)
			if err != nil {
				log.Printf("Error converting chapter (%s) number to int: %s", entry.Name(), chapterNumber)
			} else {
				chapter.number = number
				chapter.hasMatchedName = true

			}
		}
	}

	return &chapter
}

func (c *Chapter) copyPage(zw *zip.Writer, page *zip.File, firstPageNumber int) error {
	var pageName string
	_, ext := extractFileNameInfo(page)
	if ext == "" {
		pageName = fmt.Sprintf("%04d", firstPageNumber)
	} else {
		pageName = fmt.Sprintf("%04d.%s", firstPageNumber, ext)
	}

	pw, err := zw.Create(pageName)
	if err != nil {
		return fmt.Errorf("failed to create page (%s) inside of new archive : %w", pageName, err)
	}
	pageFile, err := page.Open()
	defer func() {
		err := pageFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(pw, pageFile)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", page.Name, err)
	}
	return nil
}

func (c *Chapter) CopyChapter(zw *zip.Writer, firstPageNumber int) (int, error) {
	r, err := zip.OpenReader(c.absolutePath)
	if err != nil {
		return -1, fmt.Errorf("failed to open cbz (%s): %w", c.absolutePath, err)
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			panic(err)
		}
	}(r)

	for _, zPage := range r.File {
		if strings.HasSuffix(zPage.Name, ".xml") {
			continue
		}
		c.pages = append(c.pages, NewPage(zPage))
	}

	c.sortPages()
	for _, page := range c.pages {
		err := c.copyPage(zw, page.file, firstPageNumber)
		if err != nil {
			return -1, err
		}
		firstPageNumber++
	}

	return firstPageNumber, nil
}

func (c *Chapter) sortPages() {
	sort.Slice(c.pages, func(i, j int) bool {

		return c.pages[i].number < c.pages[j].number
	})
}
