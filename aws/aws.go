package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	autoscalingapi "github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/dtan4/esnctl/aws/autoscaling"
	"github.com/pkg/errors"
)

var (
	// AutoScaling represents AutoScaling service client
	AutoScaling *autoscaling.Client
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

	return nil
}
