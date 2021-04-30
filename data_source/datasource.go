// data source from RDBMS, NoSQL
package data_source

import (
	"log"
	"time"

	parsing "github.com/shdkej/note-server/parsing"
)

// entrypoint use data source function
type DB struct {
	Store DataSource
}

type Note struct {
	FileName  string
	Tag       string
	TagLine   string
	CreatedAt string
	UpdatedAt string
}

// implemented functions to use the data source
type DataSource interface {
	Init() error
	Ping() error
	Hits(string) (int64, error)
	GetStruct(string, string) (Note, error)
	SetStruct(string, Note) error
	Delete(Note) error
	Update(string, string, interface{}) error
	GetTags(string) ([]string, error)
	AddTag(string, string) error
}

const tagPrefix string = "tag:"

func (v *DB) Init() {
	err := v.Store.Init()
	if err != nil {
		log.Println("init failed", err)
	}
	err = v.Store.Ping()
	if err != nil {
		log.Println("ping failed", err)
	}

	value, _ := v.Store.GetStruct(tagPrefix, "##")
	if (value == Note{}) {
		log.Println("Parsing start ...")
		tags, err := parsing.GetTagAll()
		if err != nil {
			log.Println("ERROR:", err)
		}
		notes, err := ListToNote(tags)
		if err != nil {
			log.Println("ERROR:", err)
		}
		for _, tag := range notes {
			err := v.PutTag(tag)
			if err != nil {
				log.Println("ERROR:", err)
			}
		}
	}
	log.Println("Init Completed")
}

func (v *DB) Hits(s string) int64 {
	hits, err := v.Store.Hits(s)
	if err != nil {
		log.Fatal(err)
		return hits
	}
	return hits
}

func (c *DB) GetTags() ([]Note, error) {
	m, err := c.Store.GetTags(tagPrefix)
	if err != nil {
		return []Note{}, err
	}

	var notes []Note
	for _, value := range m {
		tag, err := c.Store.GetStruct(tagPrefix, value)
		if err != nil {
			log.Println(err)
			return notes, err
		}
		notes = append(notes, tag)
	}

	return notes, nil
}

func (c *DB) GetEverything(prefix string) ([]Note, error) {
	prefix = prefix + ":"
	log.Println(prefix)
	m, err := c.Store.GetTags(prefix)
	if err != nil {
		return []Note{}, err
	}

	var notes []Note
	for _, value := range m {
		tag, err := c.Store.GetStruct(prefix, value)
		if err != nil {
			log.Println(err)
			return notes, err
		}
		notes = append(notes, tag)
	}

	return notes, nil
}

func (v *DB) Get(prefix string, title string) (Note, error) {
	prefix = prefix + ":"
	m, err := v.Store.GetStruct(prefix, title)
	if err != nil {
		return Note{}, err
	}
	return m, nil
}

func (v *DB) GetTag(title string) (Note, error) {
	m, err := v.Store.GetStruct(tagPrefix, title)
	if err != nil {
		return Note{}, err
	}
	return m, nil
}

func (v *DB) Put(prefix string, value Note) error {
	now := time.Now().Format("2006-01-02")
	value.CreatedAt = now
	value.UpdatedAt = now
	prefix = prefix + ":"
	err := v.Store.SetStruct(prefix, value)
	if err != nil {
		return err
	}
	return nil
}

func (v *DB) PutTag(value Note) error {
	now := time.Now().Format("2006-01-02")
	value.CreatedAt = now
	value.UpdatedAt = now

	err := v.Store.SetStruct(tagPrefix, value)
	if err != nil {
		return err
	}
	return nil
}

func (v *DB) DeleteTag(value Note) error {
	err := v.Store.Delete(value)
	if err != nil {
		return err
	}
	return nil
}

func (v *DB) UpdateTag(key string, tags interface{}) error {
	now := time.Now().Format("2006-01-02")
	err := v.Store.Update(key, "Tag", tags)
	err = v.Store.Update(key, "UpdatedAt", now)
	if err != nil {
		return err
	}
	return nil
}

func (v *DB) PutTagForSearch(key string, value string) error {
	err := v.Store.AddTag(key, value)
	if err != nil {
		return err
	}
	return nil
}

/*
// get tag and return paragraph
func (c *DB) GetTagParagraph(tag string) ([]string, error) {
	var keys []string
	tag = "*" + tag

	var result []string
	keys, err := c.Store.GetTags(tag)
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

func ListToNote(list map[string][]string) ([]Note, error) {
	var items []Note
	// "tag" : ["file path", "tagline"]
	for key, value := range list {
		var tag Note
		if value == nil {
			continue
		}
		if len(value) < 2 {
			tag = Note{
				Tag:      key,
				FileName: value[0],
			}
		} else {
			tag = Note{
				FileName: value[0],
				TagLine:  value[1],
				Tag:      key,
			}
		}
		items = append(items, tag)
	}

	return items, nil
}
