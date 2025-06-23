package main

import (
	"bytes"
	"encoding/base64"

	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type Response struct {
	Kind  string `json:"kind"`
	Tasks []struct {
		Id        string `json:"id"`
		AccountId string `json:"accountId"`
		Title     string `json:"title"`
	} `json:"data"`
}

const (
	WrikeTokenEnv  = "wrike_token"
	JiraDomainEnv  = "jira_domain"
	JiraTokenEnv   = "jira_token"
	JiraUserEnv    = "jira_user"
	WrikeApiUrl    = "https://www.wrike.com/api/v4/tasks"
	WrikeTicketUrl = "https://www.wrike.com/open.htm?id=%s"
	JiraApiUrl     = "https://%s/rest/api/3/issue/%s"
)

func getWrikeAccessToken() string {
	return os.Getenv(WrikeTokenEnv)
}

func getJiraAccessToken() string {
	return os.Getenv(JiraTokenEnv)
}

func getJiraUser() string {
	return os.Getenv(JiraUserEnv)
}

func getJiraDomain() string {
	return os.Getenv(JiraDomainEnv)
}

func getTitleForTicket(id string) string {
	if getWrikeAccessToken() != "" {
		return getTitleForWrikeTicket(id)
	} else if getJiraAccessToken() != "" && getJiraUser() != "" && getJiraDomain() != "" {
		return getTitleForJiraTicket(id)
	} else {
		fmt.Printf("Missing ticket manager info %s", id)
		return ""
	}
}

func getTitleForWrikeTicket(issueKey string) string {
	// Creation multipart/form-data
	bodyReq := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyReq)
	fw, _ := writer.CreateFormField("permalink")
	io.Copy(fw, strings.NewReader(fmt.Sprintf(WrikeTicketUrl, issueKey)))
	writer.Close()

	// Create req with url and body param
	req, err := http.NewRequest(http.MethodGet, WrikeApiUrl, bytes.NewReader(bodyReq.Bytes()))

	if err != nil {
		os.Exit(1)
	}

	// Add Header Params
	req.Header.Add("Authorization", "Bearer "+getWrikeAccessToken())
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, getErr := client.Do(req)
	if getErr != nil {
		os.Exit(1)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		os.Exit(1)
	}

	result := Response{}
	jsonerr := json.Unmarshal(body, &result)
	if jsonerr != nil {
		os.Exit(1)
	}

	if len(result.Tasks) > 0 {
		return result.Tasks[0].Title
	} else {
		fmt.Printf("Found no task for id : %s \n", issueKey)
		return ""
	}
}

func getTitleForJiraTicket(issueKey string) string {
	url := fmt.Sprintf("https://%s/rest/api/3/issue/%s", getJiraDomain(), issueKey)

	// Create HTTP Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error during query creation :", err)
		os.Exit(1)
	}

	// Authentification Basic (email:token)
	auth := base64.StdEncoding.EncodeToString([]byte(getJiraUser() + ":" + getJiraAccessToken()))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Accept", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request :", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error : HTTP status %d\n", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		os.Exit(1)
	}

	// Read and pars la réponse JSON
	var result struct {
		Fields struct {
			Summary string `json:"summary"`
		} `json:"fields"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error parsing JSON :", err)
		os.Exit(1)
	}
	return result.Fields.Summary
}
