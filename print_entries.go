package main

import (
	"fmt"
	"sort"
)

func displayEntries(entries []Entry) {
	fmt.Printf("Displayable Commits by type :\n")
	for j := 0; j < len(entries); j++ {
		entry := entries[j]
		keys := make([]string, 0, len(entry.commitMap))
		for k := range entry.commitMap {
			keys = append(keys, k)
		}
		fmt.Printf("type : %s\t\ttitle: %s\t\tmessage count: %d\n", entry.id, entry.name, len(keys))
	}
	fmt.Printf("\n")
}

func getBasicResult(entries []Entry) string {
	var result string

	for typeIndex := 0; typeIndex < len(entries); typeIndex++ {
		entry := entries[typeIndex]

		if len(entry.commitMap) > 0 {
			result += "\t" + entry.name
			result += "\n\n"

			keys := getSortedKeys(entry)
			for j := 0; j < len(keys); j++ {
				key := keys[j]
				commitList := entry.commitMap[key]
				result += "\t\t" + key

				if len(commitList) > 1 {
					for msgIndex := 0; msgIndex < len(commitList); msgIndex++ {
						result += "\n"
						result += "\t\t\t- " + commitList[msgIndex]
					}
				} else {
					result += commitList[0]
				}
				result += "\n"
			}
			result += "\n\n"
		}
	}

	if len(result) == 0{
		fmt.Printf("\n\n === No Changelog Generated === \n\n")
	}

	return result
}

func getSortedKeys(entry Entry) []string {
	keys := make([]string, 0, len(entry.commitMap))
	for k := range entry.commitMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}