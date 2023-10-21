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
	Version string  `json:"version"`
	Creator Creator `json:"creator"`
	Pages   []Page  `json:"pages,omitempty"`
	Entries []Entry `json:"entries"`
}

type Creator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Page struct {
	StartedDateTime string      `json:"startedDateTime"`
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	PageTimings     PageTimings `json:"pageTimings"`
}

type PageTimings struct {
	OnContentLoad float64 `json:"onContentLoad"`
	OnLoad        float64 `json:"onLoad"`
}

type Entry struct {
	Pageref         string                 `json:"pageref,omitempty"`
	StartedDateTime string                 `json:"startedDateTime"`
	Time            float64                `json:"time"`
	Request         Request                `json:"request"`
	Response        Response               `json:"response"`
	Cache           Cache                  `json:"cache"`
	Timings         Timings                `json:"timings"`
	ServerIPAddress string                 `json:"serverIPAddress,omitempty"`
	Initiator       map[string]interface{} `json:"_initiator"`
	Priority        string                 `json:"_priority"`
	ResourceType    string                 `json:"_resourceType"`
	Connection      string                 `json:"connection"`
	UnknownFields   map[string]interface{} `json:"-"`
}

func (e *Entry) UnmarshalJSON(data []byte) error {
	var objMap map[string]interface{}
	if err := json.Unmarshal(data, &objMap); err != nil {
		return err
	}

	// Manually assign known fields
	if val, ok := objMap["pageref"].(string); ok {
		e.Pageref = val
	}
	if val, ok := objMap["startedDateTime"].(string); ok {
		e.StartedDateTime = val
	}
	if val, ok := objMap["time"].(float64); ok {
		e.Time = val
	}

	var marshalAndUnmarshal = func(key string, dest interface{}) error {
		if val, ok := objMap[key]; ok {
			bytes, err := json.Marshal(val)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(bytes, dest); err != nil {
				return err
			}
		}
		return nil
	}

	if err := marshalAndUnmarshal("request", &e.Request); err != nil {
		return err
	}
	if err := marshalAndUnmarshal("response", &e.Response); err != nil {
		return err
	}
	if err := marshalAndUnmarshal("cache", &e.Cache); err != nil {
		return err
	}
	if err := marshalAndUnmarshal("timings", &e.Timings); err != nil {
		return err
	}
	if err := marshalAndUnmarshal("serverIPAddress", &e.ServerIPAddress); err != nil {
		return err
	}
	if err := marshalAndUnmarshal("_initiator", &e.Initiator); err != nil {
		return err
	}
	if err := marshalAndUnmarshal("_priority", &e.Priority); err != nil {
		return err
	}
	if err := marshalAndUnmarshal("_resourceType", &e.ResourceType); err != nil {
		return err
	}
	if err := marshalAndUnmarshal("connection", &e.Connection); err != nil {
		return err
	}

	// Create and populate UnknownFields
	e.UnknownFields = make(map[string]interface{})
	for key, val := range objMap {
		switch key {
		case "pageref", "startedDateTime", "time", "request", "response", "cache", "timings", "serverIPAddress", "_initiator", "_priority", "_resourceType", "connection":
			// Ignore known fields
		default:
			e.UnknownFields[key] = val
		}
	}

	return nil
}

func (e *Entry) MarshalJSON() ([]byte, error) {
	type Alias Entry
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	// First, marshal known fields
	knownFields, err := json.Marshal(aux.Alias)
	if err != nil {
		return nil, err
	}

	// Then, marshal unknown fields
	unknownFields, err := json.Marshal(e.UnknownFields)
	if err != nil {
		return nil, err
	}

	// Unmarshal known and unknown fields into maps
	var knownMap map[string]interface{}
	var unknownMap map[string]interface{}
	if err := json.Unmarshal(knownFields, &knownMap); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(unknownFields, &unknownMap); err != nil {
		return nil, err
	}

	// Merge the maps
	for k, v := range unknownMap {
		knownMap[k] = v
	}

	// Finally, marshal the merged map back to JSON
	merged, err := json.Marshal(knownMap)
	if err != nil {
		return nil, err
	}

	return merged, nil
}

type Request struct {
	Method      string   `json:"method"`
	URL         string   `json:"url"`
	HTTPVersion string   `json:"httpVersion"`
	Cookies     []Cookie `json:"cookies"`
	Headers     []Header `json:"headers"`
	QueryString []NVP    `json:"queryString"`
	PostData    PostData `json:"postData"`
	HeadersSize int      `json:"headersSize"`
	BodySize    int      `json:"bodySize"`
}

type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Path     string `json:"path"`
	Domain   string `json:"domain"`
	SameSite string `json:"sameSite"`
	Expires  string `json:"expires,omitempty"`
	HTTPOnly bool   `json:"httpOnly"`
	Secure   bool   `json:"secure"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type NVP struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PostData struct {
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}

type Response struct {
	Status       int      `json:"status"`
	StatusText   string   `json:"statusText"`
	HTTPVersion  string   `json:"httpVersion"`
	Cookies      []Cookie `json:"cookies"`
	Headers      []Header `json:"headers"`
	Content      Content  `json:"content"`
	RedirectURL  string   `json:"redirectURL"`
	HeadersSize  int      `json:"headersSize"`
	BodySize     int      `json:"bodySize"`
	TransferSize int      `json:"_transferSize"`
	Error        string   `json:"_error,omitempty"`
}

type Content struct {
	Size     int    `json:"size"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
	Encoding string `json:"encoding,omitempty"`
}

