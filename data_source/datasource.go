package data_source

type Tag struct {
	FileName    string
	FileContent string
	Tag         string
	TagLine     string
}

type DB struct {
	db DataSource
}

type DataSource interface {
	Init() error
	Ping() error
	Hits(string) (int64, error)
	GetStruct(string) (Tag, error)
	SetStruct(Tag) error
	GetTagParagraph(string) ([]string, error)
}

func (v DB) Init() error {
	err := v.db.Init()
	return err
}

func (v DB) Hits(s string) (int64, error) {
	hits, err := v.Hits(s)
	return hits, err
}

func (v DB) PutTags() error {
	values, err := getTagAll()
	if err != nil {
		return err
	}
	for key, tagline := range values {
		tag := Tag{
			FileName:    tagline[0],
			FileContent: "0",
			Tag:         key,
			TagLine:     tagline[1],
		}
		if len(tagline) == 0 {
			continue
		}
		v.db.SetStruct(tag)
	}
	return nil
}
