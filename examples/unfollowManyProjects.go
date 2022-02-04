package main

import (
	"context"
	"log"

	"github.com/olukotun-ts/circleci-client-go/circleci"
)

func main() {
	ctx := context.Background()
	c := circleci.NewClient(nil)

	slugs := []string{
		"gh/olukotun-ts/confluent-kafka-go",
		"gh/olukotun-ts/circleci-demo",
		"gh/olukotun-ts/circleci-demo-ruby-rails",
	}

	responses, err := c.Projects.UnfollowMany(ctx, slugs)
	if err != nil {
		log.Fatal(err)
	}

	for _, resp := range responses {
		log.Print("Request completed for: ", resp.Project.Slug)
		log.Print("Follow status: ", resp.Following)
	}
}
