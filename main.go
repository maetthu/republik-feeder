package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eduncan911/podcast"
	"github.com/gorilla/feeds"
	"github.com/maetthu/republik-feeder/lib/client"
	"github.com/patrickmn/go-cache"
)

const baseURL = "https://www.republik.ch"

//go:embed all:assets
var assets embed.FS

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

func buildURL(r *http.Request, path string) (*url.URL, error) {
	if env := os.Getenv("REPUBLIK_FEEDER_URL"); env != "" {
		u, err := url.Parse(env)

		if err != nil {
			return nil, err
		}

		u.Path = u.Path + path
		return u, nil
	}

	if path == "" {
		path = "/"
	}

	u := url.URL{
		Path:   path,
		Host:   r.Host,
		Scheme: "http",
	}

	return &u, nil
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	img, _ := buildURL(r, "/assets/cover.png")
	feed.AddImage(img.String())

	c := client.NewClient(os.Getenv("REPUBLIK_FEEDER_COOKIE"))
	docs, err := c.Fetch(client.Filter{Feed: true, HasAudio: true, AudioSourceKind: "readAloud"}, 20)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func assetHandler(path string) http.Handler {
	fsys, err := fs.Sub(assets, "assets")

	if err != nil {
		log.Fatal(err)
	}

	filesystem := http.FS(fsys)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, path)
		_, err := filesystem.Open(path)

		if errors.Is(err, os.ErrNotExist) {
			path = fmt.Sprintf("%s.html", path)
		}

		r.URL.Path = path
		http.FileServer(filesystem).ServeHTTP(w, r)
	})
}

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s <listen-address>\n", os.Args[0])
		os.Exit(1)
	}

	sizeCache = cache.New(24*time.Hour, 24*time.Hour)

	http.HandleFunc("/articles", articlesHandler)
	http.HandleFunc("/podcast", podcastHandler)
	http.Handle("/assets/", assetHandler("/assets"))

	log.Fatal(http.ListenAndServe(os.Args[1], nil))
}
