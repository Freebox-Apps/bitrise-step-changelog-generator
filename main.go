package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	RepoDirectoryEnv = "repo_dir"
	DebugEnv         = "debug_basic"
	DebugSlackEnv    = "debug_slack"
	DebugKeyOk       = "yes"
	BasicChangelogMaxLength = 16384
)

func main() {
	commitStrList := getCommitStringList()
	fmt.Printf("Found %d commit candidates\n", len(commitStrList))
	prefixStrList := extractTypeList()
	entries := createEntries(prefixStrList)
	fillCommitInfo(commitStrList, entries)
	displayEntries(entries)
	unicodeResult := getBasicResult(entries)

	if len(unicodeResult) > BasicChangelogMaxLength {
		unicodeResult = [:BasicChangelogMaxLength]
	}
	
	slackResult := getSlackResult(entries)

	if isDebugBasic() {
		fmt.Printf("%s", unicodeResult)
	}

	if isDebugSlack() {
		fmt.Printf("%s", slackResult)
	}

	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_BASIC", "--value", unicodeResult).CombinedOutput()
	if getWrikeAccessToken() != "" {
		exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_SLACK", "--value", slackResult).CombinedOutput()
	} else {
		exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_SLACK", "--value", unicodeResult).CombinedOutput()
	}
	
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	} else {
		os.Exit(0) //Step as "successful"
	}
}

func isDebugBasic() bool {
	return os.Getenv(DebugEnv) == DebugKeyOk
}

func isDebugSlack() bool {
	return os.Getenv(DebugSlackEnv) == DebugKeyOk
}
