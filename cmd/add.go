package cmd

import (
	"fmt"
	"log"
	"net/http"
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
	Use:           "add",
	Short:         "Add Elasticsearch node",
	RunE:          doAdd,
}

var addOpts = struct {
	autoScalingGroup string
	clusterURL       string
	delta            int
	region           string
}{}

func doAdd(cmd *cobra.Command, args []string) error {
	if addOpts.clusterURL == "" {
		return errors.New("Elasticsearch cluster URL (--cluster-url) must be specified")
	}

	if addOpts.autoScalingGroup == "" {
		return errors.New("Auto Scaling Group (--group) must be specified")
	}

	if addOpts.delta < 1 {
		return errors.New("number to add instances must be greater than 0")
	}

	httpClient := &http.Client{}

	client, err := es.New(addOpts.clusterURL, httpClient)
	if err != nil {
		return errors.Wrap(err, "failed to create Elasitcsearch API client")
	}

	if err := aws.Initialize(addOpts.region); err != nil {
		return errors.Wrap(err, "failed to initialize AWS service clients")
	}

	log.Println("===> Disabling shard reallocation...")

	if err := client.DisableReallocation(); err != nil {
		return errors.Wrap(err, "failed to disable reallocation")
	}

	log.Printf("===> Launching %d instances on %s...\n", addOpts.delta, addOpts.autoScalingGroup)

	desiredCapacity, err := aws.AutoScaling.IncreaseInstances(addOpts.autoScalingGroup, addOpts.delta)
	if err != nil {
		return errors.Wrap(err, "failed to increase instance")
	}

	log.Println("===> Waiting for nodes join to Elasticsearch cluster...")

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
		time.Sleep(addSleepSeconds * time.Second)
	}

	log.Println("===> Enabling shard reallocation...")

	if err := client.EnableReallocation(); err != nil {
		return errors.Wrap(err, "failed to enable reallocation")
	}

	log.Println("===> Finished!")

	return nil
}

func init() {
	RootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVar(&addOpts.autoScalingGroup, "group", "", "Auto Scaling Group")
	addCmd.Flags().StringVar(&addOpts.clusterURL, "cluster-url", "", "Elasticsearch cluster URL")
	addCmd.Flags().IntVarP(&addOpts.delta, "number", "n", 0, "Number to add instances")
	addCmd.Flags().StringVar(&addOpts.region, "region", "", "AWS region")
}
