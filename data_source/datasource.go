// data source from RDBMS, NoSQL
package data_source

import (
	parsing "github.com/shdkej/note-server/parsing"
	"log"
)

// entrypoint use data source function
type DB struct {
	Store DataSource
}

type Tag struct {
	FileName string
	Tag      string
	TagLine  string
}

// implemented functions to use the data source
type DataSource interface {
	Init() error
	Ping() error
	Hits(string) (int64, error)
	GetStruct(string) (Tag, error)
	SetStruct(Tag) error
	Delete(Tag) error
}

func (v *DB) Init() {
	err := v.Store.Init()
	if err != nil {
		log.Println("init failed", err)
	}
	err = v.Store.Ping()
	if err != nil {
		log.Println("ping failed", err)
	}

	if _, err := v.Store.GetStruct(""); err != nil {
		tags, err := parsing.GetTagAll()
		if err != nil {
			log.Println(err)
		}
		tag, err := ListToTag(tags)
		if err != nil {
			log.Println(err)
		}
		v.PutTags(tag)
	}
}

func (v *DB) Hits(s string) int64 {
	hits, err := v.Store.Hits(s)
	if err != nil {
		log.Fatal(err)
		return hits
	}
	return hits
}

func (c *DB) GetTag(title string) (Tag, error) {
	m, err := c.Store.GetStruct(title)
	if err != nil {
		return Tag{}, err
	}
	return m, nil
}

func (v *DB) PutTags(values []Tag) error {
	for _, tag := range values {
		v.Store.SetStruct(tag)
	}
	return nil
}

/*
// get tag and return paragraph
func (c *DB) GetTagParagraph(tag string) ([]string, error) {
	var keys []string
	tag = "*" + tag

	var result []string
	keys, err := c.Store.GetTagList(tag)
	if err != nil {
		return []string{"error"}, err
	}
	for _, key := range keys {
		list_value, err := c.Store.GetSet(key)
		if err != nil {
			return []string{"Empty"}, err
		}
		result = append(result, list_value...)
	}

	return result, nil
}
*/

func ListToTag(list map[string][]string) ([]Tag, error) {
	var items []Tag
	// "tag" : ["file path", "tagline"]
	for key, value := range list {
		var tag Tag
		if value == nil {
			continue
		}
		if len(value) < 2 {
			tag = Tag{
				Tag:      key,
				FileName: value[0],
			}
		} else {
			tag = Tag{
				FileName: value[0],
				TagLine:  value[1],
				Tag:      key,
			}
		}
		items = append(items, tag)
	}

	return items, nil
}
