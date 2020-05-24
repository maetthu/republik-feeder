package client

import (
	"context"
	"fmt"
	"time"

	"github.com/machinebox/graphql"
)

const API_URL = "https://api.republik.ch/graphql"

type Document struct {
	ID   string
	Meta struct {
		Title                       string
		Description                 string
		PublishDate                 string
		Path                        string
		Template                    string
		EstimatedReadingMinutes     int
		EstimatedConsumptionMinutes int
		AudioSource                 struct {
			MP3 string
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

func (d *Document) PubDate() time.Time {
	t, _ := time.Parse(time.RFC3339, d.Meta.PublishDate)
	return t
}

type Response struct {
	Documents struct {
		Nodes []struct {
			Entity Document
		}
	}
}

type Client struct {
	sid string
}

func (c *Client) Fetch(limit int) ([]Document, error) {
	qc := graphql.NewClient(API_URL)
	req := graphql.NewRequest(`
		query ($limit: Int!) {
		    documents: search(
		        filters: [
		            { key: "template", not: true, value: "section" }
		            { key: "template", not: true, value: "format" }
		            { key: "template", not: true, value: "front" }
		        ]
		        filter: { feed: true }
		        sort: { key: publishedAt, direction: DESC }
		        first: $limit
		    ) {
		        nodes {
		            entity {
		                ... on Document {
		                    id
		                    meta {
		                        title
		                        description
		                        publishDate
		                        path
		                        template
		                        estimatedReadingMinutes
		                        estimatedConsumptionMinutes
		                        audioSource {
		                            mp3
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
	)

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

func NewClient(SID string) *Client {
	return &Client{
		sid: SID,
	}
}
