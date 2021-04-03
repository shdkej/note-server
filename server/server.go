package server

import (
	"log"
	"strings"

	db "github.com/shdkej/note-server/data_source"
)

type Server interface {
	Init()
	RunServer()
	HealthCheck()
}

func RandomHandler() string {
	return "test"
}

func TagHandler(db db.DataSource) string {
	tags, err := db.GetAllKey("####")
	if err != nil {
		log.Fatal(err)
	}
	var taglines string
	taglines += "<h1>Tags</h1>"

	for _, tag := range tags {
		taglines += "<a href='/tag/" + strings.Trim(tag, "# ") + "'>" + tag + "</a><br/><p>" + tag + "</p>"
	}
	if err != nil {
		log.Fatal(err)
	}
	return taglines
}

func TagOneHandler(db db.DataSource, r string) []string {
	hits, err := db.Hits(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(hits)

	content, err := db.GetTagParagraph(r)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func HealthCheck() {

}
