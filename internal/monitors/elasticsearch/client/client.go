package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const nodeStatsEndpoint = "_nodes/_local/stats/transport,http,process,jvm,indices,thread_pool"

type esClient struct {
	host       string
	port       string
	scheme     string
	username   string
	password   string
	httpClient *http.Client
}

// ESHttpClient holds methods hitting various ES stats endpoints
type ESHttpClient interface {
	GetNodeAndThreadPoolStats() (*NodeStatsOutput, string, error)
}

// NewESClient creates a new esClient
func NewESClient(host string, port string, useHTTPS bool, username string, password string) ESHttpClient {
	httpClient := &http.Client{}
	scheme := "http"

	if useHTTPS {
		scheme = "https"
	}

	return &esClient{
		host:       host,
		port:       port,
		scheme:     scheme,
		username:   username,
		password:   password,
		httpClient: httpClient,
	}
}

// Method to fetch node stats
func (c *esClient) GetNodeAndThreadPoolStats() (*NodeStatsOutput, string, error) {
	url := fmt.Sprintf("%s://%s:%s/%s", c.scheme, c.host, c.port, nodeStatsEndpoint)

	body, err := get(url, c.username, c.password, *c.httpClient)

	if err != nil {
		return nil, url, err
	}

	var nodeStatsOutput NodeStatsOutput
	err = json.Unmarshal(body, &nodeStatsOutput)

	if err != nil {
		return nil, url, err
	}

	return &nodeStatsOutput, url, nil
}

// Fetches response from the given URL
func get(url string, username string, password string, httpClient http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)
	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
