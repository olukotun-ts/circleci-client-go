package cciclient

import (
	"context"
	"net/http"
)

func New(token string) *ClientWithResponses {
	return NewClientWithResponsesAndRequestEditorFunc(
		"https://circleci.com/api/v2",
		CircleCIRequestEditor(token),
	)
}

func CircleCIRequestEditor(token string) RequestEditorFn {
	return func(req *http.Request, ctx context.Context) error {
		req.Header.Add("Accept", "application/json")
		req.Header.Add("circle-token", token)
		return nil
	}
}
