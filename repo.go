package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/seletskiy/hierr"
)

var (
	reParserGitRemoteURL = regexp.MustCompile(
		`://([^@]+@)?([^/]+)[\/:]([^/]+)/([[:alnum:]][\w\-\.]*?)(\.git)?$`,
	)

	reParserWebURL = regexp.MustCompile(
		`://([^/]+)/((users|projects)/([^/]+))/repos/([^/]+)`,
	)
)

type Repo struct {
	host        string
	project     string
	projectType string
	name        string
}

func GetRepo(url string) (*Repo, error) {
	var repo *Repo
	var err error

	if url == "" {
		if !isGitRepository() {
			return nil, fmt.Errorf(
				"current directory is not git working directory, " +
					"you should cd to project with git repository " +
					"or specify repository url using -u flag",
			)
		}

		repo, err = GetRepoFromGit()
		if err != nil {
			return nil, hierr.Errorf(
				err,
				"can't parse information "+
					"about remote repository using specified URL",
			)
		}
	} else {
		repo, err = GetRepoFromWebURL(url)
		if err != nil {
			return nil, hierr.Errorf(
				err, "can't parse information "+
					"about remote repository using specified url",
			)
		}
	}

	return repo, nil
}

func GetRepoFromGit() (*Repo, error) {
	remote, err := getGitRemoteOrigin()
	if err != nil {
		return nil, err
	}

	matches := reParserGitRemoteURL.FindStringSubmatch(remote)
	if len(matches) == 0 {
		return nil, hierr.Errorf(
			errors.New(remote),
			"repository does not seem to be hosted in Bitbucket (Stash)",
		)
	}

	repo := &Repo{
		host:        matches[2],
		project:     matches[3],
		name:        matches[4],
		projectType: "projects",
	}

	if strings.HasPrefix(repo.project, "~") {
		repo.projectType = "users"
		repo.project = strings.TrimPrefix(repo.project, "~")
	}

	return repo, nil
}

func GetRepoFromWebURL(url string) (*Repo, error) {
	matches := reParserWebURL.FindStringSubmatch(url)
	if len(matches) == 0 {
		return nil, fmt.Errorf(
			"repository from given uri does not seem " +
				"to be hosted in Bitbucket (Stash)",
		)
	}

	repo := &Repo{
		host:        matches[1],
		project:     matches[4],
		projectType: matches[3],
		name:        matches[5],
	}

	return repo, nil
}

func getRepoWebURL(host, project, projectType, repo string) string {
	return "http://" + host + "/" +
		projectType + "/" + project + "/repos/" + repo
}
