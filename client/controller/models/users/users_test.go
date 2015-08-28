package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/deis/deis/client/controller/api"
	"github.com/deis/deis/client/controller/client"
	"github.com/deis/deis/version"
)

const usersFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "id": 1,
            "last_login": "2014-10-19T22:01:00.601Z",
            "is_superuser": true,
            "username": "test",
            "first_name": "test",
            "last_name": "testerson",
            "email": "test@example.com",
            "is_staff": true,
            "is_active": true,
            "date_joined": "2014-10-19T22:01:00.601Z",
            "groups": [],
            "user_permissions": []
        }
    ]
}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DEIS_API_VERSION", version.APIVersion)

	if req.URL.Path == "/v1/users/" && req.Method == "GET" {
		res.Write([]byte(usersFixture))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestUsersList(t *testing.T) {
	t.Parallel()

	expected := []api.User{
		api.User{
			ID:          1,
			LastLogin:   "2014-10-19T22:01:00.601Z",
			IsSuperuser: true,
			Username:    "test",
			FirstName:   "test",
			LastName:    "testerson",
			Email:       "test@example.com",
			IsStaff:     true,
			IsActive:    true,
			DateJoined:  "2014-10-19T22:01:00.601Z",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	u, err := url.Parse(server.URL)

	if err != nil {
		t.Fatal(err)
	}

	httpClient := client.CreateHTTPClient(false)

	client := client.Client{HTTPClient: httpClient, ControllerURL: *u, Token: "abc"}

	actual, _, err := List(&client, 100)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}
