package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Har struct {
	Log Log `json:"log"`
}

type Log struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Request Request `json:"request"`
}

type Request struct {
	Cookies []Cookie `json:"cookies"`
}

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func main() {
	// Check if the user has provided a file name argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <har_file_name>")
		os.Exit(1)
	}

	// The first argument is the name of the HAR file
	harFileName := os.Args[1]

	// Read the HAR file into a byte array
	fileBytes, err := ioutil.ReadFile(harFileName)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", harFileName, err)
		os.Exit(1)
	}

	// Parse the HAR file into our struct
	var harFile Har
	err = json.Unmarshal(fileBytes, &harFile)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// Iterate through all cookies and flag (and scramble) potential session cookies
	for i, entry := range harFile.Log.Entries {
		for j, cookie := range entry.Request.Cookies {
			if isSessionCookie(cookie.Name) {
				fmt.Printf("Unsafe to share, scrambling: %s=%s\n", cookie.Name, cookie.Value)
				harFile.Log.Entries[i].Request.Cookies[j].Value = scrambleCookieValue(cookie.Value)
			}
		}
	}

	// Serialize the modified HAR file back to JSON
	modifiedBytes, err := json.MarshalIndent(harFile, "", "  ")
	if err != nil {
		fmt.Println("Error serializing to JSON:", err)
		os.Exit(1)
	}

	// Write the modified HAR file back to disk
	err = ioutil.WriteFile("safe_to_share.har", modifiedBytes, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}

	fmt.Println("Modified HAR file has been saved as safe_to_share.har")
}

// A function to check if a cookie name suggests it's a session cookie
func isSessionCookie(name string) bool {
	return name == "SESSIONID" || name == "JSESSIONID" || name == "ASP.NET_SessionId"
}

// A function to scramble a cookie value
func scrambleCookieValue(value string) string {
	bytes := []byte(value)
	for i := range bytes {
		bytes[i] = byte('X')
	}
	return string(bytes)
}
