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
	// Read the HAR file into a byte array
	fileBytes, err := ioutil.ReadFile("example.har")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	// Parse the HAR file into our struct
	var harFile Har
	err = json.Unmarshal(fileBytes, &harFile)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// Iterate through all cookies and flag potential session cookies
	for _, entry := range harFile.Log.Entries {
		for _, cookie := range entry.Request.Cookies {
			if isSessionCookie(cookie.Name) {
				fmt.Printf("Unsafe to share: %s=%s\n", cookie.Name, cookie.Value)
			}
		}
	}
}

// A simplistic function to check if a cookie name suggests it's a session cookie
func isSessionCookie(name string) bool {
	// Here, we're just doing a very basic string match. You'd probably want a more robust way
	// to identify session cookies, perhaps even a list of known session cookie names
	// or regular expressions.
	return name == "SESSIONID" || name == "JSESSIONID" || name == "ASP.NET_SessionId"
}

