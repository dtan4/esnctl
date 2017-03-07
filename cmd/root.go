package cmd

import (
	"fmt"
	"os"

	"github.com/dtan4/esnctl/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "esnctl",
	Short: "A brief description of your application",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := aws.Initialize(rootOpts.region); err != nil {
			return errors.Wrap(err, "failed to initialize AWS service clients")
		}

		return nil
	},
}

var rootOpts = struct {
	region string
}{}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if trace := os.Getenv("TRACE"); trace == "1" {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&rootOpts.region, "region", "", "AWS region")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
