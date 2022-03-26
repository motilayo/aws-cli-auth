package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"aws-cli-auth/config"
)

var (
	rootCmd = &cobra.Command{
		Use:   "aws-cli-auth-cobra",
		Short: "A tool for authenticating your AWS CLI",
		Long: `A lightweight tool that will allow an AWS IAM user to assume a role.
		It accepts a coonfig yaml file.
		It will assume a role for the max duration that the role can be assumed.
		This tool is standalone and does not require the AWS CLI.`,
	
		Run: func(cmd *cobra.Command, args []string) {
			configFilePath := cmd.Flag("config-file").Value.String()

			if configFilePath == "" {
				err:= fmt.Errorf("Config file missing. Please provide a config file")
				log.Fatal(err)
				os.Exit(1)
			} else {
				config.Configure(configFilePath)
			}
		 },
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
    rootCmd.PersistentFlags().String("config-file", "", "A yaml file containing the IAM user & role configuration details")
}