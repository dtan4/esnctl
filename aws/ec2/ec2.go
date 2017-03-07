package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
)

// Client represents a wrapper of EC2 API
type Client struct {
	api ec2iface.EC2API
}

// New creates and returns new Client object
func New(api ec2iface.EC2API) *Client {
	return &Client{
		api: api,
	}
}

// RetrieveInstanceIDFromPrivateDNS retrieves instance ID from private DNS name
func (c *Client) RetrieveInstanceIDFromPrivateDNS(privateDNS string) (string, error) {
	resp, err := c.api.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("private-dns-name"),
				Values: []*string{
					aws.String(privateDNS),
				},
			},
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve instance ID")
	}

	if len(resp.Reservations) == 0 || len(resp.Reservations[0].Instances) == 0 {
		return "", errors.Errorf("instance with %q not found", privateDNS)
	}

	return aws.StringValue(resp.Reservations[0].Instances[0].InstanceId), nil
}
