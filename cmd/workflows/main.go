package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func main() {

	name := ""
	prompt := &survey.Input{
		Message: "pull request or push",
	}
	survey.AskOne(prompt, &name)

	file := ""
	prompt = &survey.Input{
		Message: "inform a file to save:",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}
	survey.AskOne(prompt, &file)

	options := []string{}
	prompt2 := &survey.MultiSelect{
		Message: "Choose a option:",
		Options: []string{"security", "second", "third"},
	}
	survey.AskOne(prompt2, &options)
	result := setScaffold()
	result += setCondition(name)
	result += setIntegration(options)
	log.Println(result)
}

func setScaffold() string {
	text := `
name: It build with Template engine
on:`
	return text
}

func setCondition(condition string) string {
	var text string
	if condition == "pr" {
		text += `
	pull_request:
		branches: [master]
	`
	} else {
		text += `
	push:
		branches: [master]
	`
	}
	return text
}

func setIntegration(condition []string) string {
	text := `
jobs:
	integration:
`

	for _, c := range condition {
		if "security" == c {
			text += `
		name: security check
		runs-on: ubuntu-latest
`
		}
		if "test" == c {
			text += ""
		}
	}

	return text
}

func GetFile(dir string, format string) (map[string]string, error) {
	result := make(map[string]string)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != format {
			return nil
		}
		val, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name := strings.Trim(info.Name(), ".yml")
		result[name] = string(val)

		return nil
	})

	if err != nil {
		return result, err
	}

	log.Println("Get All Snippet is Done", len(result))
	return result, nil
}

func GetSnippet(name string) ([]string, error) {
	var tag []string
	ss := strings.Split(string(name), "\n")
	for _, s := range ss {
		if strings.HasPrefix(s, "- name:") {
			tag = append(tag, s)
		}
	}
	return tag, nil
}
