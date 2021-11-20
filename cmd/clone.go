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

	"github_cloner/cloner"
	"github_cloner/repositories"

	"github.com/spf13/cobra"
)

var username string
var path string

var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clones selected repos",
	Long:  "Clone the repos from a given username",
	Run: func(cmd *cobra.Command, args []string) {
		if username == "" {
			fmt.Print("Enter github username: ")
			fmt.Scanf("%s", &username)
		}

		fmt.Println("Getting repos for", username)
		baseUrl := "https://api.github.com/users/"

		reposGetter := repositories.NewRepositoryDataFetcher(baseUrl, http.DefaultClient)
		repos, err := reposGetter.GetUserRepos(username)
		if err != nil {
			fmt.Printf("Error getting repositories for user: %v", err)
			return
		}

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

func cloneRepo(url string) {
	repoCloner := cloner.CloneCommander{}
	out, err := repoCloner.CloneRepo(url, path)

	if err != nil {
		fmt.Printf("Could not clone project to given path: %v \n", err)
		fmt.Printf(string(out))
	}
}
