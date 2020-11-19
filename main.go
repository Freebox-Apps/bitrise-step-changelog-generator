package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	RepoDirectoryEnv = "repo_dir"
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

func isDebug() bool {
	return os.Getenv(DebugEnv) == DebugKeyOk
}
