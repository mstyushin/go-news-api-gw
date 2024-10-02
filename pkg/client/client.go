package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpClient struct {
	cl *http.Client
}

func New() *HttpClient {
	return &HttpClient{
		cl: &http.Client{},
	}
}

// send HTTP GET to the specified URL and return response body as []byte
// will propagate request_id via X-Request-Id header (if there was one)
func (c *HttpClient) GET(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Println("ERROR creating GET request")
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	if ctx.Value("request_id") != nil {
		req.Header.Add("X-Request-Id", ctx.Value("request_id").(string))
	} else {
		log.Println("missing request_id, won't propagate")
	}

	log.Println("HTTP client: GET", url)
	res, err := c.cl.Do(req)
	if err != nil {
		log.Println("ERROR issuing HTTP GET:", url)
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR reading response body")
		return nil, err
	}
	res.Body.Close()

	return body, nil
}

// send HTTP POST to the specified URL with payload as request body and return response body as []byte
// will propagate request_id via X-Request-Id header (if there was one)
func (c *HttpClient) POST(ctx context.Context, url string, payload interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("ERROR marshalling payload")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("ERROR creating POST request")
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if ctx.Value("request_id") != nil {
		req.Header.Add("X-Request-Id", ctx.Value("request_id").(string))
	} else {
		log.Println("missing request_id, won't propagate")
	}

	log.Println("HTTP client: POST", url)
	res, err := c.cl.Do(req)
	if err != nil {
		log.Println("ERROR issuing HTTP POST:", url)
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, errors.New("Bad response")
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ERROR reading response body")
		return nil, err
	}
	res.Body.Close()

	return responseBody, nil
}
