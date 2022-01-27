/*
v make http call
- parse response
- follow project
- unfollow project (no api endpoint)
- write test
- add ci

go mod init
go mod tidy
go build -o
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const APIv1 string = "https://circleci.com/api/v1.1"
const APIv2 string = "https://circleci.com/api/v2"

type FollowProjectResponse struct {
	Following	bool	`json:"following"`
	Workflow 	bool	`json:"workflow"`
	FirstBuild  bool 	`json:"first_build"`
}

type GetProjectsResponse struct {
	Followed bool   `json:"followed"`
	Username string `json:"username"`
	Reponame string `json:"reponame"`
}

type GetProjectResponse struct {
	Slug string `json:"slug"`
	OrgName string `json:"organization_name"`
	ProjectName string `json:"name"`
	VCSInfo          struct {
		URL        string `json:"vcs_url"`
		DefaultBranch string `json:"default_branch"`
		Provider      string `json:"provider"`
	} `json:"vcs_info"`
}

// {
// 	"slug" : "gh/olukotun-ts/confluent-kafka-go",
// 	"organization_name" : "olukotun-ts",
// 	"name" : "confluent-kafka-go",
// 	"vcs_info" : {
// 	  "vcs_url" : "https://github.com/olukotun-ts/confluent-kafka-go",
// 	  "default_branch" : "master",
// 	  "provider" : "GitHub"
// 	}
//   }

// type Project struct {
// 	Slug             string  `json:"slug"`
// 	Organization string  `json:"organization_name"`
// 	Name             string  `json:"name"`
// 	VCS          VCSInfo `json:"vcs_info"`
// }
// type VCSInfo struct {
// 	Url        string `json:"vcs_url"`
// 	DefaultBranch string `json:"default_branch"`
// 	Provider      string `json:"provider"`
// }

type GetUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

func main() {
	loadEnvironmentVariables()
	// followProject()
	getProject()
}

func followProject() {
/*
- create project object
- add vcs type to project objecte
- add org name to project object

- get vcs type, org name 
- make request
- return error if not 200 response
*/

	client := &http.Client{}

	url := fmt.Sprintf("%s%s", APIv1, "/project/gh/olukotun-ts/name-button/follow")

	reqBody, err := json.Marshal(map[string]string{
		"branch": "master",
	})
	if err != nil {
		print(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Print("Error creating request", err)
	}
	req.Header.Set("circle-token", os.Getenv("CIRCLE_TOKEN"))
	req.Header.Set("content-type", "appliation/json")

	// Test setting multiple keys on header
	// req.Header = http.Header{
	// 	"Content-Type": []string{"application/json"},
	// 	"Circle-Token": os.Getenv("CIRCLE_TOKEN")
	// 	// "Circle-Token": []string{"Bearer Token"},
	// }

	res, err := client.Do(req)
	if err != nil {
		log.Print("Error completing request", err)
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Print("Error reading response body", err)
	}

	if res.StatusCode == 200 {
		var response FollowProjectResponse
		err = json.Unmarshal(bodyBytes, &response)
		if err != nil {
			log.Print("Error unmarshalling response", err)
		}

		log.Print(res.Body)
		log.Print("Response:", response)
	}

	log.Print("Status:", res.Status)
	log.Print("Response body:", string(bodyBytes))
}

func getProject() {
	url := "https://circleci.com/api/v2/project/gh/olukotun-ts/confluent-kafka-go"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("circle-token", os.Getenv("CIRCLE_TOKEN"))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	// fmt.Println(res)
	// fmt.Println(string(body))

	var response GetProjectResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Print("Error unmarshalling response:", err)
	}

	log.Print("Unmarshalled response.slug:", response.Slug)
	log.Print("Unmarshalled response.org:", response.OrgName)
	log.Print("Unmarshalled response.proj:", response.ProjectName)
	log.Print("Unmarshalled response.vcs.url:", response.VCSInfo.URL)
}

func loadEnvironmentVariables() {
	godotenv.Load(".env")
}
