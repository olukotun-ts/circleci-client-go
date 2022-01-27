package circleci

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	// "string"
)

type ProjectsService service

type Project struct {
	Slug            string  `json:"slug"`
	Organization 	string  `json:"organization_name"`
	Name            string  `json:"name"`
	VCS          	VCSInfo `json:"vcs_info"`
}

type VCSInfo struct {
	URL        		string `json:"vcs_url"`
	DefaultBranch 	string `json:"default_branch"`
	Provider      	string `json:"provider"`
}

// todo: Include in svc.Follow() response
type FollowProjectResponse struct {
	Project
	Following	bool	`json:"following"`
	Workflow	bool	`json:"workflow"`
	FirstBuild	bool	`json:"first_build"`
}

func (svc *ProjectsService) Get(ctx context.Context, projectSlug string) (*Project, error) {
	url := fmt.Sprintf("%sproject/%s", svc.client.v2api, projectSlug)
	req, err := svc.client.NewRequest("GET", url, nil)
	if err != nil {
		log.Print("Error unmarshalling proj:", err)
		return nil, err
	}

	res, err := svc.client.Do(ctx, req)
	if err != nil {
		log.Print("Error unmarshalling proj:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode == 200 {
		var proj Project
		err = json.Unmarshal(body, &proj)
		if err != nil {
			log.Print("Error unmarshalling proj:", err)
			return nil, err
		}

		return &proj, nil
	}


	return nil, errors.New(fmt.Sprintf("Expected 200 status code; got %v instead", res.StatusCode))
}

// todo: Return response with Project attached
func (svc *ProjectsService) Follow(ctx context.Context, projectSlug string, branch string) error {
	url := fmt.Sprintf("%sproject/%s/follow", svc.client.v1api, projectSlug)

	reqBody, _ := json.Marshal(map[string]string{
		"branch": branch,
	})
	req, _ := svc.client.NewRequest("POST", url, bytes.NewBuffer(reqBody))

	res, err := svc.client.Do(ctx, req)
	if err != nil {
		log.Print("Error completing request:", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Print("Following project at ", projectSlug)
		return nil
	}

	return errors.New(fmt.Sprintf("Expected 200 status code; got %v instead", res.StatusCode))
}

func (svc *ProjectsService) Unfollow(ctx context.Context, projectSlug string) error {
	url := fmt.Sprintf("%sproject/%s/unfollow", svc.client.v1api, projectSlug)

	req, _ := svc.client.NewRequest("POST", url, nil)

	res, err := svc.client.Do(ctx, req)
	if err != nil {
		log.Print("Error completing request:", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Print("Unfollowed project at ", projectSlug)
		return nil
	}

	return errors.New(fmt.Sprintf("Expected 200 status code; got %v instead", res.StatusCode))
}
