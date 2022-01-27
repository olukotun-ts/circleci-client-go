package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"go-client-circleci/circleci"
)

func main() {
	ctx := context.Background()
	
	circleci.LoadEnvironmentVariables()

	c := circleci.NewClient()
	proj, _ := c.Projects.Get(ctx, "confluence-kafka-go")
	if err != nil {
		log.Print(err)
	}

	log.Print("Unmarshalled proj.slug:", proj.Slug)
	log.Print("Unmarshalled proj.org:", proj.Organization)
	log.Print("Unmarshalled proj.proj:", proj.Name)
	log.Print("Unmarshalled proj.vcs.url:", proj.VCS.URL)
}
