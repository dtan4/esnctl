package cmd

import (
	"github.com/dtan4/esnctl/es"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
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

	// TODO: wait all shard are escaped from the given node

	if err := client.Shutdown(nodeName); err != nil {
		return errors.Wrap(err, "failed to shutdown node")
	}

	// TODO: detach instance from ASG

	return nil
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
