package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v2"
)

// Client represents Elasticsearch API client
type Client struct {
	client          *elastic.Client
	clusterEndpoint string
	httpClient      *http.Client
}

// NewClient creates new Client object
func NewClient(clusterURL string, httpClient *http.Client) (*Client, error) {
	u, err := url.Parse(clusterURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse cluster URL")
	}

	var clusterEndpoint string

	if u.User == nil {
		clusterEndpoint = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	} else {
		password, _ := u.User.Password()
		clusterEndpoint = fmt.Sprintf("%s://%s:%s@%s", u.Scheme, u.User.Username(), password, u.Host)
	}

	var client *elastic.Client

	client, err = elastic.NewClient(
		elastic.SetURL(clusterEndpoint),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Elasticsearch API client")
	}

	return &Client{
		client:          client,
		clusterEndpoint: clusterEndpoint,
		httpClient:      httpClient,
	}, nil
}

// DisableReallocation enables shard reallocation
// Modifies cluster.routing.allocation.enable to "none"
// https://www.elastic.co/guide/en/elasticsearch/reference/1.5/cluster-update-settings.html
func (c *Client) DisableReallocation() error {
	endpoint := c.clusterEndpoint + "/_cluster/settings"

	req, err := http.NewRequest("PUT", endpoint, strings.NewReader(`{"transient":{"cluster.routing.allocation.enable":"none"}}`))
	if err != nil {
		return errors.Wrap(err, "failed to make DisableReallocation request")
	}
	defer req.Body.Close()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute DisableReallocation request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		return errors.Errorf("failed to execute DisableReallocation request. code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}

// EnableReallocation enables shard reallocation
// Modifies cluster.routing.allocation.enable to "all"
// https://www.elastic.co/guide/en/elasticsearch/reference/1.5/cluster-update-settings.html
func (c *Client) EnableReallocation() error {
	endpoint := c.clusterEndpoint + "/_cluster/settings"

	req, err := http.NewRequest("PUT", endpoint, strings.NewReader(`{"transient":{"cluster.routing.allocation.enable":"all"}}`))
	if err != nil {
		return errors.Wrap(err, "failed to make EnableReallocation request")
	}
	defer req.Body.Close()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute EnableReallocation request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		return errors.Errorf("failed to execute DisableReallocation request. code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}

// ExcludeNodeFromAllocation excludes the given node from shard allocation group
// https://www.elastic.co/guide/en/elasticsearch/reference/current/allocation-filtering.html
func (c *Client) ExcludeNodeFromAllocation(nodeName string) error {
	endpoint := c.clusterEndpoint + "/_cluster/settings"
	reqBody := fmt.Sprintf(`{"transient":{"cluster.routing.allocation.exclude._name":"%s"}}`, nodeName)

	req, err := http.NewRequest("PUT", endpoint, strings.NewReader(reqBody))
	if err != nil {
		return errors.Wrap(err, "failed to make ExcludeNodeFromAllocation request")
	}
	defer req.Body.Close()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute ExcludeNodeFromAllocation request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		return errors.Errorf("failed to execute ExcludeNodeFromAllocation request. code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}

// ListNodes returns the list of node names
func (c *Client) ListNodes() ([]string, error) {
	nodesInfo, err := c.client.NodesInfo().Do()
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to execute NodesInfo API")
	}

	nodes := []string{}

	for _, node := range nodesInfo.Nodes {
		nodes = append(nodes, node.Name)
	}

	return nodes, nil
}

// ListShardsOnNode returns the list of shards on the given node
func (c *Client) ListShardsOnNode(nodeName string) ([]string, error) {
	endpoint := c.clusterEndpoint + "/_cat/shards/"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to make cat-shards request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to execute Shutdown request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return []string{}, errors.Errorf("failed to execute Shutdown request. code: %d, body: %s", resp.StatusCode, body)
	}

	lines := strings.Split(string(body), "\n")

	shardsOnNode := []string{}

	for _, line := range lines {
		if strings.Contains(line, nodeName) {
			shardsOnNode = append(shardsOnNode, line)
		}
	}

	return shardsOnNode, nil
}

// Shutdown shutdowns the given node
func (c *Client) Shutdown(nodeName string) error {
	endpoint := c.clusterEndpoint + "/_cluster/nodes/" + nodeName + "/_shutdown"

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return errors.Wrap(err, "failed to make Shutdown request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to execute Shutdown request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}

		return errors.Errorf("failed to execute Shutdown request. code: %d, body: %s", resp.StatusCode, body)
	}

	return nil
}
