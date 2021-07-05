package cloner

import (
	"testing"
	"fmt"
	"os/exec"
	"os"
	"strings"
)

type MockCloneCommander struct{}

func (c MockCloneCommander) CloneRepo(url, path string) ([]byte, error) {
	cs := []string{"-test.run=TestClone", "--", "git"}
	args := []string{"-C", path, "clone", url}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_TEST_OUTPUT=1"}
	out, err := cmd.CombinedOutput()
	return out, err
}

func TestClone(t *testing.T) {
	if os.Getenv("GO_WANT_TEST_OUTPUT") != "1" {
		return
	}

	defer os.Exit(0)

	args := os.Args

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	repoName := args[7]

	if !strings.Contains(repoName, ".git") {
		fmt.Fprintf(os.Stderr, "Not a git repository\n")
		os.Exit(128)
	}

	return
}

func TestCloneRepo(t *testing.T) {
	cases := []struct{
		Title         string
		url           string
		path          string
		hasError      bool
		expectedError string
	} {
		{
			"With no url",
			"",
			"/home",
			true,
			"exit status 128",
		},
		{
			"With invalid url",
			"www.invalid.com",
			"/home",
			true,
			"exit status 128",
		},
		{
			"With right path and url",
			"gitrepo.git",
			"/home",
			false,
			"",
		},
	}

	for _, test := range cases {
		t.Run(test.Title, func(t *testing.T) {
			commander = MockCloneCommander{}
			defer func() { commander = CloneCommander{} }()

			_, err := commander.CloneRepo(test.url, test.path)

			if !test.hasError && err != nil {
				t.Errorf("got unexpected error: %s", err.Error())
			}

			if test.hasError && err.Error() != test.expectedError {
				t.Errorf("got unexpected error: %s", err.Error())
			}
		})
	}
}