package elbv2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/pkg/errors"
)

// Client represents a wrapper of ALB API
type Client struct {
	api elbv2iface.ELBV2API
}

// New creates and returns new Client object
func New(api elbv2iface.ELBV2API) *Client {
	return &Client{
		api: api,
	}
}

// DetachInstance detaches the given instance from the given target group
func (c *Client) DetachInstance(targetGroupARN, instanceID string) error {
	_, err := c.api.DeregisterTargets(&elbv2.DeregisterTargetsInput{
		TargetGroupArn: aws.String(targetGroupARN),
		Targets: []*elbv2.TargetDescription{
			&elbv2.TargetDescription{
				Id: aws.String(instanceID),
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "failed to detach instance")
	}

	return nil
}
