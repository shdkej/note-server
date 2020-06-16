package data_source

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

type Dynamodb struct {
	svc       *dynamodb.DynamoDB
	item      []Tag
	TableName string
}

func (conn *Dynamodb) Init() error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	conn.svc = dynamodb.New(sess)

	return nil
}

func (conn *Dynamodb) getTable() error {
	input := &dynamodb.ListTablesInput{}
	tableName := ""
	for {
		result, err := conn.svc.ListTables(input)
		if err != nil {
			log.Fatal(err)
			return err
		}
		for _, n := range result.TableNames {
			fmt.Println(*n)
			tableName = *n
			conn.TableName = tableName
			return nil
		}

		input.ExclusiveStartTableName = result.LastEvaluatedTableName
		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return nil
}

func (conn *Dynamodb) put(tag Tag) error {
	av, err := dynamodbattribute.MarshalMap(tag)
	if err != nil {
		log.Fatal(err)
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(conn.TableName),
	}

	_, err = conn.svc.PutItem(input)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Succeessfully added : " + tag.Tag)
	return nil
}

func (conn *Dynamodb) get(key string) Tag {
	result, err := conn.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(conn.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Tag": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	item := Tag{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get: ", item.Tag)

	return item
}

func (conn *Dynamodb) deleteItem(item Tag) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Tag": {
				S: aws.String(item.Tag),
			},
		},
		TableName: aws.String(conn.TableName),
	}
	_, err := conn.svc.DeleteItem(input)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("Deleted : " + item.Tag)
	return nil
}

func (conn *Dynamodb) loadData(filename string) error {
	jsonData, err := os.Open(filename)
	defer jsonData.Close()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return err
	}

	var items []Tag
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
			TableName: aws.String(conn.TableName),
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

func (conn *Dynamodb) putTags() error {
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
		conn.put(tag)
	}
	return nil
}
