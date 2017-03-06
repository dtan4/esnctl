package v1

import (
	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v2"
)

// Client represents Elasticsearch API client
type Client struct {
	client *elastic.Client
}

// NewClient creates new Client object
func NewClient(clusterURL string) (*Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(clusterURL))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Elasticsearch API client")
	}

	return &Client{
		client: client,
	}, nil
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
