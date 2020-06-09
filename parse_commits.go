package main

import (
	"os"
	"regexp"
	"strings"
)

const (
	CommitEnv       = "commit_list"
	CommitSeparator = "\n\n"
	ScopeRegex      = "([\\(].*[\\)][\\:])"
	ScopeValueRegex = "[\\(](.*)[\\)][\\:]"
)

type Entry struct {
	id        string
	name        string
	commitMap map[string][]string // Key is scope and value is list commits related
}

func extractCommitList() []string {
	return strings.Split(os.Getenv(CommitEnv), CommitSeparator)
}

func fillCommitInfo(commits []string, entries []Entry) {
	for i := 0; i < len(commits); i++ { //looping from 0 to the length of the array
		message := commits[i]
		for j := 0; j < len(entries); j++ {
			_type := entries[j].id
			if strings.HasPrefix(message, _type) {
				noTypeMessage := strings.Trim(message, _type)
				scope := extractScope(noTypeMessage)
				noScopeMessage := cleanScope(noTypeMessage)
				createOrAppendCommit(entries[j].commitMap, scope, noScopeMessage)
			}
		}
	}
}

func createOrAppendCommit(commitMap map[string][]string, scope string, message string) {
	commitMap[scope] = append(commitMap[scope], message)
}

func extractScope(message string) string {
	regex := regexp.MustCompile(ScopeValueRegex)
	scope := regex.FindStringSubmatch(message)
	if len(scope) > 1 {
		return scope[1]
	} else{
		return "Unknown"
	}
}

func cleanScope(message string) string {
	regex := regexp.MustCompile(ScopeRegex)
	scope := regex.FindStringSubmatch(message)
	if len(scope) > 1 {
		noScopeMessage := strings.TrimLeft(message, scope[1])
		return strings.Trim(noScopeMessage, " ")
	} else{
		return message
	}
}