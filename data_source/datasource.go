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
	SetInitial() error
	Hits(string) (int64, error)
	Get(string) (string, error)
	Set(string, string) error
	Append(string, string) error
	GetStruct(string) (Tag, error)
	SetStruct(Tag) error
	GetSet(string) ([]string, error)
	PushSet(string, string) error
	GetAllKey(string) ([]string, error)
	GetTagParagraph(string) ([]string, error)
	PutTags() error
}
