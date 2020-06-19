package elasticsearch

import (
	"context"
	"log"

	"github.com/olivere/elastic/v7"
	db "github.com/shdkej/note-server/data_source"
)

type Elastic struct {
	client *elastic.Client
}

func (e *Elastic) Init() error {
	var err error
	e.client, err = elastic.NewClient()
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = e.client.IndexExists("note").Do(context.Background())
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (e *Elastic) Put(tag db.Tag) error {
	_, err := e.client.Index().
		Index("note").
		BodyJson(tag).
		Do(context.Background())

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (e *Elastic) Get() error {
	return nil
}

func (e *Elastic) Update() error {
	return nil
}
