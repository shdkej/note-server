package main

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Tags struct {
	note string
	tag  string
}

var tagPrefix = "##"

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
		ss := strings.Split(string(val), "\n\n")
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

func getTaglineAll() ([]string, error) {
	result := []string{}
	err := filepath.Walk(wikiDir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".md" {
			return nil
		}
		filename := path
		val, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		ss := strings.Split(string(val), "\n\n")
		for _, s := range ss {
			if strings.HasPrefix(s, tagPrefix) {
				result = append(result, s)
			}
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	log.Println("Get Tagline All Done")
	return result, nil
}

func getFileAll() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(wikiDir)
	if err != nil {
		log.Fatal(err)
		return make([]os.FileInfo, 0), err
	}
	return files, nil
}

func makeTagSet(filename string) (map[string][]string, error) {
	result := make(map[string][]string)

	val, err := ioutil.ReadFile(filename)
	if err != nil {
		return result, err
	}

	var tag string
	ss := strings.Split(string(val), "\n\n")
	for _, s := range ss {
		tag = getTag(s)
		if strings.HasPrefix(s, tagPrefix) {
			result[tag] = []string{filename, s}
		}
	}
	return result, nil
}

func getRandom(arg interface{}) []int {
	var value int
	switch arg.(type) {
	case string:
		value = len(arg.(string))
		log.Println("check")
	case int:
		value = arg.(int)
	default:
		value = 10
	}

	var numbers []int
	for i := 0; i < 5; i++ {
		numbers = append(numbers, rand.Intn(value+i))
	}
	log.Println(numbers)
	return numbers
}

func makeCSVForm(index int, key string, value []string) []string {
	//current := time.Now().String()
	tagline := strings.ReplaceAll(value[1], "\n", " ")
	tagline = "\"" + tagline + "\""
	result := []string{strconv.Itoa(index), tagline}
	// if i saw the tag.
	// change current
	return result
}

func toCSV(tags map[string][]string) error {
	log.Println(tags)
	file, err := os.OpenFile("tags.csv", os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	// TagName Latest-Update Latest-Read Weight
	// How to check update date?
	err = writer.Write([]string{"id", "description"})
	if err != nil {
		log.Fatal(err)
		return err
	}
	var count int
	for key, value := range tags {
		r := makeCSVForm(count, key, value)
		err := writer.Write(r)
		if err != nil {
			log.Fatal(err)
			return err
		}
		count++
	}
	return nil
}
