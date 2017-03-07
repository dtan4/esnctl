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
}

// NewClient creates new Client object
func NewClient(clusterURL string) (*Client, error) {
	u, err := url.Parse(clusterURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse cluster URL")
	}
	clusterEndpoint := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	var client *elastic.Client

	if u.User == nil {
		client, err = elastic.NewClient(
			elastic.SetURL(clusterEndpoint),
			elastic.SetSniff(false),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create Elasticsearch API client")
		}
	} else {
		password, _ := u.User.Password()
		client, err = elastic.NewClient(
			elastic.SetURL(clusterEndpoint),
			elastic.SetBasicAuth(u.User.Username(), password),
			elastic.SetSniff(false),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create Elasticsearch API client")
		}
	}

	return &Client{
		client:          client,
		clusterEndpoint: clusterEndpoint,
	}, nil
}

// DisableReallocation enables shard reallocation
// Modifies cluster.routing.allocation.enable to "none"
// https://www.elastic.co/guide/en/elasticsearch/reference/1.5/cluster-update-settings.html
func (c *Client) DisableReallocation() error {
	httpClient := &http.Client{}
	endpoint := c.clusterEndpoint + "/_cluster/settings"

	req, err := http.NewRequest("PUT", endpoint, strings.NewReader(`{"transient":{"cluster.routing.allocation.enable":"none"}}`))
	if err != nil {
		return errors.Wrap(err, "failed to make DisableReallocation request")
	}
	defer req.Body.Close()

	resp, err := httpClient.Do(req)
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
	httpClient := &http.Client{}
	endpoint := c.clusterEndpoint + "/_cluster/settings"

	req, err := http.NewRequest("PUT", endpoint, strings.NewReader(`{"transient":{"cluster.routing.allocation.enable":"all"}}`))
	if err != nil {
		return errors.Wrap(err, "failed to make EnableReallocation request")
	}
	defer req.Body.Close()

	resp, err := httpClient.Do(req)
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
	httpClient := &http.Client{}
	endpoint := c.clusterEndpoint + "/_cluster/settings"
	reqBody := fmt.Sprintf(`{"transient":{"cluster.routing.allocation.exclude._name":"%s"}}`, nodeName)

	req, err := http.NewRequest("PUT", endpoint, strings.NewReader(reqBody))
	if err != nil {
		return errors.Wrap(err, "failed to make ExcludeNodeFromAllocation request")
	}
	defer req.Body.Close()

	resp, err := httpClient.Do(req)
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

// Shutdown shutdowns the given node
func (c *Client) Shutdown(nodeName string) error {
	httpClient := &http.Client{}
	endpoint := c.clusterEndpoint + "/_cluster/nodes/" + nodeName + "/_shutdown"

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return errors.Wrap(err, "failed to make Shutdown request")
	}

	resp, err := httpClient.Do(req)
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
