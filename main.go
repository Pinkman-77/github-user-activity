package main

import (
        "encoding/json"
        "fmt"
        "net/http"
		"io/ioutil"
		"os"

)

func fetchActivity(username string) ([]byte, error) {
        url := "https://api.github.com/users/" + username + "/events"
        resp, err := http.Get(url)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("fetch failed with status %d", resp.StatusCode)
		}

      
        return ioutil.ReadAll(resp.Body)
}

func printActivity(data []byte) error {
        var events []interface{}
        err := json.Unmarshal(data, &events)
        if err != nil {
                return err
        }
 // Iterate over each event
 for _, event := range events {
	// Type assert the event to a map[string]interface{} to access fields
	m := event.(map[string]interface{})

	// Extract common fields
	eventType := m["type"].(string)
	repoName := m["repo"].(map[string]interface{})["name"].(string)
	actorLogin := m["actor"].(map[string]interface{})["login"].(string)
	
	// Handle different event types
	switch eventType {
	case "PushEvent":
			// Extract commit information
			commits := m["payload"].(map[string]interface{})["commits"].([]interface{})
			fmt.Printf("- Pushed %d commits to %s (%s)\n", len(commits), repoName, actorLogin)
	case "IssuesEvent":
			// Extract issue action
			action := m["payload"].(map[string]interface{})["action"].(string)
			fmt.Printf("- %s an issue in %s (%s)\n", action, repoName, actorLogin)
	case "WatchEvent":
			fmt.Printf("- Starred %s (%s)\n", repoName, actorLogin)
	default:
			fmt.Printf("- %s event in %s (%s)\n", eventType, repoName, actorLogin)
	}
}
return nil
	}
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <username>\n", os.Args[0])
		os.Exit(1)
	}

	username := os.Args[1]
	data, err := fetchActivity(username)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	err = printActivity(data)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)

	}
}
	



	

