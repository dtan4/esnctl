package cmd

import (
	"fmt"

	"github.com/dtan4/esnctl/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	RunE:  doVersion,
}

func doVersion(cmd *cobra.Command, args []string) error {
	fmt.Println(version.String())

	return nil
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
