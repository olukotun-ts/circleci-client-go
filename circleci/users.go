package circleci

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type UsersService service

type User struct {
	Login 			string  `json:"login"`
	Name            string  `json:"name"`
	UUID            string  `json:"id"`
}

func (svc *UsersService) GetCurrentUser(ctx context.Context) (*User, error) {
	url := fmt.Sprintf("%sme", svc.client.v2api)
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
		var user User
		err = json.Unmarshal(body, &user)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}

	return nil, fmt.Errorf("Expected 200 status code; got %v instead", res.StatusCode)
}
