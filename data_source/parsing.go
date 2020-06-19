package data_source

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Tags struct {
	note string
	tag  string
}

const tagPrefix = "##"

var wikiDir = os.Getenv("VIMWIKI")

func getTagList() ([]string, error) {
	result := []string{}
	err := filepath.Walk(wikiDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".md" {
			return nil
		}
		filename := path
		val, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer val.Close()

		scanner := bufio.NewScanner(val)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, tagPrefix) {
				result = append(result, line)
			}
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	log.Println("Get Tag All Done")
	return result, nil
}

func getTagAll() (map[string][]string, error) {
	result := make(map[string][]string)
	err := filepath.Walk(wikiDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".md" {
			return nil
		}
		val, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var tag string
		ss := strings.Split(string(val), "\n##")
		for _, s := range ss {
			tag = getTag(s)
			if strings.HasPrefix(s, tagPrefix) {
				result[tag] = []string{path, s}
			}
		}

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return result, err
	}
	log.Println("Get Tag All Done")
	return result, nil
}

func getTag(tagline string) string {
	tag := strings.Split(tagline, "\n")
	if len(tag) == 0 {
		return "error"
	}

	return tag[0]
}

func getFileAll() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(wikiDir)
	if err != nil {
		log.Fatal(err)
		return make([]os.FileInfo, 0), err
	}
	return files, nil
}
