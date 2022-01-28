package main

import (
	"context"
	"log"

	"github.com/olukotun-ts/circleci-client-go/circleci"
)

func main() {
	ctx := context.Background()
	c := circleci.NewClient(nil)
	res, err := c.Projects.Unfollow(ctx, "gh/olukotun-ts/name-button")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Completed request with: ", res.Status)
}
