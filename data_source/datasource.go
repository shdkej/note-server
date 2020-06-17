package data_source

type Tag struct {
	FileName    string
	FileContent string
	Tag         string
	TagLine     string
}

type Database struct {
	Source DataSource
}

type DataSource interface {
	Init() error
	Ping() error
	Hits(string) (int64, error)
	GetStruct(string) (Tag, error)
	SetStruct(Tag) error
	GetAllKey(string) ([]string, error)
	GetTagParagraph(string) ([]string, error)
	PutTags() error
}
