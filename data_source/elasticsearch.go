package data_source

import (
	"context"
	"log"

	"github.com/olivere/elastic/v7"
)

type Elastic struct {
	client *elastic.Client
	index  string
}

func (e *Elastic) Init() error {
	var err error
	e.client, err = elastic.NewClient()
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = e.client.IndexExists(e.index).Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (e *Elastic) Ping() error {
	return nil
}

func (e *Elastic) Hits(key string) error {
	return nil
}

func (e *Elastic) GetStruct(key string) error {
	return nil
}

func (e *Elastic) SetStruct(tag Tag) error {
	_, err := e.client.Index().
		Index(e.index).
		BodyJson(tag).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (e *Elastic) Delete() error {
	return nil
}
