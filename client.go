package circleci

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(req *http.Request, ctx context.Context) error

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example.
	Server string

	// HTTP client with any customized settings, such as certificate chains.
	Client http.Client

	// A callback for modifying requests which are generated before sending over
	// the network.
	RequestEditor RequestEditorFn
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses returns a ClientWithResponses with a default Client:
func NewClientWithResponses(server string) *ClientWithResponses {
	return &ClientWithResponses{
		ClientInterface: &Client{
			Client: http.Client{},
			Server: server,
		},
	}
}

// NewClientWithResponsesAndRequestEditorFunc takes in a RequestEditorFn callback function and returns a ClientWithResponses with a default Client:
func NewClientWithResponsesAndRequestEditorFunc(server string, reqEditorFn RequestEditorFn) *ClientWithResponses {
	return &ClientWithResponses{
		ClientInterface: &Client{
			Client:        http.Client{},
			Server:        server,
			RequestEditor: reqEditorFn,
		},
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetCurrentUser request
	GetCurrentUser(ctx context.Context) (*http.Response, error)
}

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
