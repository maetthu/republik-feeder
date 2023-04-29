# republik-feeder

![build](https://github.com/maetthu/republik-feeder/workflows/build/badge.svg) 
![release](https://github.com/maetthu/republik-feeder/workflows/release/badge.svg)

Quick and simple RSS service for [republik.ch](https://www.republik.ch) content. Listens for HTTP requests and returns 
an RSS feed for most recent articles and narrated articles as a podcast feed.

> **Note**
> A paid [republik.ch](https://www.republik.ch) account is required for this tool to work. 
> If the feed returns an empty list, make sure your subscription is already active. 

## Feeds

* /articles the most recent articles

```
$ http http://localhost:8080/articles
HTTP/1.1 200 OK
Content-Type: text/xml; charset=utf-8
Date: Sun, 24 May 2020 20:54:45 GMT
Transfer-Encoding: chunked

<?xml version="1.0" encoding="UTF-8"?><rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/">
  <channel>
    <title>Republik - RSS Feed</title>
    <link>https://www.republik.ch</link>
    <description>Republik - RSS Feed</description>
    <managingEditor>kontakt@republik.ch (Republik)</managingEditor>
    <pubDate>Sun, 24 May 2020 22:54:45 +0200</pubDate>
    <item>
      ...
    </item>
  </channel>
</rss>
```
* /podcast narrated articles as a podcast feed
```
$ http http://localhost:8080/podcast
HTTP/1.1 200 OK
Content-Type: text/xml; charset=utf-8
Date: Fri, 28 Oct 2022 19:26:57 GMT
Transfer-Encoding: chunked

<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" version="2.0">
  <channel>
    <title>Republik - RSS Feed</title>
    <link>https://www.republik.ch</link>
    <description>Republik - Podcast Feed</description>
    <generator>go podcast v1.3.1 (github.com/eduncan911/podcast)</generator>
    <language>de-CH</language>
    <lastBuildDate>Fri, 28 Oct 2022 21:26:56 +0200</lastBuildDate>
    <pubDate>Fri, 28 Oct 2022 21:26:56 +0200</pubDate>
    <item>
        ...
    </item>
  </channel>
</rss>
```

## Installation

* REPUBLIK_FEEDER_COOKIE: Contents of the connect.sid cookie after login in browser.
* REPUBLIK_FEEDER_URL: (optional) Root URL where this service is reachable, e.g https://example.org/feeds/republik  
* REPUBLIK_FEEDER_ARTICLE_LIMIT: (optional) Number of articles in feed (default: 20)  

### Manual

* [Download tarball from release page](https://github.com/maetthu/republik-feeder/releases)
* Run

``` 
$ export REPUBLIK_FEEDER_COOKIE="s%3Asomething..." REPUBLIK_FEEDER_URL="https://example.org/feeds/republik" REPUBLIK_FEEDER_ARTICLE_LIMIT="25"
$ ./republik-feeder :8080
```

### Docker

```
$ export REPUBLIK_FEEDER_COOKIE="s%3Asomething..." REPUBLIK_FEEDER_URL="https://example.org/feeds/republik" REPUBLIK_FEEDER_ARTICLE_LIMIT="25"
$ docker run -p 8080:8080 -e REPUBLIK_FEEDER_COOKIE=$REPUBLIK_FEEDER_COOKIE -e $REPUBLIK_FEEDER_URL=REPUBLIK_FEEDER_URL ghcr.io/maetthu/republik-feeder/republik-feeder:latest
```

### Docker compose

```
version: '3'
services:
  republik-feeder:
    image: "ghcr.io/maetthu/republik-feeder/republik-feeder:latest"
    user: "65534:65534"
    ports:
      - "8080:8080"
    environment:
      REPUBLIK_FEEDER_COOKIE: "s%3Asomething..."
      REPUBLIK_FEEDER_URL: "https://example.org/feeds/republik"
      REPUBLIK_FEEDER_ARTICLE_LIMIT: "25"
```

