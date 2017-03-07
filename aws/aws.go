package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	autoscalingapi "github.com/aws/aws-sdk-go/service/autoscaling"
	ec2api "github.com/aws/aws-sdk-go/service/ec2"
	elbv2api "github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/dtan4/esnctl/aws/autoscaling"
	"github.com/dtan4/esnctl/aws/ec2"
	"github.com/dtan4/esnctl/aws/elbv2"
	"github.com/pkg/errors"
)

var (
	// AutoScaling represents AutoScaling service client
	AutoScaling *autoscaling.Client
	// EC2 represents EC2 service client
	EC2 *ec2.Client
	// ELBv2 represents ELBV2 service client
	ELBv2 *elbv2.Client
)

// Initialize creates AWS service client objects
func Initialize(region string) error {
	var (
		sess *session.Session
		err  error
	)

	if region == "" {
		sess, err = session.NewSession()
		if err != nil {
			return errors.Wrap(err, "failed to create new AWS session")
		}
	} else {
		sess, err = session.NewSession(&aws.Config{Region: aws.String(region)})
		if err != nil {
			return errors.Wrap(err, "failed to create new AWS session")
		}
	}

	AutoScaling = autoscaling.New(autoscalingapi.New(sess))
	EC2 = ec2.New(ec2api.New(sess))
	ELBv2 = elbv2.New(elbv2api.New(sess))

	return nil
}
