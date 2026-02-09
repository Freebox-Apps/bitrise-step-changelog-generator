package main

import (
	"fmt"
	"os"
	"os/exec"
	"maps"
	"slices"
	"sort"
)

const (
	RepoDirectoryEnv        = "repo_dir"
	DebugEnv                = "debug_basic"
	DebugSlackEnv           = "debug_slack"
	DebugHtmlEnv            = "debug_html"
	DebugKeyOk              = "yes"
	BasicChangelogMaxLength = 16384
)

type LogEntry struct {
	title string
	ticketId string
}

type LogCategory struct {
	index int
	title string
	ticketId string
	entries []LogEntry
}

type TypeSection struct {
	index int
	title string
	categories []LogCategory
}

type ChangeLog struct {
	sections []TypeSection
}

func getJiraChangeLog(entries []Entry) ChangeLog {
	typeMap := make(map[string]TypeSection)
	for j := 0; j < len(entries); j++ {
		entry := entries[j]
		catMap := make(map[string]LogCategory)
		for k, commits := range slices.Collect(maps.Values(entry.commitMap)) {
			for l := 0; l < len(commits); l++ {
				commit := commits[l]
				for m := 0; m < len(commit.ticketIds); m++ {
					ticketId := commit.ticketIds[m]
					ticket := getTicketInfo(ticketId)
					category, exists := catMap[ticket.parentTicketId]
					if !exists {
						var entries []LogEntry
						index := k
						if ticket.parentTicketId == "" {
							index = 999
						}
						category = LogCategory{index, ticket.category, ticket.parentTicketId, entries}
					}
					logEntry := LogEntry{ticket.title, ticketId}
					category.entries = append(category.entries, logEntry)
					catMap[ticket.parentTicketId] = category
				}
			}
		}
		catValues := slices.Collect(maps.Values(catMap))
		sort.Slice(catValues, func(a, b int) bool {
			return catValues[a].index < catValues[b].index
		})
		typeMap[entry.name] = TypeSection{j, entry.name, catValues}
	}

	typeValues := slices.Collect(maps.Values(typeMap))
	sort.Slice(typeValues, func(a, b int) bool {
		return typeValues[a].index < typeValues[b].index
	})
	return ChangeLog{typeValues}
}

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
	htmlResult := getNewHtmlResult(getJiraChangeLog(entries))

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
	exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_SLACK", "--value", slackResult).CombinedOutput()
	exec.Command("bitrise", "envman", "add", "--key", "CHANGELOG_HTML", "--value", htmlResult).CombinedOutput()

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
