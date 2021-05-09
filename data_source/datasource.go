// data source from RDBMS, NoSQL
package data_source

import (
	"log"
	"os"
	"strings"
	"time"

	parsing "github.com/shdkej/note-server/parsing"
)

// entrypoint use data source function
type DB struct {
	Store  DataSource
	prefix string
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
	GetStruct(string) (map[string]string, error)
	SetStruct(string, Note) error
	Update(string, string, interface{}) error
	Delete(string) error
	Scan(string) ([]string, error)
}

func (v *DB) SetPrefix(prefix string) {
	var table string = os.Getenv("TABLE")
	if prefix != "" {
		table = prefix
	}
	if table == "" {
		table = "tag"
	}
	v.prefix = table + ":"
	log.Println(v.prefix)
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

	v.SetPrefix("")

	value, _ := v.Get("##")
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
			err := v.Put(tag)
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

func (v *DB) GetEverything() ([]Note, error) {
	m, err := v.Store.Scan(v.prefix)
	if err != nil {
		return []Note{}, err
	}

	var notes []Note
	for _, value := range m {
		tag, err := v.Get(value)
		if err != nil {
			log.Println(err)
			return notes, err
		}
		notes = append(notes, tag)
	}

	return notes, nil
}

func (v *DB) Get(title string) (Note, error) {
	if title == "tag:" || title == "" {
		return Note{}, nil
	}

	if !strings.HasPrefix(title, v.prefix) {
		title = v.prefix + title
	}
	m, err := v.Store.GetStruct(title)
	if err != nil {
		return Note{}, err
	}

	tag := Note{}
	for key, value := range m {
		switch key {
		case "Tag":
			tag.Tag = value
		case "TagLine":
			tag.TagLine = value
		case "FileName":
			tag.FileName = value
		case "CreatedAt":
			tag.CreatedAt = value
		case "UpdatedAt":
			tag.UpdatedAt = value
		}
	}
	return tag, nil
}

func (v *DB) Put(value Note) error {
	now := time.Now().Format("2006-01-02")
	value.CreatedAt = now
	value.UpdatedAt = now
	err := v.Store.SetStruct(v.prefix, value)
	if err != nil {
		return err
	}
	log.Println("put complete", value.Tag)
	return nil
}

func (v *DB) Delete(key string) error {
	err := v.Store.Delete(v.prefix + key)
	if err != nil {
		return err
	}
	return nil
}

func (v *DB) Update(key string, tags interface{}) error {
	now := time.Now().Format("2006-01-02")
	err := v.Store.Update(v.prefix+key, "Tag", tags)
	err = v.Store.Update(v.prefix+key, "UpdatedAt", now)
	if err != nil {
		return err
	}
	return nil
}

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
