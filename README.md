# http-batch
[![Go Report Card](https://goreportcard.com/badge/github.com/logjamdev/http-batch)](https://goreportcard.com/report/github.com/logjamdev/http-batch)

A simple http batch extension for the logjam load testing tool

The aim of this module is to provide the ability to send a batch of HTTP requests concurrently and in parallel, providing responses in a readable format with information that we care about. In addition to this, the pacakage tries to make constructing the array of requests as easy and readable as possible.

## Features

* Pure Go library
* Simple and Readable API
* Less than 200 lines of code

## How it works

Behind the scenes it uses an array of easily readable Request data structs that eventually maps out to http.NewRequests(). Then it sets a concurrency limit based on the size of the batch of requests and sends them concurrently and in parallel. Once the requests are finished it orders and formats the received responses. This may make a bit more sense after reading the code and following the below example:

## Simplest Working Example

```go
package main

import (
    "fmt"

    "github.com/logjamdev/http-batch"
)

func main() {
    var requests = []Request{
    	{
	    Method:  "GET",
	    URL:     "https://httpbin.org/get",
	    Headers: nil,
            Body:    nil,
	},
	{
	    Method: "POST",
	    URL:    "https://httpbin.org/post",
	    Headers: map[string][]string{
	        "Accept": {"application/json"},
	    },
	    Body: map[string]string{
		"name":  "Test API Guy",
		"email": "testapiguy@email.com",
	    },
	},
    }

    batch, err := Batch(requests)
    if err != nil {
        fmt.Println(err)
    }
    responses := batch.Send()
    fmt.Println(responses)
}
```
