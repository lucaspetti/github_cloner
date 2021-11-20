package repositories

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetUserRepos(t *testing.T) {
	username := "mock_username"
	baseUrl := "http://localhost:0/"

	cases := []struct {
		Title                string
		jsonResponse         []byte
		httpError            error
		want                 []Repo
		expectedErrorMessage string
	}{
		{
			"successful request",
			[]byte(`[{"id": 1, "ssh_url": "ssh/url", "name": "repo_name"}]`),
			nil,
			[]Repo{{1, "ssh/url", "repo_name"}},
			"",
		},
		{
			"wrong body received",
			[]byte(`{}`),
			nil,
			nil,
			"json: cannot unmarshal object into Go value of type []repositories.Repo",
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			response := &http.Response{
				Body: ioutil.NopCloser(bytes.NewReader([]byte(test.jsonResponse))),
			}
			httpClient := MockHttpClient{
				response,
				test.httpError,
			}

			getter := NewRepositoryDataFetcher(baseUrl, httpClient)

			got, err := getter.GetUserRepos(username)
			want := test.want

			if err != nil && err.Error() != test.expectedErrorMessage {
				t.Errorf("got unexpected error: %v", err)
			}

			if want != nil && got[0] != want[0] {
				t.Errorf("got unexpected value. Want %v, got %v", want, got)
			}
		})
	}
}

type MockHttpClient struct {
	mockResponse *http.Response
	err          error
}

func (c MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.mockResponse, nil
}
