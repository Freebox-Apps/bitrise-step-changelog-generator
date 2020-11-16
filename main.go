package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Freebox-CI/bitrise-step-changelog-generator/git"
)

const (
	RepoDirectoryEnv = "repo_dir"
	CommitStartEnv   = "start_commit"
	CommitEndEnv     = "end_commit"
	DebugEnv         = "debug"
	DebugKeyOk       = "yes"
)

func main() {
	isDebug := isDebug()

	commitStrList := getCommitStringList()
	fmt.Printf("Found %d commit candidates\n", len(commitStrList))
	prefixStrList := extractTypeList()
	entries := createEntries(prefixStrList)
	fillCommitInfo(commitStrList, entries)
	displayEntries(entries)
	unicodeResult := getBasicResult(entries)

	if isDebug {
		fmt.Printf("%s", unicodeResult)
	}

	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_BASIC", "--value", unicodeResult).CombinedOutput()
	exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_SLACK", "--value", getSlackResult(entries)).CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	} else {
		os.Exit(0) //Step as "successful"
	}
}

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

	// fetch tags
	var gitCmd, _ = git.New(dir)
	errTag := gitCmd.FetchTags().Run()
	if errTag != nil {
		fmt.Printf("Failed to fetch tags for this repository")
		os.Exit(1)
	}

	// get logs
	logCmd := gitCmd.Log("%s%n%b", commitStart, commitEnd, "--no-merges", "--children")
	var output, errLog = logCmd.RunAndReturnTrimmedOutput()
	if errLog != nil {
		fmt.Printf("Failed get logs for this repository")
		os.Exit(1)
	}

	if isDebug() {
		fmt.Printf("Git Executing command: %s\n", logCmd.GetCmd().Args)
	}
	fmt.Printf("\n")
	return output
}

func isDebug() bool {
	return os.Getenv(DebugEnv) == DebugKeyOk
}
