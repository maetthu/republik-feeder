package main

import (
	"fmt"
	"github.com/eduncan911/podcast"
	"github.com/patrickmn/go-cache"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
	"github.com/maetthu/republik-feeder/internal/lib/client"
)

const baseURL = "https://www.republik.ch"

var sizeCache *cache.Cache

func getSize(URL string) (int64, error) {
	if size, found := sizeCache.Get(URL); found {
		return size.(int64), nil
	}

	r, err := http.Head(URL)

	if err != nil {
		return 0, err
	}

	sizeCache.Set(URL, r.ContentLength, cache.DefaultExpiration)
	return r.ContentLength, nil
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(err.Error()))
}

func articlesHandler(w http.ResponseWriter, r *http.Request) {
	feed := &feeds.Feed{
		Title:       "Republik - RSS Feed",
		Link:        &feeds.Link{Href: baseURL},
		Description: "Republik - RSS Feed",
		Author:      &feeds.Author{Name: "Republik", Email: "kontakt@republik.ch"},
		Created:     time.Now(),
	}

	c := client.NewClient(os.Getenv("REPUBLIK_FEEDER_COOKIE"))
	docs, err := c.Fetch(client.Filter{Feed: true}, 20)

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

func podcastHandler(w http.ResponseWriter, r *http.Request) {
	date := time.Now()

	feed := podcast.New(
		"Republik - Podcast Feed",
		baseURL,
		"Republik - Podcast Feed",
		&date,
		&date,
	)

	feed.Language = "de-CH"

	c := client.NewClient(os.Getenv("REPUBLIK_FEEDER_COOKIE"))
	docs, err := c.Fetch(client.Filter{Feed: true, HasAudio: true, AudioSourceKind: "readAloud"}, 20)

	if err != nil {
		writeError(w, err)
		return
	}

	for _, d := range docs {
		title := d.Meta.Title

		if d.Meta.Format.Meta.Title != "" {
			title = fmt.Sprintf("[%s] %s", d.Meta.Format.Meta.Title, d.Meta.Title)
		}

		pubdate := d.PubDate()

		item := podcast.Item{
			Title:       title,
			Link:        baseURL + d.Meta.Path,
			Description: d.Meta.Description,
			PubDate:     &pubdate,
			IDuration:   strconv.Itoa(int(d.Meta.AudioSource.DurationMs / 1000)),
		}

		// fetch file size
		if size, err := getSize(d.Meta.AudioSource.MP3); err == nil {
			item.AddEnclosure(d.Meta.AudioSource.MP3, podcast.MP3, size)
			_, _ = feed.AddItem(item)
		}
	}

	if err := feed.Encode(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <listen-address>\n", os.Args[0])
		os.Exit(1)
	}

	sizeCache = cache.New(24*time.Hour, 24*time.Hour)

	http.HandleFunc("/articles", articlesHandler)
	http.HandleFunc("/podcast", podcastHandler)

	log.Fatal(http.ListenAndServe(os.Args[1], nil))
}
