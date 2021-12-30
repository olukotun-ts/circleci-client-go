package circleci

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (c *Client) GetCurrentUser(ctx context.Context) (*http.Response, error) {
	req, err := NewGetCurrentUserRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if c.RequestEditor != nil {
		err = c.RequestEditor(req, ctx)
		if err != nil {
			return nil, err
		}
	}
	return c.Client.Do(req)
}

// NewGetCurrentUserRequest generates requests for GetCurrentUser
func NewGetCurrentUserRequest(server string) (*http.Request, error) {
	var err error

	queryUrl := fmt.Sprintf("%s/me", server)

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// GetCurrentUserWithResponse request returning *GetCurrentUserResponse
func (c *ClientWithResponses) GetCurrentUserWithResponse(ctx context.Context) (*getCurrentUserResponse, error) {
	rsp, err := c.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return ParsegetCurrentUserResponse(rsp)
}

// ParsegetCurrentUserResponse parses an HTTP response from a GetCurrentUserWithResponse call
func ParsegetCurrentUserResponse(rsp *http.Response) (*getCurrentUserResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return nil, err
	}

	response := &getCurrentUserResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		response.JSON200 = &struct {
			Id    string `json:"id"`
			Login string `json:"login"`
			Name  string `json:"name"`
		}{}
		if err := json.Unmarshal(bodyBytes, response.JSON200); err != nil {
			return nil, err
		}
	}

	return response, nil
}

// Status returns HTTPResponse.Status
func (r getCurrentUserResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r getCurrentUserResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type getCurrentUserResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Id    string `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	}
}
