package data_source

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
)

type Tag struct {
	Name      string
	Tags      string
	UpdatedAt string
}

type Elastic struct {
	client *elastic.Client
	index  string
	ctx    context.Context
}

func (e *Elastic) Init() error {
	var err error
	host := os.Getenv("ELASTICSEARCH_HOST")
	if host == "" {
		host = "localhost"
	}
	e.client, err = elastic.NewClient(elastic.SetURL("http://" + host + ":9200"))
	if err != nil {
		log.Fatal(err)
		return err
	}
	e.index = "analyze"
	e.ctx = context.Background()
	exists, err := e.client.IndexExists(e.index).Do(e.ctx)
	if err != nil {
		log.Fatal("check es index exist ", err)
		return err
	}

	if !exists {
		_, err := e.client.CreateIndex(e.index).Do(e.ctx)
		if err != nil {
			log.Fatal("create es index", err)
			return err
		}
	}
	return nil
}

func (e *Elastic) Ping() error {
	return nil
}

func (e *Elastic) Hits(key string) error {
	return nil
}

func (e *Elastic) SetIndex(index string) string {
	e.index = index
	return index
}

func (e *Elastic) GetAll() ([]Tag, error) {
	searchResult, err := e.client.
		Search().
		Index(e.index).
		Pretty(true).
		Size(30).
		Do(e.ctx)

	log.Println(searchResult.TotalHits())
	var note []Tag
	if err != nil {
		log.Println("Get occured Error ", err)
		return note, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var n Tag
		err := json.Unmarshal(hit.Source, &n)
		if err != nil {
			log.Println("Failed Unmarshal", err)
		}

		note = append(note, n)
	}

	return note, nil
}

func (e *Elastic) GetSynonym(key string) ([]Tag, error) {
	query := elastic.NewMultiMatchQuery(key, "Name", "Tags")
	query = query.Analyzer("korean_analyzer")
	searchResult, err := e.client.
		Search().
		Index(e.index).
		Query(query).
		Pretty(true).
		Do(e.ctx)

	log.Println(searchResult.TotalHits())
	var note []Tag
	if err != nil {
		log.Println("Get Synonyms occured Error ", err)
		return note, err
	}

	for _, hit := range searchResult.Hits.Hits {
		var n Tag
		err := json.Unmarshal(hit.Source, &n)
		if err != nil {
			log.Println("Failed Unmarshal", err)
		}

		note = append(note, n)
	}

	return note, nil
}

func (e *Elastic) GetStruct(key string) (string, error) {
	get, err := e.client.Get().Index(e.index).Id(key).Pretty(true).Do(e.ctx)
	if err != nil {
		log.Println("Get occured Error ", err)
		return "", err
	}

	result := string(get.Source)

	return result, nil
}

func (e *Elastic) SetStruct(tag Tag) error {
	now := time.Now().Format("2006-01-02")
	tag.UpdatedAt = now
	_, err := e.client.Index().
		Index(e.index).
		Id(tag.Name).
		BodyJson(tag).
		Do(context.Background())

	if err != nil {
		log.Println("Set occured Error:", err)
		return err
	}
	return nil
}

func (e *Elastic) Update(key string, new_value string) error {
	_, err := e.client.Update().Index(e.index).Id(key).
		Doc(map[string]interface{}{"Tags": new_value}).
		Do(e.ctx)
	if err != nil {
		log.Println("Update Error:", err)
		return err
	}
	return nil
}

func (e *Elastic) Delete(key string) error {
	_, err := e.client.Delete().Index(e.index).Id(key).Do(e.ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
