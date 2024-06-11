package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type StatusValue struct {
	State       string `json:"state"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
}

type GetStatusResponse struct {
	Values []StatusValue `json:"values"`
}

func main() {
	commitId := os.Getenv("COMMIT")
	token := os.Getenv("GIT_ACCESS_TOKEN")
	domain := os.Getenv("GIT_BASE_URL")

	fmt.Println("-----------------------------------------")
	fmt.Println("Inputs: ")
	fmt.Println("Domain: " + domain)
	fmt.Println("Commit Hash: " + commitId)
	fmt.Println("-----------------------------------------")
	fmt.Println("")

	fmt.Println("-----------------------------------------")
	fmt.Println("Getting statuses for commit: " + commitId)
	statuses := getStatuses(token, domain, commitId)
	fmt.Printf("Total status count is: %d", len(statuses))
	fmt.Println()
	for i := 0; i < len(statuses); i++ {
		fmt.Printf("Status %d name: %s", i, statuses[i].Name)
		fmt.Println()
	}
	fmt.Println("-----------------------------------------")
	fmt.Println("")

	fmt.Println("-----------------------------------------")
	fmt.Println("Updating Statuses...")
	for i := 0; i < len(statuses); i++ {
		status := statuses[i]
		if status.State == "INPROGRESS" {
			updated := invalidateStatus(token, domain, commitId, status)
			fmt.Printf("Updating status for "+status.Name+": %d", updated)
			fmt.Println()
		}
	}
	fmt.Println("-----------------------------------------")
	fmt.Println("")
}

func getStatuses(token string, domain string, commitId string) []StatusValue {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://"+domain+"/rest/build-status/1.0/commits/"+commitId, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
		return []StatusValue{}
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return []StatusValue{}
	}

	var responseObject GetStatusResponse
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Values
}

func invalidateStatus(token string, domain string, commitId string, status StatusValue) bool {
	status.State = "FAILED"
	reqBody, _ := json.Marshal(status)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://"+domain+"/rest/build-status/1.0/commits/"+commitId, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
		return false
	}

	_, readError := io.ReadAll(response.Body)
	if readError != nil {
		log.Fatal(readError)
		return false
	}

	if response.StatusCode >= 200 && response.StatusCode < 210 {
		return true
	}

	return false
}
