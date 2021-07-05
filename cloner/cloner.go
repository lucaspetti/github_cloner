package cloner

import (
	"os/exec"
)

type Commander interface {
	CloneRepo(url, path string) ([]byte, error)
}

var commander Commander

type CloneCommander struct{}

func (c CloneCommander) CloneRepo(url, path string) ([]byte, error) {
	args := []string{"-C", path, "clone", url}
	cmd := exec.Command("git", args...)

	return cmd.CombinedOutput()
}
