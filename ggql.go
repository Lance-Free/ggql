package ggql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

// Request represents an HTTP request to a specific endpoint with optional headers.
type Request struct {
	Endpoint           string
	Headers, Variables map[string]string
}

// NewRequest initializes a new Request object with the specified endpoint and an empty header map.
func NewRequest(endpoint string) Request {
	return Request{
		Endpoint:  endpoint,
		Headers:   make(map[string]string),
		Variables: make(map[string]string),
	}
}

// AddHeader adds a header to the request. It takes a key-value pair and updates the
// Headers map in the Request struct. The updated Request is then returned.
func (request Request) AddHeader(key, value string) Request {
	request.Headers[key] = value
	return request
}

// AddHeaders appends the key-value pairs in the provided headers map to the Request's Headers map.
func (request Request) AddHeaders(headers map[string]string) Request {
	for key, value := range headers {
		request.Headers[key] = value
	}
	return request
}

// RemoveHeaders removes the specified headers from the Request's Headers map.
// The keys parameter specifies the keys of the headers to be removed.
// The function returns the modified Request.
func (request Request) RemoveHeaders(keys ...string) Request {
	for _, key := range keys {
		delete(request.Headers, key)
	}
	return request
}

func (request Request) ClearHeaders() Request {
	request.Headers = make(map[string]string)
	return request
}

// AddVariable adds a variable to the request. It takes a key-value pair and updates the
// Variables map in the Request struct. The updated Request is then returned.
func (request Request) AddVariable(key, value string) Request {
	request.Variables[key] = value
	return request
}

// RemoveVariables removes the specified variables from the Request's Variables map.
// The keys parameter specifies the keys of the variables to be removed.
// The function returns the modified Request.
func (request Request) RemoveVariables(keys ...string) Request {
	for _, key := range keys {
		delete(request.Variables, key)
	}
	return request
}

// ClearVariables clears the Variables map in the Request struct by creating
// a new empty map. It returns the updated Request.
func (request Request) ClearVariables() Request {
	request.Variables = make(map[string]string)
	return request
}

// AddVariables appends the key-value pairs in the provided variables map to the Request's Variables map.
func (request Request) AddVariables(variables map[string]string) Request {
	for key, value := range variables {
		request.Variables[key] = value
	}
	return request
}

type content struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

// Query sends a POST request to the specified endpoint with the provided query.
// It returns the result of the query as a gjson.Result and any error encountered.
// The request is created with a "Content-Type" header set to "application/json".
// Any headers previously added to the Request object are also included in the request.
// The response body is read into a buffer and parsed as a gjson.Result.
// The response body is automatically closed once it has been read.
func (request Request) Query(query string) (gjson.Result, error) {
	c := content{
		Query:     query,
		Variables: request.Variables,
	}

	var reqBuf bytes.Buffer
	err := json.NewEncoder(&reqBuf).Encode(c)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("encoding content: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, request.Endpoint, &reqBuf)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("sending request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	var resBuf bytes.Buffer
	_, err = resBuf.ReadFrom(res.Body)
	if err != nil {
		return gjson.Result{}, fmt.Errorf("reading response: %w", err)
	}

	return gjson.ParseBytes(resBuf.Bytes()), nil
}
