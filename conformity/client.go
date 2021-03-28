package conformity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var UrlBase *string
var HttpClient *http.Client

type Client struct {
	region     string
	BaseUrl    *url.URL
	httpClient *http.Client
	headers    http.Header
}

// NewClient initializes a new Client instance to communicate with the Conformity api
func NewClient(h http.Header, authToken, region string) *Client {
	client := &Client{region: region}

	client.httpClient = &http.Client{
		Timeout: time.Second * 30,
	}

	urlstr := "http://localhost:8080"
	if u, err := url.Parse(urlstr); err != nil {
		panic("Could not init Provider client to Conformity")
	} else {
		client.BaseUrl = u
	}

	if h != nil {
		client.headers = h
	}

	return client
}

func (c *Client) doGetGroups(baseURL string) (g []map[string]interface{}, err error) {
	path := baseURL
	log.Println("[DEBUG] Something happened twice!")
	log.Printf("[DEBUG] Client is %v\n", c)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	for name, value := range c.headers {
		req.Header.Add(name, strings.Join(value, ""))
	}
	log.Printf("Request is %v\n", req)
	log.Printf("Calling %s\n", req.URL.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to issue get groups API call: %s", err.Error())
	}
	defer resp.Body.Close()
	log.Printf("Raw output: %v\n", resp)
	if resp.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type of HTTP response is invalid")
	}
	resJson := make(map[string][]map[string]interface{})
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &resJson)
	log.Printf("Groups: %v\n", resJson["data"])
	fgroups := flattenGroupsData(resJson["data"])
	return fgroups, err
}

func flattenGroupsData(groups []map[string]interface{}) []map[string]interface{} {
	if groups != nil {
		fgroups := make([]map[string]interface{}, len(groups))

		for i, group := range groups {
			fgroup := make(map[string]interface{})

			fgroup["type"] = group["type"]
			fgroup["id"] = group["id"]
			groupAttributes := group["attributes"].(map[string]interface{})
			fgroup["attributes_name"] = groupAttributes["name"]
			fgroup["attributes_tags"] = groupAttributes["tags"]
			fgroup["attributes_created_date"] = groupAttributes["created-date"]
			fgroup["attributes_last_modified_date"] = groupAttributes["last-modified-date"]
			groupRelationshipsOrgainzationData := group["relationships"].(map[string]interface{})["organisation"].(map[string]interface{})["data"].(map[string]interface{})
			fgroup["relationships_organisation_data_type"] = groupRelationshipsOrgainzationData["type"]
			fgroup["relationships_organisation_data_id"] = groupRelationshipsOrgainzationData["id"]
			groupAccountsData := group["relationships"].(map[string]interface{})["accounts"].(map[string]interface{})["data"]
			fgroup["relationships_accounts_data"] = groupAccountsData

			fgroups[i] = fgroup
		}
		return fgroups
	}

	return make([]map[string]interface{}, 0)
}
