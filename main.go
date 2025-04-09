package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	RepoDirectoryEnv        = "repo_dir"
	DebugEnv                = "debug_basic"
	DebugSlackEnv           = "debug_slack"
	DebugHtmlEnv            = "debug_html"
	DebugKeyOk              = "yes"
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
		unicodeResult = unicodeResult[:BasicChangelogMaxLength]
	}

	slackResult := getMarkdownResult(entries)
	htmlResult := getHtmlResult(entries)

	if isDebugBasic() || isDebugSlack() || isDebugHtml() {
		fmt.Printf("\n    -------- Debug output(s) --------\n\n")
	}

	if isDebugBasic() {
		fmt.Printf("\t---------------- Unicode Result ----------------\n\n%s\n\n", unicodeResult)
	}

	if isDebugSlack() {
		fmt.Printf("\t---------------- Slack Result ----------------\n\n%s\n\n", slackResult)
	}

	if isDebugHtml() {
		fmt.Printf("\t---------------- HTML Result ----------------\n\n%s\n\n", htmlResult)
	}
	if isDebugBasic() || isDebugSlack() || isDebugHtml() {
		fmt.Printf("    -------------------------------\n\n")
	}

	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_BASIC", "--value", unicodeResult).CombinedOutput()
	if getWrikeAccessToken() != "" {
		exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_SLACK", "--value", slackResult).CombinedOutput()
		exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_HTML", "--value", htmlResult).CombinedOutput()
	} else {
		exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_SLACK", "--value", unicodeResult).CombinedOutput()
		exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_HTML", "--value", unicodeResult).CombinedOutput()
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

func isDebugHtml() bool {
	return os.Getenv(DebugHtmlEnv) == DebugKeyOk
}
