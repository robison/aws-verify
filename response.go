package main

import "encoding/json"
import "net/http"

/**
 * Create a new Response instance
 */
func NewResponse(code int, success bool) *Response {
	return &Response{
		V:       "1.0",
		Code:    code,
		Success: success,
	}
}

/**
 * Structure for a JSON response to verification requests
 */
type Response struct {
	V        string                 `json:"v"`
	Code     int                    `json:"code"`
	Success  bool                   `json:"success"`
	Error    string                 `json:"error"`
	Document map[string]interface{} `json:"document"`
}

/**
 * Parse the JSON identity document in the signature's content and attach it to
 * the response object. This allows the client to parse the response JSON and extract
 * the Document as a native Object/Map
 */
func (response *Response) AddDocument(content []byte) error {
	return json.Unmarshal(content, &response.Document)
}

/**
 * Add an error's message string to the Response
 */
func (response *Response) AddError(err error) {
	response.Error = err.Error()
}

/**
 * Marshal the Response to JSON and write it to the outgoing HTTP message
 */
func (response *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(response.Code)

	return json.NewEncoder(w).Encode(response)
}
