package main

import (
	"context"
	"log"

	"github.com/olukotun-ts/circleci-client-go/circleci"
)

func main() {
	ctx := context.Background()
	c := circleci.NewClient(nil)
	user, err := c.Users.GetCurrentUser(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("User name:", user.Name)
	log.Print("User login:", user.Login)
	log.Print("User ID:", user.UUID)
}
