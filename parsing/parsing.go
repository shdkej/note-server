package parsing

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const tagPrefix = "##"

var wikiDir = os.Getenv("VIMWIKI")

func GetTagAll() (map[string][]string, error) {
	if wikiDir == "" {
		wikiDir = "/home/sh/wiki-blog/content"
	}

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
			tag = GetTagByTagline(s)
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

	log.Println("Get All Tag is Done", len(result))
	return result, nil
}

func GetSnippet(dir string, format string) (map[string][]string, error) {
	result := make(map[string][]string)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != format {
			return nil
		}
		val, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name := strings.Trim(info.Name(), ".yml")
		result[name] = []string{path, string(val)}

		return nil
	})

	if err != nil {
		return result, err
	}

	log.Println("Get All Snippet is Done", len(result))
	return result, nil
}

func GetTagByTagline(tagline string) string {
	tag := strings.Split(tagline, "\n")
	if len(tag) == 0 {
		return "error"
	}

	return tag[0]
}

func AppendToDictionary(file string, keyword string, text string) error {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	val, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	ss := strings.Split(string(val), "\n")
	content := ""
	isExist := false
	for _, line := range ss {
		if strings.HasPrefix(line, keyword) {
			line = line + "," + text
			isExist = true
		}
		content += line
	}

	if isExist == false {
		content += keyword + "," + text
	}

	err = WriteToFile(file, content)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func WriteToFile(file, content string) error {
	err := ioutil.WriteFile(file, []byte(content), 0600)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func LoadJson(filename string) ([]interface{}, error) {
	var items []interface{}
	jsonData, err := os.Open(filename)
	defer jsonData.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return items, err
	}

	err = json.NewDecoder(jsonData).Decode(&items)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return items, err
	}

	log.Println("Load Json Complete")
	return items, nil
}
