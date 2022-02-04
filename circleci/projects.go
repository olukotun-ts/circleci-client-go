package circleci

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ProjectsService service

type Project struct {
	Slug         string  `json:"slug"`
	Organization string  `json:"organization_name"`
	Name         string  `json:"name"`
	VCS          VCSInfo `json:"vcs_info"`
}

type VCSInfo struct {
	URL           string `json:"vcs_url"`
	DefaultBranch string `json:"default_branch"`
	Provider      string `json:"provider"`
}

type ProjectResponse struct {
	Project   *Project
	Following bool `json:"following"`
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
// - Test return codes for already-(un)followed projects, missing projects, no permission, no config
// - Experiment with diags
func (svc *ProjectsService) Follow(ctx context.Context, projectSlug string, branch string) (*ProjectResponse, error) {
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

	switch res.StatusCode {
	case 200:
		resBody, _ := ioutil.ReadAll(res.Body)
		var response ProjectResponse
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			return nil, err
		}

		proj, err := svc.Get(ctx, projectSlug)
		if err != nil {
			return nil, err
		}
		response.Project = proj

		return &response, nil
	case 429:
		// For HTTP APIs, when a request is throttled, you will receive HTTP status code 429.
		// If your integration requires that a throttled request is completed, then you should retry these requests after a delay, using an exponential backoff.
		// In most cases, the HTTP 429 response code will be accompanied by the Retry-After HTTP header.
		// When this header is present, your integration should wait for the period of time specified by the header value before retrying a request.
		fallthrough
	default:
		return nil, fmt.Errorf("Expected 200 status code; got %v instead", res.StatusCode)
	}

}

// Why branch string instead of list of corresponding branches?
//		- encourage naming conventions
// 		- onboarding should either use master for all projects to follow
// 		- or create a same-named branch on all projects for POCs
func (svc *ProjectsService) FollowMany(ctx context.Context, projectSlugs []string, branch string) ([]*ProjectResponse, error) {
	responses := []*ProjectResponse{}
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

func (svc *ProjectsService) Unfollow(ctx context.Context, projectSlug string) (*ProjectResponse, error) {
	url := fmt.Sprintf("%sproject/%s/unfollow", svc.client.v1api, projectSlug)

	req, _ := svc.client.NewRequest("POST", url, nil)

	res, err := svc.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		resBody, _ := ioutil.ReadAll(res.Body)
		var response ProjectResponse
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			return nil, err
		}

		proj, err := svc.Get(ctx, projectSlug)
		if err != nil {
			return nil, err
		}
		response.Project = proj

		return &response, nil
	case 429:
		// For HTTP APIs, when a request is throttled, you will receive HTTP status code 429.
		// If your integration requires that a throttled request is completed, then you should retry these requests after a delay, using an exponential backoff.
		// In most cases, the HTTP 429 response code will be accompanied by the Retry-After HTTP header.
		// When this header is present, your integration should wait for the period of time specified by the header value before retrying a request.
		fallthrough
	default:
		return nil, fmt.Errorf("Expected 200 status code; got %v instead", res.StatusCode)
	}
}

func (svc *ProjectsService) UnfollowMany(ctx context.Context, projectSlugs []string) ([]*ProjectResponse, error) {
	responses := []*ProjectResponse{}
	for _, slug := range projectSlugs {
		resp, err := svc.Unfollow(ctx, slug)
		if err != nil {
			return nil, err
		}

		responses = append(responses, resp)
	}

	return responses, nil
}
