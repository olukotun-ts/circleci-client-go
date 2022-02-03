package circleci

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
/*
https://github.com/hashicorp/terraform-provider-hashicups/blob/implement-create/hashicups/resource_order.go#L83-L88
		oi := hc.OrderItem{
			Coffee: hc.Coffee{
				ID: coffee["id"].(int),
			},
			Quantity: i["quantity"].(int),
		}
*/
type FollowProjectResponse struct {
	Project
	Following	bool	`json:"following"`
	Workflow	bool	`json:"workflow"`
	FirstBuild	bool	`json:"first_build"`
}

// Todo: Prompt for reponame only and reconstruct slug in function from VCS provider, org name, and reponame
func (svc *ProjectsService) Get(ctx context.Context, projectSlug string) (*Project, error) {
	url := fmt.Sprintf("%sproject/%s", svc.client.v2api, projectSlug)
	req, err := svc.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := svc.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode == 200 {
		var proj Project
		err = json.Unmarshal(body, &proj)
		if err != nil {
			return nil, err
		}

		return &proj, nil
	}

	return nil, fmt.Errorf("Expected 200 status code; got %v instead", res.StatusCode)
}

// Todo:
	// - Return response with Project attached
	// - Test return codes for already-(un)followed projects, missing projects, no permission, no config
func (svc *ProjectsService) Follow(ctx context.Context, projectSlug string, branch string) (*http.Response, error) {
	url := fmt.Sprintf("%sproject/%s/follow", svc.client.v1api, projectSlug)

	reqBody, _ := json.Marshal(map[string]string{
		"branch": branch,
	})
	req, _ := svc.client.NewRequest("POST", url, bytes.NewBuffer(reqBody))

	res, err := svc.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return res, nil
	}

	return nil, fmt.Errorf("Expected 200 status code; got %v instead", res.StatusCode)
}

// Why branch string instead of list of corresponding branches?
//		- encourage naming conventions
// 		- onboarding should either use master for all projects to follow
// 		- or create a same-named branch on all projects for POCs
func (svc *ProjectsService) FollowMany(ctx context.Context, projectSlugs []string, branch string) ([]*http.Response, error) {
	responses := []*http.Response{}
	// Todo: Explore implementation with a Go routine and concurrent requests.
	for _, slug := range projectSlugs {
		resp, err := svc.Follow(ctx, slug, branch)
		if err != nil {
			return nil, err
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

func (svc *ProjectsService) Unfollow(ctx context.Context, projectSlug string) (*http.Response, error) {
	url := fmt.Sprintf("%sproject/%s/unfollow", svc.client.v1api, projectSlug)
	
	req, _ := svc.client.NewRequest("POST", url, nil)
	
	res, err := svc.client.Do(ctx, req)
	if err != nil {
		log.Print("Error completing request:", err)
		return nil, err
	}
	defer res.Body.Close()
	
	if res.StatusCode == 200 {
		return res, nil
	}
	
	return nil, fmt.Errorf("Expected 200 status code; got %v instead", res.StatusCode)
}

func (svc *ProjectsService) UnfollowMany(ctx context.Context, projectSlugs []string) ([]*http.Response, error) {
	responses := []*http.Response{}
	for _, slug := range projectSlugs {
		resp, err := svc.Unfollow(ctx, slug)
		if err != nil {
			return nil, err
		}

		responses = append(responses, resp)
	}

	return responses, nil
}
