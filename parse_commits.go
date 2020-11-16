package main

import (
	"regexp"
	"strings"
)

const (
	CommitSeparator = "\n\n"
	ScopeRegex      = "([\\(].*[\\)][\\:])"
	ScopeValueRegex = "[\\(](.*)[\\)][\\:]"
)

type Entry struct {
	id        string
	name      string
	commitMap map[string][]Commit // Key is scope and value is list commits related
}

type Commit struct {
	message   string
	ticketIds []string
}

func extractCommitListFromString(commits string) []string {
	return strings.Split(commits, CommitSeparator)
}

func fillCommitInfo(commits []string, entries []Entry) {
	for i := 0; i < len(commits); i++ { //looping from 0 to the length of the array
		message := commits[i]
		for j := 0; j < len(entries); j++ {
			_type := entries[j].id
			if strings.HasPrefix(message, _type) {
				firstLine := strings.Split(message, "\n")[0]
				noTypeMessage := strings.TrimLeft(firstLine, _type)
				scope := extractScope(noTypeMessage)
				noScopeMessage := cleanScope(noTypeMessage)
				ticketIds := extractSovledTickets(message)
				commits := Commit{noScopeMessage, ticketIds}
				createOrAppendCommit(entries[j].commitMap, scope, commits)
			}
		}
	}
}

func createOrAppendCommit(commitMap map[string][]Commit, scope string, commit Commit) {
	commitMap[scope] = append(commitMap[scope], commit)
}

func extractScope(message string) string {
	regex := regexp.MustCompile(ScopeValueRegex)
	scope := regex.FindStringSubmatch(message)
	if len(scope) > 1 && scope[1] != "" {
		return scope[1]
	} else {
		return "Unknown"
	}
}

func cleanScope(message string) string {
	regex := regexp.MustCompile(ScopeRegex)
	scope := regex.FindStringSubmatch(message)
	if len(scope) > 1 {
		noScopeMessage := strings.TrimLeft(message, scope[1])
		return strings.Trim(noScopeMessage, " ")
	} else {
		return message
	}
}
