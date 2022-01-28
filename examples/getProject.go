package main

import (
	"context"
	"log"

	"github.com/olukotun-ts/circleci-client-go/circleci"
)

func main() {
	ctx := context.Background()
	c := circleci.NewClient(nil)
	proj, err := c.Projects.Get(ctx, "gh/olukotun-ts/name-button")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Project slug:", proj.Slug)
	log.Print("Project org:", proj.Organization)
	log.Print("Project name:", proj.Name)
	log.Print("Project vcs.url:", proj.VCS.URL)
}
