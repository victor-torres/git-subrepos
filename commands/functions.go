package commands

import (
	"fmt"
	"strings"

	"git-subrepos/git"
	"git-subrepos/repos"

	"github.com/fatih/color"
)

func Sync(config repos.Config) error {
	PrintRepositoryCounter(config)
	for repoName, repo := range config.Repos {
		fmt.Println("Working on", repoName)
		target := repos.ParseTarget(repo)

		exists := git.Exists(repo)
		if !exists {
			// Repository does not exist, let's clone it!
			fmt.Println("Repository does not exist at", repo.Path)
			fmt.Printf("Cloning from %s (%s: %s)...\n", repo.URL, target.Type, target.DisplayName)
			err := git.Clone(repo)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Checking out %s \"%s\"...\n", target.Type, target.DisplayName)
		err := git.Checkout(repo)
		if err != nil {
			return err
		}

		fmt.Println()
	}

	return Status(config)
}

func Status(config repos.Config) error {
	PrintRepositoryCounter(config)

	maxRepoNameLength := 0
	for repoName, _ := range config.Repos {
		if len(repoName) > maxRepoNameLength {
			maxRepoNameLength = len(repoName)
		}
	}

	for repoName, repo := range config.Repos {
		tabBuilder := &strings.Builder{}
		for i := 0; i < (maxRepoNameLength + 4 - len(repoName)); i++ {
			tabBuilder.WriteString(" ")
		}
		tabString := tabBuilder.String()

		status, isDirty, err := git.Status(repo)
		if err != nil {
			fmt.Printf("%s%sRepository does not exist at %s", repoName, tabString, repo.Path)
			continue
		}

		target := repos.ParseTarget(repo)
		dirtyStatus := ParseDirtyStatus(status, isDirty, target)

		fmt.Printf("%s%s%s %s %s\n", repoName, tabString, dirtyStatus.Icon, status, color.RedString(dirtyStatus.Reason))
	}
	return nil
}

func Run(config repos.Config, args []string) error {
	PrintRepositoryCounter(config)
	for repoName, repo := range config.Repos {
		fmt.Printf("Running on %s at %s\n", repoName, repo.Path)
		fmt.Println()
	}
	return nil
}
