package data_source

import (
	"fmt"
	"os"

	"github.com/fatih/structs"
	"github.com/go-redis/redis"
)

type Redis struct {
	redis *redis.Client
}

func (c *Redis) Init() error {
	host := os.Getenv("REDIS_HOST")
	c.redis = redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: "",
		DB:       0,
	})

	c.Ping()

	return nil
}

func (c *Redis) Ping() error {
	pong, err := c.redis.Ping().Result()
	if err != nil {
		return err
	}

	fmt.Println(pong, err)
	return nil
}

func (c *Redis) SetInitial() error {
	err := c.PutTags()
	if err != nil {
		return err
	}

	return nil
}

func (c *Redis) Hits(page string) (int64, error) {
	hitstring := page + ":"
	hits, err := c.redis.Incr(hitstring).Result()
	if err != nil {
		return 0, err
	}
	return hits, nil
}

func (c *Redis) Set(keyword string, value string) error {
	err := c.redis.Set(keyword, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Redis) Append(keyword string, value string) error {
	err := c.redis.Append(keyword, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Redis) Get(keyword string) (string, error) {
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

func (c *Redis) PushSet(keyword string, value string) error {
	err := c.redis.RPush(keyword, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Redis) GetSet(keyword string) ([]string, error) {
	val, err := c.redis.LRange(keyword, 0, 100).Result()
	if err != nil {
		return val, err
	}
	return val, err
}

func (c *Redis) SetStruct(tag Tag) error {
	const objectPrefix string = "tag:"

	articleM := structs.Map(tag)

	err := c.redis.HMSet(objectPrefix+tag.Tag, articleM).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Redis) GetStruct(title string) (Tag, error) {
	const objectPrefix string = "tag:"

	title = objectPrefix + title
	m, err := c.redis.HGetAll(title).Result()
	if err == redis.Nil {
		fmt.Printf("Article does not exist")
	} else if err != nil {
		return Tag{}, err
	}

	tag := Tag{}
	for key, value := range m {
		switch key {
		case "Tag":
			tag.Tag = value
		case "TagLine":
			tag.TagLine = value
		case "FileName":
			tag.FileName = value
		case "FileContent":
			tag.FileContent = value
		}
	}

	return tag, nil
}

func (c *Redis) GetAllKey(tag string) ([]string, error) {
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
		switch v.(type) {
		case string:
			result[i] = v.(string)
		default:
			fmt.Println(v)
		}
	}

	return result, nil
}

func (c *Redis) GetTagParagraph(tag string) ([]string, error) {
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
		list_value, err := c.GetSet(key)
		if err != nil {
			return []string{"Empty"}, err
		}
		result = append(result, list_value...)
	}

	return result, nil
}

func (c *Redis) PutTags() error {
	values, err := getTagAll()
	if err != nil {
		return err
	}
	for key, tagline := range values {
		/*
			tag := Tag{
				FileName:    tagline[0],
				FileContent: "0",
				Tag:         key,
				TagLine:     tagline[1],
			}
		*/
		if len(tagline) == 0 {
			continue
		}
		c.PushSet(key, tagline[0])
	}
	return nil
}
