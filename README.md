# republik-feeder

![build](https://github.com/maetthu/republik-feeder/workflows/build/badge.svg) 
![release](https://github.com/maetthu/republik-feeder/workflows/release/badge.svg)
[![](https://images.microbadger.com/badges/version/maetthu/republik-feeder.svg)](https://microbadger.com/images/maetthu/republik-feeder "Get your own version badge on microbadger.com")
[![](https://images.microbadger.com/badges/image/maetthu/republik-feeder.svg)](https://microbadger.com/images/maetthu/republik-feeder "Get your own image badge on microbadger.com")

Quick and simple RSS service for [republik.ch](https://www.republik.ch) content. Listens for HTTP requests and returns an RSS feed of the most recent articles.

```
$ http http://localhost:8080/
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

## Installation

Environment:

* REPUBLIK_FEEDER_COOKIE: Contents of the connect.sid cookie after login in browser.

### Manual

* [Download tarball from release page](https://github.com/maetthu/republik-feeder/releases)
* Run

``` 
$ export REPUBLIK_FEEDER_COOKIE="..."
$ ./republik-feeder :8080
```

### Docker

```
$ export REPUBLIK_FEEDER_COOKIE="..."
$ docker run -p 8080:8080 -e REPUBLIK_FEEDER_COOKIE=$REPUBLIK_FEEDER_COOKIE maetthu/republik-feeder:latest
```

### Docker compose

```
version: '3'
services:
  republik-feeder:
    image: "maetthu/republik-feeder:latest"
    user: "65534:65534"
    ports:
      - "8080:8080"
    environment:
      REPUBLIK_FEEDER_COOKIE: "..."
```
