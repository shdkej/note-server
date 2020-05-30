package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/structs"
	"github.com/go-redis/redis"
	"github.com/russross/blackfriday/v2"
)

type Client struct {
	redis *redis.Client
}

func (c *Client) NewClient() {
	host := os.Getenv("REDIS_HOST")
	c.redis = redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "",
		DB:       0,
	})
}

func (c *Client) ping() error {
	pong, err := c.redis.Ping().Result()
	if err != nil {
		return err
	}

	fmt.Println(pong, err)
	return nil
}

func (c *Client) setInitial() error {
	dir := "/home/sh/vimwiki/"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".md" {
			return nil
		}
		filename := path
		val, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		output := blackfriday.Run(val)
		article := Article{
			Title:    info.Name(),
			Category: info.Name(),
			Content:  string(output),
		}
		err = c.setStruct(article)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	//err = c.putTags()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) hits(page string) (int64, error) {
	hitstring := page + ":"
	hits, err := c.redis.Incr(hitstring).Result()
	if err != nil {
		return 0, err
	}
	return hits, nil
}

func (c *Client) set(keyword string, value string) error {
	err := c.redis.Set(keyword, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) append(keyword string, value string) error {
	err := c.redis.Append(keyword, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) get(keyword string) (string, error) {
	key := keyword
	val, err := c.redis.Get(key).Result()
	if err == redis.Nil {
		fmt.Printf("%s does not exist\n", key)
	} else if err != nil {
		return val, err
	} //else {
	//fmt.Printf("%s = %s\n", key, val)
	//}

	return val, nil
}

func (c *Client) rpush(keyword string, value string) error {
	err := c.redis.RPush(keyword, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) getSet(keyword string) ([]string, error) {
	val, err := c.redis.LRange(keyword, 0, 100).Result()
	if err != nil {
		return val, err
	}
	return val, err
}

type Article struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
}

func (c *Client) setStruct(article Article) error {
	const objectPrefix string = "article:"

	articleM := structs.Map(article)

	err := c.redis.HMSet(objectPrefix+article.Title, articleM).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getStruct(title string) (Article, error) {
	const objectPrefix string = "article:"

	title = objectPrefix + title
	m, err := c.redis.HGetAll(title).Result()
	if err == redis.Nil {
		fmt.Printf("Article does not exist")
	} else if err != nil {
		return Article{}, err
	}

	arti := Article{}
	for key, value := range m {
		switch key {
		case "Title":
			arti.Title = value
		case "Category":
			arti.Category = value
		case "Content":
			arti.Content = value
		}
	}

	return arti, nil
}

func (c *Client) getAllKey(tag string) ([]string, error) {
	var cursor uint64
	var keys []string
	tag = "*" + tag
	// Scan keys
	for {
		var err error
		var t []string
		t, cursor, err = c.redis.Scan(cursor, tag, 1000).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, t...)
		if cursor == 0 {
			break
		}
	}
	if len(keys) == 0 {
		return []string{"Empty"}, nil
	}
	value, err := c.redis.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}
	//interface to []string
	result := make([]string, len(value))
	for i, v := range value {
		result[i] = v.(string)
	}

	return result, nil
}

func (c *Client) getTagParagraph(tag string) ([]string, error) {
	var cursor uint64
	var keys []string
	tag = "*" + tag
	// Scan keys
	for {
		var err error
		var t []string
		t, cursor, err = c.redis.Scan(cursor, tag, 1000).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, t...)
		if cursor == 0 {
			break
		}
	}
	if len(keys) == 0 {
		return []string{"Empty"}, nil
	}

	var result []string
	for _, key := range keys {
		list_value, err := c.getSet(key)
		if err != nil {
			return []string{"Empty"}, err
		}
		result = append(result, list_value...)
	}

	return result, nil
}

func (c *Client) putTags() error {
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
		if len(tagline) == 0 {
			continue
		}
		fmt.Printf("Title: %s", tagline)
		//c.rpush(tag, tagline)
	}
	return nil
}
