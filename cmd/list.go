package cmd

import (
	"fmt"
	"net/http"

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

var listOpts = struct {
	clusterURL string
}{}

func doList(cmd *cobra.Command, args []string) error {
	if listOpts.clusterURL == "" {
		return errors.New("Elasticsearch cluster (--cluster-url) must be specified")
	}

	httpClient := &http.Client{}

	client, err := es.New(listOpts.clusterURL, httpClient)
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

	listCmd.Flags().StringVar(&listOpts.clusterURL, "cluster-url", "", "Elasticsearch cluster URL")
}
