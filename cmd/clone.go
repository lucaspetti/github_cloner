/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github_cloner/cloner"
	"github.com/spf13/cobra"
	// "github.com/google/go-github/github"
)

var username string
var path string

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clones selected repos",
	Long: "Clone the repos from a given username",
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" {
			fmt.Print("Enter github username: ")
			fmt.Scanf("%s", &username)
		}

		repos := getUserRepos(username)

		cloneRepos(repos)
		fmt.Println("Done")
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	cloneCmd.Flags().StringVar(
		&username, "user", "", "username to look for repos",
	)

	cloneCmd.Flags().StringVar(
		&path, "path", "/home", "path to clone repos into",
	)
}

type Repo struct {
	ID     int    `json:"id"`
	SSHURL string `json:"ssh_url"`
	Name   string
}

func getUserRepos(username string) []Repo {
	fmt.Println("Getting repos for", username)
	url := fmt.Sprint("https://api.github.com/users/", username, "/repos")
	responseBytes := getReposData(url)

	repos := []Repo{}

	if err := json.Unmarshal(responseBytes, &repos); err != nil {
		fmt.Printf("Could not unmarshal response: %v", err)
	}

	return repos
}

func getReposData(url string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		fmt.Printf("Could not send request: %v", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "github_cloner CLI")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Could not retrieve data: %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Could not read response body: %v", err)
	}

	return responseBytes
}

func cloneRepos(repos []Repo) {
	for i := range repos {
		repo := repos[i]

		var cloneProj string

		fmt.Printf("Clone %s into %s? y/n (press q to quit) ", repo.Name, path)
		fmt.Scanf("%s", &cloneProj)

		if cloneProj == "y" {
			sshURL := repo.SSHURL
			fmt.Println("Cloning...")
			cloneRepo(sshURL)
			fmt.Printf("Cloned %s into %s \n", repo.Name, path)
		} else if cloneProj == "q" {
			return
		}
	}
}

func cloneRepo(url string) {
	repoCloner := cloner.CloneCommander{}
	out, err := repoCloner.CloneRepo(url, path)

	if err != nil {
		fmt.Printf("Could not clone project to given path: %v \n", err)
		fmt.Printf(string(out))
	}
}
