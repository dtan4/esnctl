package ec2

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/dtan4/esnctl/aws/mock"
	"github.com/golang/mock/gomock"
)

func TestRetrieveInstanceIDFromPrivateDNS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	api := mock.NewMockEC2API(ctrl)
	api.EXPECT().DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("private-dns-name"),
				Values: []*string{
					aws.String("ip-10-0-1-23.ap-northeast-1.compute.internal"),
				},
			},
		},
	}).Return(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&ec2.Reservation{
				Instances: []*ec2.Instance{
					&ec2.Instance{
						InstanceId:     aws.String("i-1234abcd"),
						PrivateDnsName: aws.String("ip-10-0-1-23.ap-northeast-1.compute.internal"),
					},
				},
			},
		},
	}, nil)

	client := &Client{
		api: api,
	}

	privateDNS := "ip-10-0-1-23.ap-northeast-1.compute.internal"
	expected := "i-1234abcd"

	got, err := client.RetrieveInstanceIDFromPrivateDNS(privateDNS)
	if err != nil {
		t.Errorf("error should not be raised: %s", err)
	}

	if got != expected {
		t.Errorf("instance ID does not match. expected: %q, got: %q", expected, got)
	}
}
