package data_source

import (
	"log"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/go-redis/redis"
)

type Redis struct {
	redis *redis.Client
}

func (c *Redis) Init() error {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
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

	log.Println(pong, err)
	return nil
}

func (c *Redis) Hits(page string) (int64, error) {
	hitstring := page + ":"
	hits, err := c.redis.Incr(hitstring).Result()
	if err != nil {
		log.Println("hutswing")
		return 0, err
	}
	return hits, nil
}

func (c *Redis) SetStruct(prefix string, tag Note) error {
	articleM := structs.Map(tag)

	err := c.redis.HMSet(prefix+tag.Tag, articleM).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Redis) GetStruct(prefix string, key string) (Note, error) {
	if key == "tag:" || key == "" {
		return Note{}, nil
	}
	if !strings.HasPrefix(key, prefix) {
		key = prefix + key
	}
	m, err := c.redis.HGetAll(key).Result()
	if err == redis.Nil {
		log.Printf("Article does not exist")
	} else if err != nil {
		log.Printf("Some Error Occured, %s", key)
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
		}
	}

	return tag, nil
}

func (c *Redis) GetTags(prefix string) ([]string, error) {
	var cursor uint64
	var keys []string

	// Scan keys
	for {
		var err error
		var t []string
		t, cursor, err = c.redis.Scan(cursor, prefix+"*", 1000).Result()
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
	return keys, nil
}

func (c *Redis) GetSet(keyword string) ([]string, error) {
	val, err := c.redis.LRange(keyword, 0, 100).Result()
	if err != nil {
		return val, err
	}
	return val, err
}

func (c *Redis) Delete(tag Note) error {
	c.redis.Del(tagPrefix + tag.Tag)
	return nil
}

func (c *Redis) Update(key string, field string, value interface{}) error {
	_, err := c.redis.HSet(key, field, value).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *Redis) AddTag(key string, value string) error {
	_, err := c.redis.HSet(tagPrefix+key, "tag", value).Result()
	if err != nil {
		return err
	}
	return nil
}
