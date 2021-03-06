package main

import (
	"context"
	"log"

	"github.com/olukotun-ts/circleci-client-go/circleci"
)

func main() {
	ctx := context.Background()
	c := circleci.NewClient(nil)
	res, err := c.Projects.Follow(ctx, "gh/olukotun-ts/name-button", "master")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Request completed with: ", res.Status)
}
