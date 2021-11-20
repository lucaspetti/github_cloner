package repositories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Repo struct {
	ID     int    `json:"id"`
	SSHURL string `json:"ssh_url"`
	Name   string
}

type RepositoryDataFetcher struct {
	BaseUrl    string
	httpClient HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RepositoryDataGetter interface {
	GetReposData(url string) []byte
}

func NewRepositoryDataFetcher(baseUrl string, httpClient HttpClient) *RepositoryDataFetcher {
	return &RepositoryDataFetcher{baseUrl, httpClient}
}

type UnexpectedBodyError struct {
	errorMessage string
}

func (e *UnexpectedBodyError) Error() string {
	return e.errorMessage
}

// GetUserRepos returns a slice of repos for a given username
func (f RepositoryDataFetcher) GetUserRepos(username string) ([]Repo, error) {
	url := fmt.Sprint(f.BaseUrl, username, "/repos")

	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		fmt.Printf("Could not build request: %v", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "github_cloner CLI")

	response, err := f.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	repos := []Repo{}

	if err := json.Unmarshal(responseBytes, &repos); err != nil {
		return nil, &UnexpectedBodyError{err.Error()}
	}

	return repos, nil
}
