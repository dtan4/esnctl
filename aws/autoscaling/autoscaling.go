package autoscaling

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/pkg/errors"
)

// Client represents a wrapper of AutoScaling API
type Client struct {
	api autoscalingiface.AutoScalingAPI
}

// New creates and returns new Client object
func New(api autoscalingiface.AutoScalingAPI) *Client {
	return &Client{
		api: api,
	}
}

// DetachInstance detaches instance from the given ASG
func (c *Client) DetachInstance(groupName, instanceID string) error {
	_, err := c.api.DetachInstances(&autoscaling.DetachInstancesInput{
		AutoScalingGroupName: aws.String(groupName),
		InstanceIds: []*string{
			aws.String(instanceID),
		},
		ShouldDecrementDesiredCapacity: aws.Bool(true),
	})
	if err != nil {
		return errors.Wrap(err, "failed to detach instance")
	}

	return nil
}

// IncreaseInstances increases the number of instance
func (c *Client) IncreaseInstances(groupName string, delta int) (int, error) {
	resp, err := c.api.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(groupName),
		},
	})
	if err != nil {
		return -1, errors.Wrap(err, "failed to get AutoScaling Groups")
	}

	if len(resp.AutoScalingGroups) == 0 {
		return -1, errors.Errorf("AutoScaling Groups %q does not exist", groupName)
	}
	asg := resp.AutoScalingGroups[0]

	currentDesiredCapacity := aws.Int64Value(asg.DesiredCapacity)
	targetDesiredCapacity := currentDesiredCapacity + int64(delta)

	_, err = c.api.SetDesiredCapacity(&autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(groupName),
		DesiredCapacity:      aws.Int64(targetDesiredCapacity),
	})
	if err != nil {
		return -1, errors.Wrap(err, "failed to increase desired capacity")
	}

	return int(targetDesiredCapacity), nil
}

// RetrieveTargetGroup retrieves target group ARN attached to the given ASG
func (c *Client) RetrieveTargetGroup(groupName string) (string, error) {
	resp, err := c.api.DescribeLoadBalancerTargetGroups(&autoscaling.DescribeLoadBalancerTargetGroupsInput{
		AutoScalingGroupName: aws.String(groupName),
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to retirve attached target group")
	}

	if len(resp.LoadBalancerTargetGroups) == 0 {
		return "", errors.Errorf("no target group is attached to %q", groupName)
	}

	return aws.StringValue(resp.LoadBalancerTargetGroups[0].LoadBalancerTargetGroupARN), nil
}
