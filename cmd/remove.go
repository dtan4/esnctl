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

func doRemove(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return errors.New("cluster URL and node name must be specified")
	}
	clusterURL := args[0]
	nodeName := args[1]

	client, err := es.New(clusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create Elasitcsearch API client")
	}

	// TODO: remove instance from ELB

	// TODO: wait for connection draining

	if err := client.ExcludeNodeFromAllocation(nodeName); err != nil {
		return errors.Wrap(err, "failed to exclude node from allocation group")
	}

	retryCount := 0

	for {
		shards, err := client.ListShardsOnNode(nodeName)
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
		time.Sleep(removeSleepSeconds)
	}

	if err := client.Shutdown(nodeName); err != nil {
		return errors.Wrap(err, "failed to shutdown node")
	}

	// TODO: detach instance from ASG

	return nil
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
