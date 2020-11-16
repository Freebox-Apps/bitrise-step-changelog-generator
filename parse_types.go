package main

import (
	"os"
	"regexp"
	"strings"
)

const (
	TypeEnv       = "type_prefix"
	TypeSeparator = "|"
	TypeRegex     = "^(.*)[\\(](.*)[\\)]"
)

func extractTypeList() []string {
	return strings.Split(os.Getenv(TypeEnv), TypeSeparator)
}

func createEntries(prefixList []string) []Entry {
	var result []Entry
	regex := regexp.MustCompile(TypeRegex)
	for i := 0; i < len(prefixList); i++ {
		prefix := prefixList[i]
		typeConf := regex.FindStringSubmatch(prefix)

		if len(typeConf) > 1 {
			result = append(result, makeTypeEntry(typeConf[1], typeConf[2]))
		} else {
			result = append(result, makeTypeEntry(prefix, prefix))
		}
	}
	return result
}

func makeTypeEntry(id string, label string) Entry {
	return Entry{id, label, make(map[string][]Commit)}
}
