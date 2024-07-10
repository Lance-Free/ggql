package ggql

import (
	"bytes"
	"encoding/json"
	"github.com/samber/mo"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

// Request represents an HTTP request to a specific endpoint with optional headers.
type Request struct {
	Endpoint, Request string
	Headers           map[string]string
	Variables         map[string]any
}

// NewRequest initializes a new Request object with the specified endpoint and an empty header map.
func NewRequest(endpoint string) Request {
	return Request{
		Endpoint:  endpoint,
		Headers:   make(map[string]string),
		Variables: make(map[string]any),
	}
}

// AddHeader adds a header to the request. It takes a key-value pair and updates the
// Headers map in the Request struct. The updated Request is then returned.
func (request Request) AddHeader(key, value string) Request {
	request.Headers[key] = value
	return request
}

// AddHeaders appends the key-value pairs in the provided headers map to the
// Request's Headers map. It iterates over each key-value pair in the headers
// map and adds it to the Headers map of the Request struct. The modified
// Request is then returned.
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

// ClearHeaders resets the Headers map in the Request struct by creating a new empty map.
// It returns the updated Request.
func (request Request) ClearHeaders() Request {
	request.Headers = make(map[string]string)
	return request
}

// AddVariable adds a variable to the request. It takes a key-value pair and updates the
// Variables map in the Request struct. The updated Request is then returned.
func (request Request) AddVariable(key string, value any) Request {
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
	request.Variables = make(map[string]any)
	return request
}

// AddVariables appends the key-value pairs in the provided variables map to the Request's Variables map.
// It iterates through the variables map and assigns each key-value pair to the corresponding key in the Request's Variables map.
// The updated Request struct is then returned.
func (request Request) AddVariables(variables map[string]any) Request {
	for key, value := range variables {
		request.Variables[key] = value
	}
	return request
}

// Query sets the query for the request. It updates the Request field of the
// Request struct and returns the modified Request.
func (request Request) Query(query string) Request {
	request.Request = query
	return request
}

// content represents the request payload for an HTTP request sent to a GraphQL endpoint.
// It contains a query string and a map of variables.
type content struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

// Do sends an HTTP POST request to the specified endpoint with the query/mutation from the Request.
// It encodes the request payload, sets the "Content-Type" header to "application/json",
// sends the request, reads the response body, and returns the parsed response as a gjson.Result.
// If the request is empty, it returns an error indicating that no query/mutation is provided.
// If there is an error encoding the request payload, creating the request, sending the request,
// or reading the response, it returns an error with the corresponding error message.
// The response is always closed before returning.
func (request Request) Do() mo.Result[gjson.Result] {
	if request.Request == "" {
		return mo.Errf[gjson.Result]("no query/mutation provided")
	}

	c := content{
		Query:     request.Request,
		Variables: request.Variables,
	}

	var reqBuf bytes.Buffer
	err := json.NewEncoder(&reqBuf).Encode(c)
	if err != nil {
		return mo.Errf[gjson.Result]("encoding request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, request.Endpoint, &reqBuf)
	if err != nil {
		return mo.Errf[gjson.Result]("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return mo.Errf[gjson.Result]("sending request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	var resBuf bytes.Buffer
	_, err = resBuf.ReadFrom(res.Body)
	if err != nil {
		return mo.Errf[gjson.Result]("reading response: %w", err)
	}

	return mo.Ok[gjson.Result](gjson.ParseBytes(resBuf.Bytes()))
}
