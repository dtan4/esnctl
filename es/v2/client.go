package v2

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"gopkg.in/olivere/elastic.v3"
)

// Client represents Elasticsearch API client
type Client struct {
	client *elastic.Client
}

// NewClient creates new Client object
func NewClient(clusterURL string) (*Client, error) {
	u, err := url.Parse(clusterURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse cluster URL")
	}
	plainURL := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	var client *elastic.Client

	if u.User == nil {
		client, err = elastic.NewClient(
			elastic.SetURL(plainURL),
			elastic.SetSniff(false),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create Elasticsearch API client")
		}
	} else {
		password, _ := u.User.Password()
		client, err = elastic.NewClient(
			elastic.SetURL(plainURL),
			elastic.SetBasicAuth(u.User.Username(), password),
			elastic.SetSniff(false),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create Elasticsearch API client")
		}
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
