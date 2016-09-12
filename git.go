package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/reconquest/executil-go"
)

func isGitRepository() bool {
	_, err := os.Stat(".git/config")
	return err == nil
}

func getGitRemoteOrigin() (string, error) {
	output, _, err := executil.Run(
		exec.Command("git", "remote", "get-url", "origin"),
	)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(output), "\n"), nil
}
