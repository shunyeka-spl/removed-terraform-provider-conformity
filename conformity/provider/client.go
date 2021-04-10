package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const BaseURL string = "https://ap-southeast-2-api.cloudconformity.com/v1"

type Client struct {
	Region     *string
	BaseUrl    *url.URL
	HTTPClient *http.Client
	Headers    http.Header
	Token      *string
}

type ErrorResponse struct {
	Errors []struct {
		Status int    `json:"status"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

// NewClient initializes a new Client instance to communicate with the Conformity api
func NewClient(h http.Header, token, region *string) *Client {
	client := Client{Region: region, Token: token}

	client.HTTPClient = &http.Client{
		Timeout: time.Second * 30,
	}

	if u, err := url.Parse(BaseURL); err != nil {
		panic("Could not init Provider client to Conformity")
	} else {
		client.BaseUrl = u
	}

	if h != nil {
		client.Headers = h
	}

	return &client
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	req.Header = c.Headers
	log.Printf("[DEBUG] Request is %v\n", req)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Printf("Error while making http request: %v\n", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	errorResponse := ErrorResponse{}
	err = json.Unmarshal(body, &errorResponse)
	if err == nil && len(errorResponse.Errors) > 0 {
		log.Printf("[DEBUG] Response is error %v\n", errorResponse.Errors[0].Detail)
		return nil, err
	}
	bodyString := string(body[:])
	log.Printf("[DEBUG] Response body is %v\n", bodyString)
	if err != nil {
		log.Printf("Error while reading response body: %v\n", err)
		return nil, err
	}

	return body, err
}
