package circleci

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func testNew(t *testing.T) *ClientWithResponses {
	t.Helper()
	godotenv.Load("/Users/olukotun-ts/.env")
	token := os.Getenv("CIRCLE_TOKEN")
	if token == "" {
		t.Fatal("CIRCLE_TOKEN env var not set")
	}
	return New(token)
}

func Test_GetCurrentUser(t *testing.T) {
	circle := testNew(t)
	resp, err := circle.GetCurrentUserWithResponse(context.TODO())
	if err != nil {
		t.Fatalf("Expected no error %s", err)
	}
	if resp.StatusCode() != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode())
	}
	if resp.JSON200 == nil {
		t.Fatalf("Expected non-nil JSON200, got %#v", resp.JSON200)
	}
	if resp.JSON200.Id == "" {
		t.Errorf("Expected nonzero Id, got %s", resp.JSON200.Id)
	}
	if resp.JSON200.Login == "" {
		t.Errorf("Expected nonzero Login, got %s", resp.JSON200.Login)
	}
	// if resp.JSON200.Name == "" {
	// 	t.Errorf("Expected nonzero Name, got %s", resp.JSON200.Name)
	// }
}
