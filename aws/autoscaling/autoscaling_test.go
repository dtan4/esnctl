package autoscaling

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/dtan4/esnctl/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestIncreaseInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockAutoScalingAPI(ctrl)
	api.EXPECT().DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String("elasticsearch"),
		},
	}).Return(&autoscaling.DescribeAutoScalingGroupsOutput{
		AutoScalingGroups: []*autoscaling.Group{
			&autoscaling.Group{
				AutoScalingGroupName: aws.String("elasticsearch"),
				DesiredCapacity:      aws.Int64(3),
			},
		},
	}, nil)
	api.EXPECT().SetDesiredCapacity(&autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String("elasticsearch"),
		DesiredCapacity:      aws.Int64(5),
	}).Return(&autoscaling.SetDesiredCapacityOutput{}, nil)

	client := &Client{
		api: api,
	}

	groupName := "elasticsearch"
	delta := 2

	got, err := client.IncreaseInstances(groupName, delta)
	if err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	expected := 5

	if got != expected {
		t.Errorf("desired capacity does not match. expected: %d, got: %d", expected, got)
	}
}

func TestRetrieveTargetGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockAutoScalingAPI(ctrl)
	api.EXPECT().DescribeLoadBalancerTargetGroups(&autoscaling.DescribeLoadBalancerTargetGroupsInput{
		AutoScalingGroupName: aws.String("elasticsearch"),
	}).Return(&autoscaling.DescribeLoadBalancerTargetGroupsOutput{
		LoadBalancerTargetGroups: []*autoscaling.LoadBalancerTargetGroupState{
			&autoscaling.LoadBalancerTargetGroupState{
				LoadBalancerTargetGroupARN: aws.String("arn:aws:elasticloadbalancing:ap-northeast-1:012345678901:targetgroup/elasticsearch/0123abcd5678efab"),
			},
		},
	}, nil)

	client := &Client{
		api: api,
	}

	groupName := "elasticsearch"
	expected := "arn:aws:elasticloadbalancing:ap-northeast-1:012345678901:targetgroup/elasticsearch/0123abcd5678efab"

	got, err := client.RetrieveTargetGroup(groupName)
	if err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	if got != expected {
		t.Errorf("target group ARN does not match. expected: %q, got: %q", expected, got)
	}
}
