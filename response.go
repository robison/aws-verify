package main

import "encoding/json"
import "net/http"

// NewResponse creates a new Response instance
func NewResponse(code int, success bool) *Response {
	return &Response{
		V:       "1.0",
		Code:    code,
		Success: success,
	}
}


// Response is an interface for the JSON response to verification requests
type Response struct {
	V        string                 `json:"v"`
	Code     int                    `json:"code"`
	Success  bool                   `json:"success"`
	Error    string                 `json:"error"`
	Document map[string]interface{} `json:"document"`
}

// AddDocument parses the JSON identity document in the signature's content and
// attaches it to the response object. This allows the client to parse the response
// JSON and extract the Document as a native Object/Map
func (response *Response) AddDocument(content []byte) error {
	return json.Unmarshal(content, &response.Document)
}

// AddError adds an error message string to the Response
func (response *Response) AddError(err error) {
	response.Error = err.Error()
}

// Send sets response headers, marshals the Response to JSON, and
// writes it to the outgoing HTTP message
func (response *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(response.Code)

	return json.NewEncoder(w).Encode(response)
}
