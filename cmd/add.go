package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dtan4/esnctl/aws"
	"github.com/dtan4/esnctl/es"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	addMaxRetry     = 120
	addSleepSeconds = 5
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "add <cluster URL> <AutoScaling Group> <Number to increase (+N)>",
	Short:         "Add Elasticsearch node",
	RunE:          doAdd,
}

func doAdd(cmd *cobra.Command, args []string) error {
	if len(args) != 3 {
		return errors.New("cluster URL, AutoScaling Group and the number to increase must be specified")
	}
	clusterURL := args[0]
	autoScalingGroup := args[1]

	delta, err := strconv.Atoi(args[2])
	if err != nil {
		return errors.Wrapf(err, "invalid number to increase %q", args[2])
	}

	client, err := es.New(clusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create Elasitcsearch API client")
	}

	if err := client.DisableReallocation(); err != nil {
		return errors.Wrap(err, "failed to disable reallocation")
	}

	desiredCapacity, err := aws.AutoScaling.IncreaseInstances(autoScalingGroup, delta)
	if err != nil {
		return errors.Wrap(err, "failed to increase instance")
	}

	retryCount := 0

	for {
		nodes, err := client.ListNodes()
		if err != nil {
			return errors.Wrap(err, "failed to list nodes")
		}

		if len(nodes) == desiredCapacity {
			fmt.Print("\n")
			break
		}

		fmt.Print(".")

		if retryCount == addMaxRetry {
			return errors.New("timed out: added nodes do not join to Elasticsearch cluster")
		}

		retryCount++
		time.Sleep(addSleepSeconds)
	}

	if err := client.EnableReallocation(); err != nil {
		return errors.Wrap(err, "failed to enable reallocation")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(addCmd)
}
