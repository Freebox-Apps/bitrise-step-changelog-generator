package main

import (
	"fmt"
	"sort"
)

func displayEntries(entries []Entry) {
	for j := 0; j < len(entries); j++ {
		fmt.Printf("%s : %s\t", entries[j].id, entries[j].name)
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