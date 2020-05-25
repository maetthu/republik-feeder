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

const baseURL = "https://www.republik.ch"

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	feed := &feeds.Feed{
		Title:       "Republik - RSS Feed",
		Link:        &feeds.Link{Href: baseURL},
		Description: "Republik - RSS Feed",
		Author:      &feeds.Author{Name: "Republik", Email: "kontakt@republik.ch"},
		Created:     now,
	}

	c := client.NewClient(os.Getenv("REPUBLIK_FEEDER_COOKIE"))
	docs, err := c.Fetch(20)

	if err != nil {
		writeError(w, err)
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
			Link:        &feeds.Link{Href: baseURL + d.Meta.Path},
			Description: d.Meta.Description,
			Created:     d.PubDate(),
		})
	}

	rss, err := feed.ToRss()

	if err != nil {
		writeError(w, err)
		return
	}

	_, _ = w.Write([]byte(rss))
}

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <listen-address>\n", os.Args[0])
		os.Exit(1)
	}

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(os.Args[1], nil))
}
