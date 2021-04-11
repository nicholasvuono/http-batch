package httpbatch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/logjammdev/utils"
)

//Request holds information to create a new http request.
type Request struct {
	Method  string
	URL     string
	Headers map[string][]string
	Body    map[string]string
}

//batch holds a slice(array) of requests to be sent concurrently and in parallel.
type batch struct {
	Requests []*http.Request
}

//formattedResponse holds readable fields of a http.Response object.
type formattedResponse struct {
	Status        string
	Header        map[string][]string
	Body          string
	ContentLength int64
}

//response is a data structure to help us keep track of http.Responses received during the batch request.
type response struct {
	Index    int
	Response http.Response
	Err      error
}

//Batch is a constructor function for the batch struct.
func Batch(requests []Request) (*batch, error) {
	batch := batch{}
	err := batch.SetRequests(requests)
	return &batch, err
}

//GetRequests is a getter function for the batch's Requests field.
func (b *batch) GetRequests() []*http.Request {
	return b.Requests
}

//SetRequests is a setter function for batch's Requests field.
func (b *batch) SetRequests(requests []Request) error {
	reqs := []*http.Request{}
	for _, req := range requests {
		body, err := json.Marshal(req.Body)
		if err != nil {
			return err
		}
		request, err := http.NewRequest(
			req.Method,
			req.URL,
			bytes.NewBuffer(body),
		)
		if err != nil {
			return err
		}
		request.Header = req.Headers
		reqs = append(reqs, request)
	}
	b.Requests = reqs
	return nil
}

//Formats the batch struct and its fields into a string.
func (b *batch) String() string {
	return fmt.Sprintf("%#v", b)
}

//Sends a batch of requests concurrently and in parallel, gathering responses, ordering them, and formatting them.
func (b *batch) Send() []formattedResponse {

	client := http.DefaultClient

	semaphoreChan := make(chan struct{}, len(b.Requests))
	responsesChan := make(chan *response)

	defer func() {
		close(semaphoreChan)
		close(responsesChan)
	}()

	for i, req := range b.Requests {
		go func(i int, req *http.Request) {
			semaphoreChan <- struct{}{}
			res, err := client.Do(req)
			utils.Explain(err)
			response := &response{i, *res, err}
			responsesChan <- response
			<-semaphoreChan
		}(i, req)
	}

	var responses []response

	for {
		response := <-responsesChan
		responses = append(responses, *response)
		if len(responses) == len(b.Requests) {
			break
		}
	}

	sort.Slice(responses, func(i, j int) bool {
		return responses[i].Index < responses[j].Index
	})

	return format(responses)
}

//Formats responses into a more readable format, capturing the specific fields we care about.
func format(responses []response) []formattedResponse {
	formattedResponses := []formattedResponse{}
	for _, res := range responses {
		body, err := ioutil.ReadAll(res.Response.Body)
		utils.Explain(err)
		formattedResponse := formattedResponse{
			Status:        res.Response.Status,
			Header:        res.Response.Header,
			Body:          string(body),
			ContentLength: res.Response.ContentLength,
		}
		formattedResponses = append(formattedResponses, formattedResponse)
	}
	return formattedResponses
}
