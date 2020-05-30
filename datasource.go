package main

type Tag struct {
	fileName    string
	fileContent string
	tag         string
	tagLine     string
}

type datasource interface {
	get(string) (string, error)
	set(string, string) error
	getStruct() (Article, error)
	setStruct() error
	append(string, string) error
	getSet(string) ([]string, error)
	pushSet(string, string) error
	getAllKey(string) ([]string, error)
	getTagParagraph(string) (string, error)
	putTags() error
}
