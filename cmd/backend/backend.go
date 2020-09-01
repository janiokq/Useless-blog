package main

import (
	"github.com/janiokq/Useless-blog/api/backend"
	"github.com/janiokq/Useless-blog/cmd/cmd"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd *cobra.Command
	rootCmd = &cobra.Command{
		Use:   "janiokq",
		Short: "Useless-blog backend module",
		Long:  `Useless-blog distributed backend module`,
		Run: func(cmd *cobra.Command, args []string) {
			//  Do Stuff Here
		},
	}
	cmd.Execute(rootCmd)
	backend.Run()
}