type Cache struct {
	BeforeRequest CacheInfo `json:"beforeRequest"`
	AfterRequest  CacheInfo `json:"afterRequest"`
}

type CacheInfo struct {
	Expires    string `json:"expires,omitempty"`
	LastAccess string `json:"lastAccess"`
	ETag       string `json:"etag,omitempty"`
	HitCount   int    `json:"hitCount"`
}

type Timings struct {
	Blocked         float64 `json:"blocked"`
	BlockedQueueing float64 `json:"_blocked_queueing"`
	DNS             float64 `json:"dns"`
	Connect         float64 `json:"connect"`
	Send            float64 `json:"send"`
	Wait            float64 `json:"wait"`
	Receive         float64 `json:"receive"`
	SSL             float64 `json:"ssl,omitempty"`
}

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

func compareMaps(map1, map2 map[string]interface{}, prefix string) {
	for k1, v1 := range map1 {
		v2, ok := map2[k1]
		newPrefix := prefix + k1 + "."

		if !ok {
			fmt.Printf("Key missing in second map: %s%s\n", prefix, k1)
			continue
		}

		switch v1Typed := v1.(type) {
		case map[string]interface{}:
			v2Typed, ok := v2.(map[string]interface{})
			if !ok {
				fmt.Printf("Type mismatch: %s%s\n", prefix, k1)
				continue
			}
			compareMaps(v1Typed, v2Typed, newPrefix)

		case []interface{}:
			v2Typed, ok := v2.([]interface{})
			if !ok {
				fmt.Printf("Type mismatch: %s%s\n", prefix, k1)
				continue
			}
			compareSlices(v1Typed, v2Typed, newPrefix)

		default:
			if v1 != v2 {
				fmt.Printf("Value mismatch at %s%s: %v != %v\n", prefix, k1, v1, v2)
			}
		}
	}

	for k2 := range map2 {
		if _, ok := map1[k2]; !ok {
			//fmt.Printf("Key missing in first map: %s%s\n", prefix, k2)
		}
	}
}

func compareSlices(slice1, slice2 []interface{}, prefix string) {
	if len(slice1) != len(slice2) {
		fmt.Printf("Slice length mismatch at %s: %d != %d\n", prefix, len(slice1), len(slice2))
		return
	}

	for i, v1 := range slice1 {
		v2 := slice2[i]
		newPrefix := fmt.Sprintf("%s[%d].", prefix, i)

		switch v1Typed := v1.(type) {
		case map[string]interface{}:
			v2Typed, ok := v2.(map[string]interface{})
			if !ok {
				fmt.Printf("Type mismatch in slice at %s\n", newPrefix)
				continue
			}
			compareMaps(v1Typed, v2Typed, newPrefix)

		case []interface{}:
			v2Typed, ok := v2.([]interface{})
			if !ok {
				fmt.Printf("Type mismatch in slice at %s\n", newPrefix)
				continue
			}
			compareSlices(v1Typed, v2Typed, newPrefix)

		default:
			if v1 != v2 {
				fmt.Printf("Value mismatch in slice at %s: %v != %v\n", newPrefix, v1, v2)
			}
		}
	}
}

func sanitizeHeaders(headers []Header) []Header {
	sanitizedHeaders := []Header{}
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

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <command> <har_file_name>")
		fmt.Println("Commands: sanitize, reorder")
		os.Exit(1)
	}

	command := os.Args[1]
	harFileName := os.Args[2]

	fileBytes, err := ioutil.ReadFile(harFileName)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", harFileName, err)
		os.Exit(1)
	}

	var harFile Har
	err = json.Unmarshal(fileBytes, &harFile)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err)
		os.Exit(1)
	}

	switch command {
	case "sanitize":
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

		modifiedBytes, err := json.MarshalIndent(harFile, "", "  ")
		if err != nil {
			fmt.Printf("Error serializing to JSON: %s\n", err)
			os.Exit(1)
		}

		modifiedFileName := "modified_" + harFileName
		err = ioutil.WriteFile(modifiedFileName, modifiedBytes, 0644)
		if err != nil {
			fmt.Printf("Error writing file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Modified HAR file has been saved as %s\n", modifiedFileName)

	case "reorder":
		orderedBytes, err := json.MarshalIndent(harFile, "", "  ")
		if err != nil {
			fmt.Printf("Error serializing to JSON: %s\n", err)
			os.Exit(1)
		}

		orderedFileName := "example_ordered.har"
		err = ioutil.WriteFile(orderedFileName, orderedBytes, 0644)
		if err != nil {
			fmt.Printf("Error writing file: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Ordered HAR file has been saved as %s\n", orderedFileName)

	case "compare":
		if len(os.Args) < 4 {
			fmt.Println("Usage for compare: go run main.go compare <json_file1_name> <json_file2_name>")
			os.Exit(1)
		}

		jsonFileName1 := os.Args[2]
		jsonFileName2 := os.Args[3]

		file1Bytes, err := ioutil.ReadFile(jsonFileName1)
		if err != nil {
			fmt.Printf("Error reading file %s: %s\n", jsonFileName1, err)
			os.Exit(1)
		}

		file2Bytes, err := ioutil.ReadFile(jsonFileName2)
		if err != nil {
			fmt.Printf("Error reading file %s: %s\n", jsonFileName2, err)
			os.Exit(1)
		}

		var obj1, obj2 map[string]interface{}
		json.Unmarshal(file1Bytes, &obj1)
		json.Unmarshal(file2Bytes, &obj2)

		compareMaps(obj1, obj2, "")

	default:
		fmt.Println("Invalid command. Use 'sanitize', 'reorder' or 'compare'")
		os.Exit(1)
	}
}
