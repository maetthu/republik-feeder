package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/feeds"
	"github.com/maetthu/republik-feeder/internal/lib/client"
)

const BASE_URL = "https://www.republik.ch"

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "Republik - RSS Feed",
		Link:        &feeds.Link{Href: BASE_URL},
		Description: "Republik - RSS Feed",
		Author:      &feeds.Author{Name: "Republik", Email: "kontakt@republik.ch"},
		Created:     now,
	}

	c := client.NewClient(os.Getenv("REPUBLIK_FEEDER_COOKIE"))
	docs, err := c.Fetch(20)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	for _, d := range docs {
		title := d.Meta.Title

		if d.Meta.Format.Meta.Title != "" {
			title = fmt.Sprintf("[%s] %s", d.Meta.Format.Meta.Title, d.Meta.Title)
		}

		feed.Add(&feeds.Item{
			Id:          d.ID,
			Title:       title,
			Link:        &feeds.Link{Href: BASE_URL + d.Meta.Path},
			Description: d.Meta.Description,
			Created:     d.PubDate(),
		})
	}

	rss, err := feed.ToRss()

	w.Write([]byte(rss))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
