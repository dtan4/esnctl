package cmd

import (
	"fmt"
	"time"

	"github.com/dtan4/esnctl/es"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	removeMaxRetry     = 60
	removeSleepSeconds = 5
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove node from Elasticsearch cluster",
	RunE:  doRemove,
}

var removeOpts = struct {
	autoScalingGroup string
	clusterURL       string
	nodeName         string
}{}

func doRemove(cmd *cobra.Command, args []string) error {
	if removeOpts.autoScalingGroup == "" {
		return errors.New("AutoScaling Group must be specified")
	}

	if removeOpts.clusterURL == "" {
		return errors.New("Elasticsearch cluster URL must be specified")
	}

	if removeOpts.nodeName == "" {
		return errors.New("Elasticsearch Node name must be specified")
	}

	client, err := es.New(removeOpts.clusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create Elasitcsearch API client")
	}

	// TODO: remove instance from ELB

	// TODO: wait for connection draining

	if err := client.ExcludeNodeFromAllocation(removeOpts.nodeName); err != nil {
		return errors.Wrap(err, "failed to exclude node from allocation group")
	}

	retryCount := 0

	for {
		shards, err := client.ListShardsOnNode(removeOpts.nodeName)
		if err != nil {
			return errors.Wrap(err, "failed to list shards on the given node")
		}

		if len(shards) == 0 {
			fmt.Print("\n")
			break
		}

		fmt.Print(".")

		if retryCount == removeMaxRetry {
			return errors.New("shards did not escaped from the given node")
		}

		retryCount++
		time.Sleep(removeSleepSeconds * time.Second)
	}

	if err := client.Shutdown(removeOpts.nodeName); err != nil {
		return errors.Wrap(err, "failed to shutdown node")
	}

	// TODO: detach instance from ASG

	return nil
}

func init() {
	RootCmd.AddCommand(removeCmd)

	addCmd.Flags().StringVar(&removeOpts.autoScalingGroup, "group", "", "AutoScaling Group")
	addCmd.Flags().StringVar(&removeOpts.clusterURL, "cluster-url", "", "Elasticsearch cluster URL")
	addCmd.Flags().StringVar(&removeOpts.nodeName, "node-name", "", "Elasticsearch node name to remove")
}
