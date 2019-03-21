package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	nodeStatsEndpoint = "_nodes/_local/stats/transport,http,process,jvm,indices,thread_pool"
	clusterHealthStatsEndpoint = "_cluster/health"
)

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
	GetNodeAndThreadPoolStats() (*NodeStatsOutput, error)
	GetClusterStats() (*ClusterStatsOutput, error)
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

// Method to fetch cluster stats
func (c *esClient) GetClusterStats() (*ClusterStatsOutput, error) {
	url := fmt.Sprintf("%s://%s:%s/%s", c.scheme, c.host, c.port, clusterHealthStatsEndpoint)

	body, err := get(url, c.username, c.password, *c.httpClient)

	if err != nil {
		return nil, err
	}

	var clusterStatsOutput ClusterStatsOutput
	err = json.Unmarshal(body, &clusterStatsOutput)

	if err != nil {
		return nil, err
	}

	return &clusterStatsOutput, nil
}

// Method to fetch node stats
func (c *esClient) GetNodeAndThreadPoolStats() (*NodeStatsOutput, error) {
	url := fmt.Sprintf("%s://%s:%s/%s", c.scheme, c.host, c.port, nodeStatsEndpoint)

	body, err := get(url, c.username, c.password, *c.httpClient)

	if err != nil {
		return nil, err
	}

	var nodeStatsOutput NodeStatsOutput
	err = json.Unmarshal(body, &nodeStatsOutput)

	if err != nil {
		return nil, err
	}

	return &nodeStatsOutput, nil
}

// Fetches response from the given URL
func get(url string, username string, password string, httpClient http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get url %s: %v", url, err)
	}

	req.SetBasicAuth(username, password)
	res, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("could not get url %s: %v", url, err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("could not get url %s: %v", url, err)
	}

	return body, nil
}
