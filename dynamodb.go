package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

type Item struct {
	Year    int
	Title   string
	Plot    string
	Rating  float64
	Content string
}

type Dynamodb struct {
	svc  *dynamodb.DynamoDB
	item []Item
}

/*
var wikiDir = "/home/sh/vimwiki"

func main() {
	conn := Dynamodb{}
	conn.initDB()
	tableName, err := conn.getTable()
	conn.putTags(tableName)
	if err != nil {
		log.Fatal(err)
	}
	conn.getItem(tableName, "#### 3/13", "0")
}
*/

func (conn *Dynamodb) initDB() error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	conn.svc = dynamodb.New(sess)

	return nil
}

func (conn *Dynamodb) getTable() (string, error) {
	input := &dynamodb.ListTablesInput{}
	tableName := ""
	for {
		result, err := conn.svc.ListTables(input)
		if err != nil {
			log.Fatal(err)
			return "error", err
		}
		for _, n := range result.TableNames {
			fmt.Println(*n)
			tableName = *n
			return tableName, nil
		}

		input.ExclusiveStartTableName = result.LastEvaluatedTableName
		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return tableName, nil
}

func (conn *Dynamodb) putItem(tableName string, item Item) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatal(err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = conn.svc.PutItem(input)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Succeessfully added : " + item.Title)
	return nil
}

func (conn *Dynamodb) getItem(tableName string, key string, point string) Item {
	result, err := conn.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(key),
			},
			"Rating": {
				N: aws.String(point),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	item := Item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get: ", item.Title)

	return item
}

func (conn *Dynamodb) deleteItem(tableName string, item Item) error {
	rate := fmt.Sprintf("%f", item.Rating)
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(item.Title),
			},
			"Rating": {
				N: aws.String(rate),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := conn.svc.DeleteItem(input)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Deleted : " + item.Title)
	return nil
}

func (conn *Dynamodb) loadData(tableName string, filename string) error {
	jsonData, err := os.Open(filename)
	defer jsonData.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return err
	}

	var items []Item
	err = json.NewDecoder(jsonData).Decode(&items)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return err
	}

	for _, item := range items {
		info, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Fatal(err)
			return err
		}
		input := &dynamodb.PutItemInput{
			Item:      info,
			TableName: aws.String(tableName),
		}

		_, err = conn.svc.PutItem(input)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	log.Println("Load Json Complete")

	return nil
}

func (conn *Dynamodb) putTags(tableName string) error {
	values, err := getTagAll()
	if err != nil {
		return err
	}
	for _, tagline := range values {
		//item := Item{
		//	Year:    2020,
		//	Title:   tag,
		//	Content: tagline,
		//	Rating:  0.0,
		//}
		//fmt.Println("Title: ", tag, " Content: ", tagline)
		if len(tagline) == 0 {
			continue
		}
		//fmt.Printf("Title: %s", tagline)
		//conn.putItem(tableName, item)
	}
	return nil
}
