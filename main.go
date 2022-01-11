package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const APIv1 string = "https://circleci.com/api/v1.1"
const APIv2 string = "https://circleci.com/api/v2"

func main() {
	loadEnvironmentVariables()

	client := &http.Client{}

	url := fmt.Sprintf("%s%s", APIv1, "/projects")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("Error creating request", err)
	}
	req.Header.Set("circle-token", os.Getenv("CIRCLE_TOKEN"))
	res, err := client.Do(req)
	if err != nil {
		log.Print("Error completing request", err)
	}
	defer res.Body.Close()

	log.Print(res.Status)
	log.Print(res.Body)
}

func loadEnvironmentVariables() {
	godotenv.Load("/Users/olukotun-ts/.env")
}
