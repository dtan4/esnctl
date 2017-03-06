package cmd

import (
	"fmt"

	"github.com/dtan4/esnctl/es"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "list URL",
	Short:         "List nodes",
	RunE:          doList,
}

func doList(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("cluster URL must be specified")
	}
	clusterURL := args[0]

	client, err := es.New(clusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create Elasitcsearch API client")
	}

	nodes, err := client.ListNodes()
	if err != nil {
		return errors.Wrap(err, "failed to list Elasticsearch nodes")
	}

	for _, node := range nodes {
		fmt.Println(node)
	}

	return nil
}

func init() {
	RootCmd.AddCommand(listCmd)
}
