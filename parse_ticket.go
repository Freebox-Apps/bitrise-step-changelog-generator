package main

import (
	"os"
	"regexp"
)

const (
	TicketUrlEnv  = "ticket_url"
	TicketIdRegex = "[#]([^,\\s\n]+)"
)

func getTicketURLPrefix() string {
	return os.Getenv(TicketUrlEnv)
}

func extractSovledTickets(message string) []string {
	regex := regexp.MustCompile(TicketIdRegex)
	matches := regex.FindAllStringSubmatch(message, -1)
	var result []string
	for _, match := range matches {
		result = append(result, match[1])
	}
	return result
}
