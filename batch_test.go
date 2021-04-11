package httpbatch

import (
	"fmt"
	"strconv"
	"testing"
)

var tests = []func(t *testing.T){
	TestBatch,
}

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

func TestBatch(t *testing.T) {
	batch, err := Batch(requests)
	if err != nil {
		t.Error(err)
	}
	responses := batch.Send()
	if len(responses) != len(requests) {
		t.Error("Testing Error: number of responses does not equal number of requests sent!")
	}
	fmt.Println(responses)
}

func TestAll(t *testing.T) {
	for i, test := range tests {
		t.Run(strconv.Itoa(i), test)
	}
}
