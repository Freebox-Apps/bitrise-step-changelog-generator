package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"strings"
	"io/ioutil"
	"mime/multipart"
	"bytes"
	"io"
	"os"
)

type Response struct {
	Kind		string  `json:"kind"`
	Tasks		[]struct {
		Id			string 	`json:"id"`
		AccountId	string 	`json:"accountId"`
		Title 		string 	`json:"title"`
	} `json:"data"`
}

const (
	WrikeTokenEnv  = "wrike_token"
)

func getWrikeAccessToken() string {
	return os.Getenv(WrikeTokenEnv)
}

func getTitleForTicket(id string) string {
	url := "https://www.wrike.com/api/v4/tasks"

	// Creation multipart/form-data
	bodyReq := &bytes.Buffer{}
    writer := multipart.NewWriter(bodyReq)
    fw, err := writer.CreateFormField("permalink")
    io.Copy(fw, strings.NewReader("https://www.wrike.com/open.htm?id=" + id))
    writer.Close()

	// Create req with url and body param
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(bodyReq.Bytes()))

	if err != nil {
		os.Exit(1)
	}	

	// Add Header Params
	req.Header.Add("Authorization", "Bearer " + getWrikeAccessToken())
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, getErr := client.Do(req)
	if getErr != nil {
		os.Exit(1)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
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
		fmt.Printf("Found no task for id : %s \n", id)
		return ""
	}
}





