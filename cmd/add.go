package cmd

import (
	"fmt"
	"time"

	"github.com/dtan4/esnctl/es"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "add",
	Short:         "Add Elasticsearch node",
	RunE:          doAdd,
}

func doAdd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("cluster URL must be specified")
	}
	clusterURL := args[0]

	client, err := es.New(clusterURL)
	if err != nil {
		return errors.Wrap(err, "failed to create Elasitcsearch API client")
	}

	if err := client.DisableReallocation(); err != nil {
		return errors.Wrap(err, "failed to disable reallocation")
	}
	defer client.EnableReallocation()

	fmt.Println("TODO: Add node here")

	time.Sleep(10 * time.Second)

	return nil
}

func init() {
	RootCmd.AddCommand(addCmd)
}
