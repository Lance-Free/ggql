# ggql

## Introduction
The ggql library is a lightweight and simplified HTTP client specifically designed to interact with GraphQL endpoints. The primary goal of this library is to provide a quick way for developers to prototype applications that communicate with GraphQL based services without having to handle the overhead of complex setup and request management.

## Features
Here are some of the key features provided by ggql:

- **Request Management**: The library provides a Request struct that represents an HTTP request to a GraphQL endpoint. It allows you to set up the endpoint, headers, and variables for each request.

- **Chained Operations**: You can chain several operations to build a request. This includes adding headers and variables, removing headers and variables, clearing headers and variables, and setting up the query.

- **Header and Variable Manipulation**: The library provides functions to add, remove, clear, and set headers and variables for a request. It allows granular control over the specifications of each request.

- **Request Execution**: The `Do` function can be used to send an HTTP POST request to the specified GraphQL endpoint. It takes care of encoding the request payload, setting the appropriate "Content-Type" header, sending the HTTP request, processing the response body, and returning the parsed response.

The ggql library is minimalistic by design and intended primarily for quick prototyping. It is not meant to be a full-fledged GraphQL client library with advanced features like caching, subscriptions, or complex query management. However, it provides a simple and straightforward way to interact with GraphQL endpoints for basic use cases.