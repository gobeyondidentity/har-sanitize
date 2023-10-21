package main

import (
	"encoding/json"
	"fmt"
	"github.com/nmelo/har-sanitize/har"
	"io/ioutil"
	"os"
)

func isSessionCookie(name string) bool {
	// Define a list of session cookie names to scan for
	sessionCookies := []string{
		"SESSIONID", "JSESSIONID", "ASP.NET_SessionId",
		"okta-oauth-nonce", "oktaStateToken", "okta-oauth-state",
		"srefresh", "sid",
	}

	// Check if the cookie name exists in the list
	for _, sessionCookie := range sessionCookies {
		if name == sessionCookie {
			return true
		}
	}

	return false
}

func sanitizeHeaders(headers []har.Header) []har.Header {
	sanitizedHeaders := []har.Header{}
	// List of sensitive headers that should not be shared
	sensitiveHeaders := map[string]bool{
		"Authorization": true,
		"authorization": true,
		"Cookie":        true,
		"cookie":        true,
		"Set-Cookie":    true,
		"set-cookie":    true,
		// Add more headers to sanitize here
	}

	for _, header := range headers {
		if _, isSensitive := sensitiveHeaders[header.Name]; isSensitive {
			// Skip sensitive headers
			fmt.Printf("Unsafe to include header, removing: %s=%s\n", header.Name, header.Value)
			continue
		} else {
			// Keep non-sensitive headers
			sanitizedHeaders = append(sanitizedHeaders, header)
		}
	}

	return sanitizedHeaders
}

func sanitizeHar(harFile har.Har) {
	for i, entry := range harFile.Log.Entries {

		// Sanitize request headers
		harFile.Log.Entries[i].Request.Headers = sanitizeHeaders(harFile.Log.Entries[i].Request.Headers)

		// Sanitize response headers
		harFile.Log.Entries[i].Response.Headers = sanitizeHeaders(harFile.Log.Entries[i].Response.Headers)

		for j, cookie := range entry.Request.Cookies {
			if isSessionCookie(cookie.Name) {
				fmt.Printf("Unsafe to share in Request, sanitizing: %s=%s\n", cookie.Name, cookie.Value)
				harFile.Log.Entries[i].Request.Cookies[j].Value = "SANITIZED"
			}
		}
		for j, cookie := range entry.Response.Cookies {
			if isSessionCookie(cookie.Name) {
				fmt.Printf("Unsafe to share in Response, sanitizing: %s=%s\n", cookie.Name, cookie.Value)
				harFile.Log.Entries[i].Response.Cookies[j].Value = "SANITIZED"
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <har_file_name>")
		os.Exit(1)
	}

	harFileName := os.Args[1]

	fileBytes, err := ioutil.ReadFile(harFileName)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", harFileName, err)
		os.Exit(1)
	}

	var harFile har.Har
	err = json.Unmarshal(fileBytes, &harFile)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		os.Exit(1)
	}

	sanitizeHar(harFile)

	modifiedBytes, err := json.MarshalIndent(harFile, "", "  ")
	if err != nil {
		fmt.Printf("Error serializing to JSON: %s\n", err)
		os.Exit(1)
	}

	modifiedFileName := "sanitized_" + harFileName
	err = ioutil.WriteFile(modifiedFileName, modifiedBytes, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Modified HAR file has been saved as %s\n", modifiedFileName)
}
