package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Freebox-CI/bitrise-step-changelog-generator/git"
)

const (
	CommitSeparator = "commit-separator"
	CommitStartEnv  = "start_commit"
	CommitEndEnv    = "end_commit"
)

func getCommitStringList() []string {
	// get commits raw text from repo
	paramLogCommits := getCommitLogs(os.Getenv(RepoDirectoryEnv), os.Getenv(CommitStartEnv), os.Getenv(CommitEndEnv))

	// convert to an array of strings
	commitStrList := extractCommitListFromString(paramLogCommits)
	if len(commitStrList) == 0 {
		fmt.Printf("Failed to Parse changelog give either git directory or commit list as input")
		os.Exit(1)
	}
	return commitStrList
}

func getCommitLogs(dir string, commitStart string, commitEnd string) string {
	fmt.Printf("Git Directory: %s \n", dir)
	fmt.Printf("Git Start Commit: %s \n", commitStart)
	fmt.Printf("Git End Commit: %s \n", commitEnd)

	// Check if inconsistent
	if len(dir) > 0 && len(commitStart) == 0 {
		fmt.Printf("You must provide a commit from were to start changelog generation")
		os.Exit(1)
	} else if len(dir) == 0 {
		return ""
	}

	var gitCmd, _ = git.New(dir)
	fetchTags(gitCmd, dir)
	var output = getLogs(gitCmd, commitStart, commitEnd)

	fmt.Printf("\n")
	return output
}

func fetchTags(gitCmd git.Git, dir string) {
	errTag := gitCmd.FetchTags().Run()
	if errTag != nil {
		fmt.Printf("Failed to fetch tags for this repository")
		os.Exit(1)
	}
}

func getLogs(gitCmd git.Git, commitStart string, commitEnd string) string {
	// get logs
	logCmd := gitCmd.Log("%s%n%b"+CommitSeparator, commitStart, commitEnd, "--no-merges", "--children")
	var output, errLog = logCmd.RunAndReturnTrimmedOutput()
	if errLog != nil {
		fmt.Printf("Failed get logs for this repository")
		os.Exit(1)
	}
	if isDebug() {
		fmt.Printf("Git Executing command: %s\n", logCmd.GetCmd().Args)
	}
	return output
}

func extractCommitListFromString(commits string) []string {
	commitList := strings.Split(commits, CommitSeparator)
	var output = make([]string, len(commitList))
	for i := 0; i < len(commitList); i++ {
		output[i] = strings.TrimSpace(commitList[i])
	}
	return output
}
