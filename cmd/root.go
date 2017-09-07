package cmd

import "github.com/spf13/cobra"

func init() {}

func Execute() {
	RootCmd.Execute()
}

var RootCmd = &cobra.Command{
	Use:   "influxc",
	Short: "influxc is simple cluster management command.",
	Long: `This is a simple cluster management command for influxdb
	       in Go. Internal, use docker library, but no need extra 
				 library to execute.  `,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
