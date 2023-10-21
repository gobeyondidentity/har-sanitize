package har

import "encoding/json"

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
