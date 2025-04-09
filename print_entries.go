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
	var ticketURLPrefix = getTicketURLPrefix()

	for typeIndex := 0; typeIndex < len(entries); typeIndex++ {
		entry := entries[typeIndex]

		if len(entry.commitMap) == 0 {
			continue
		}

		result += entry.name + "\n"

		keys := getSortedKeys(entry)
		for j := 0; j < len(keys); j++ {
			key := keys[j]
			commitList := entry.commitMap[key]
			result += "\n\t• " + key + ": "

			if len(commitList) > 1 {
				for msgIndex := 0; msgIndex < len(commitList); msgIndex++ {
					result += "\n"
					result += "\t\t - " + commitToString(commitList[msgIndex], ticketURLPrefix)
				}
			} else {
				result += commitToString(commitList[0], ticketURLPrefix)
			}
			result += "\n"
		}
		result += "\n\n"
	}

	if len(result) == 0 {
		fmt.Printf("\n\n === No Debug Changelog Generated === \n\n")
	}

	return result
}

func getMarkdownResult(entries []Entry) string {
	var result string
	var ticketURLPrefix = getTicketURLPrefix()

	for typeIndex := 0; typeIndex < len(entries); typeIndex++ {
		entry := entries[typeIndex]
		showTitlePart := false
		typeResult := ""

		if len(entry.commitMap) == 0 {
			continue
		}

		keys := getSortedKeys(entry)
		for j := 0; j < len(keys); j++ {
			key := keys[j]
			commitList := entry.commitMap[key]
			showSubTitlePart := false
			scopeResult := ""

			for msgIndex := 0; msgIndex < len(commitList); msgIndex++ {
				if len(commitList[msgIndex].ticketIds) == 0 {
					continue
				}

				if !showTitlePart {
					showTitlePart = true
					scopeResult += entry.name + "\n"
				}

				if !showSubTitlePart {
					showSubTitlePart = true
					scopeResult += "\n\t• " + key
				}

				scopeResult += commitToMarkdownString(commitList[msgIndex], ticketURLPrefix)
			}

			if scopeResult != "" {
				typeResult += scopeResult + "\n"
			}
		}

		if typeResult != "" {
			result += typeResult + "\n\n"
		}
	}

	if len(result) == 0 {
		fmt.Printf("\n\n === No Slack Changelog Generated === \n\n")
	}

	return result
}

func getHtmlResult(entries []Entry) string {

	var result string
	var ticketURLPrefix = getTicketURLPrefix()

	for typeIndex := 0; typeIndex < len(entries); typeIndex++ {
		entry := entries[typeIndex]
		showTitlePart := false
		typeResult := ""

		if len(entry.commitMap) == 0 {
			continue
		}

		keys := getSortedKeys(entry)
		for j := 0; j < len(keys); j++ {
			key := keys[j]
			commitList := entry.commitMap[key]
			showSubTitlePart := false
			scopeResult := ""

			for msgIndex := 0; msgIndex < len(commitList); msgIndex++ {
				if len(commitList[msgIndex].ticketIds) == 0 {
					continue
				}

				if !showTitlePart {
					showTitlePart = true
					scopeResult += "\n<h1>" + entry.name + "</h1>"
					scopeResult += "\n<ul>"
				}

				if !showSubTitlePart {
					showSubTitlePart = true
					scopeResult += "\n\t<li><i><b>" + key + "</b></i></li>"
				}

				scopeResult += commitToHtmlString(commitList[msgIndex], ticketURLPrefix)
			}

			if scopeResult != "" {
				typeResult += scopeResult
			}
		}

		if typeResult != "" {
			result += typeResult + "\n</ul>"
		}
	}

	if len(result) == 0 {
		fmt.Printf("\n\n === No HTML Changelog Generated === \n\n")
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

func commitToString(commit Commit, urlPrefix string) string {
	var result = commit.message
	var ids = commit.ticketIds

	for i := 0; i < len(ids); i++ {
		result += " #" + ids[i]
	}
	return result
}

func commitToMarkdownString(commit Commit, urlPrefix string) string {
	var result = ""
	var ids = commit.ticketIds

	for i := 0; i < len(ids); i++ {
		var id = ids[i]
		var ticketTitle = getTitleForTicket(id)
		if ticketTitle != "" {
			result += "\n\t\t - "
			result += ticketTitle + " <" + urlPrefix + id + "|#" + id + ">"
		} else {
			result += "\n\t\t - "
			if i == 0 {
				result += commit.message
			}
			result += " <" + urlPrefix + id + "|#" + id + ">"
		}
	}
	return result
}

func commitToHtmlString(commit Commit, urlPrefix string) string {
	var result = "\n\t\t<ul>"
	var ids = commit.ticketIds

	for i := 0; i < len(ids); i++ {
		var id = ids[i]
		var ticketTitle = getTitleForTicket(id)
		result += "\n\t\t\t<li>"
		if ticketTitle != "" {
			result += ticketTitle + " <a href=" + urlPrefix + id + ">#" + id + "</a>"
		} else {
			if i == 0 {
				result += commit.message
			}
			result += " <a href=" + urlPrefix + id + ">#" + id + "</a>"
		}
		result += "</li>"
	}
	result += "\n\t\t</ul>"
	return result
}
