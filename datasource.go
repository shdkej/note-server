package main

type Tag struct {
	FileName    string
	FileContent string
	Tag         string
	TagLine     string
}

type DataSource interface {
	Init() error
	ping() error
	hits(string) (int64, error)
	get(string) (string, error)
	set(string, string) error
	append(string, string) error
	getStruct(string) (Tag, error)
	setStruct(Tag) error
	getSet(string) ([]string, error)
	pushSet(string, string) error
	getAllKey(string) ([]string, error)
	getTagParagraph(string) ([]string, error)
	putTags() error
}
