package main

import (
	"context"
	"log"

	"github.com/olukotun-ts/circleci-client-go/circleci"
)

func main() {
	ctx := context.Background()
	c := circleci.NewClient(nil)
	resp, err := c.Projects.Unfollow(ctx, "gh/olukotun-ts/name-button")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Request completed for: ", resp.Project.Slug)
	log.Print("Follow status: ", resp.Following)
}
