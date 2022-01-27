package circleci

import (
	"context"
	"io"
	"log"
	"os"
	"net/http"
	"net/url"
)

const (
	APIv1 string = "https://circleci.com/api/v1.1/"
	APIv2 string = "https://circleci.com/api/v2/"
)

type Client struct {
	Projects	*ProjectsService
	
	client		*http.Client
	common 		service
	v1api		*url.URL
	v2api		*url.URL
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient ==  nil {
		httpClient = &http.Client{}
	}

	v1, _ := url.Parse(APIv1)
	v2, _ := url.Parse(APIv2)

	c := &Client{
		client: httpClient,
		v1api: v1,
		v2api: v2,
	}
	c.common.client = c
	c.Projects = (*ProjectsService)(&c.common)

	return c
}

func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Print("Error creating request", err)
		return nil, err
	}

	req.Header.Set("circle-token", os.Getenv("CIRCLE_TOKEN"))
	req.Header.Set("content-type", "appliation/json")

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		log.Print("Error completing request", err)
		return nil, err
	}
	
	return res, nil
}