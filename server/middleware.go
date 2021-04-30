package server

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type Middleware func(next HandlerFunc) HandlerFunc

func logHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		next(c)

		elapsed := time.Now().Sub(t)
		log.Printf("[%s] %q %v\n",
			c.Request.Method,
			c.Request.URL.String(),
			elapsed,
		)

		if elapsed.Milliseconds() > 200 {
			message := "found note server late response " + c.Request.URL.String()
			err := SendTelegram(message)
			log.Printf("%s, %v", message, err)
		}
	}
}

func recoverHandler(next HandlerFunc) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(c.ResponseWriter,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next(c)
	}
}

func staticHandler(next HandlerFunc) HandlerFunc {
	var (
		dir       = http.Dir(".")
		indexFile = "index.html"
	)
	return func(c *Context) {
		if c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			next(c)
			return
		}

		file := c.Request.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			next(c)
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			next(c)
			return
		}

		if fi.IsDir() {
			if !strings.HasSuffix(c.Request.URL.Path, "/") {
				http.Redirect(c.ResponseWriter, c.Request, c.Request.URL.Path+"/", http.StatusFound)
				return
			}

			file = path.Join(file, indexFile)

			f, err = dir.Open(file)
			if err != nil {
				next(c)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				next(c)
				return
			}
		}

		http.ServeContent(c.ResponseWriter, c.Request, file, fi.ModTime(), f)
	}
}

func SendTelegram(message string) error {
	var BOT_TOKEN = os.Getenv("TELEGRAM_TOKEN")
	var CHAT_ID = os.Getenv("TELEGRAM_CHAT_ID")

	address := "https://api.telegram.org/"
	address += BOT_TOKEN + "/sendmessage?"
	address += "chat_id=" + CHAT_ID + "&"
	address += "text=" + url.QueryEscape(message)

	resp, err := http.Get(address)
	if err != nil {
		log.Println(err)
		return err
	}

	defer resp.Body.Close()

	return nil
}
