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
