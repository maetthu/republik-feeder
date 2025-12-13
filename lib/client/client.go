package client

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/machinebox/graphql"
)

const apiURL = "https://api.republik.ch/graphql"

// Document is marshalled from GraphQL response for a single document
type Document struct {
	ID   string
	Meta struct {
		Title                       string
		Path                        string
		Image                       string
		Description                 string
		PublishDate                 string
		Template                    string
		EstimatedReadingMinutes     int
		EstimatedConsumptionMinutes int
		AudioCoverCrop              struct {
			X      int
			Y      int
			Width  int
			Height int
		}
		AudioSource struct {
			MP3        string
			Kind       string
			DurationMs int
		}
		Format struct {
			Meta struct {
				Path  string
				Title string
				Kind  string
			}
		}
	}
}

// PubDate returns parsed publication date of a document
func (d *Document) PubDate() time.Time {
	t, _ := time.Parse(time.RFC3339, d.Meta.PublishDate)
	return t
}

// Response is marshalled from the GraphQL response from search query
type Response struct {
	Documents struct {
		Nodes []struct {
			Entity Document
		}
	}
}

type Filter struct {
	Feed            bool
	HasAudio        bool
	AudioSourceKind string
	Format          string
}

func (f Filter) String() string {
	c := []string{}

	if f.Feed {
		c = append(c, "feed: true")
	}

	if f.HasAudio {
		c = append(c, "hasAudio: true")
	}

	if f.AudioSourceKind != "" {
		c = append(c, "audioSourceKind: "+f.AudioSourceKind)
	}

	if f.Format != "" {
		c = append(c, "format: \""+f.Format+"\"")
	}

	if len(c) == 0 {
		return ""
	}

	return fmt.Sprintf("filter: {%s}", strings.Join(c, ", "))
}

// Client is the wrapper around the republik GraphQL API
type Client struct {
	sid string
}

// Fetch fetches a list of documents
func (c *Client) Fetch(filter Filter, limit int) ([]Document, error) {
	qc := graphql.NewClient(apiURL)
	req := graphql.NewRequest(fmt.Sprintf(`
		query ($limit: Int!) {
			documents: search(
				filters: [
					{ key: "template", not: true, value: "section" }
					{ key: "template", not: true, value: "format" }
					{ key: "template", not: true, value: "front" }
				]
				%s
				sort: { key: publishedAt, direction: DESC }
				first: $limit
			) {
				nodes {
					entity {
						... on Document {
							id
							meta {
								title
								path
								image
								description
								publishDate
								template
								estimatedReadingMinutes
								estimatedConsumptionMinutes
								audioCoverCrop {
									x
									y
									width
									height
								}  
								audioSource {
									mp3
									kind
									durationMs
								}
								format {
									meta {
										path
										title
										kind
									}
								}
							}
						}
					}
				}
			}
		}`,
		filter.String(),
	))

	req.Var("limit", limit)
	req.Header.Set("Cookie", fmt.Sprintf("connect.sid=%s", c.sid))
	var resp Response

	if err := qc.Run(context.Background(), req, &resp); err != nil {
		return nil, err
	}

	d := []Document{}

	for _, n := range resp.Documents.Nodes {
		d = append(d, n.Entity)
	}

	return d, nil
}

type GraphqlResponse struct {
	Data interface{} `json:"data"`
}

type Identifier string

const (
	Figure Identifier = "FIGURE"
	Title  Identifier = "TITLE"
	Center Identifier = "CENTER"
)

type Content struct {
	Type     string      `json:"type"`
	Children []MdAstNode `json:"children"`
	// another Meta
}
type Article struct {
	Type    string   `json:"type"`
	Id      string   `json:"id"`
	RepoId  string   `json:"repoId"`
	Content Content  `json:"content"`
	Meta    struct{} `json:"meta"`
}
type ArticleResponse struct {
	Article Article `json:"article"`
}

// Read the getDocument graphql query from a seperate file since it is verbose
// Copied it from the browser networking devtools
var queryBytes, _ = os.ReadFile("./lib/client/getDocument.gql")
var GetDocumentQuery = string(queryBytes)

// Fetch the article at path including its content and metadata
// Returns nested struct representing the article
func (c *Client) FetchArticle(path string) (*ArticleResponse, error) {
	qc := graphql.NewClient(apiURL)
	req := graphql.NewRequest(GetDocumentQuery)
	req.Var("path", path)

	req.Header.Set("Cookie", fmt.Sprintf("connect.sid=%s", c.sid))
	var resp ArticleResponse
	if err := qc.Run(context.Background(), req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Returns article at path in HTML format
func (c *Client) GetArticleHTML(path string) string {
	resp, err := c.FetchArticle(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fetching article %s failed with error %s", path, err))
		return "<p>Fetching article failed </p>"
	}
	return resp.ToHTML()
}

// NewClient returns a new API client
func NewClient(SID string) *Client {
	return &Client{
		sid: SID,
	}
}
