package main

import (
	"github.com/janiokq/Useless-blog/cmd/cmd"
	"github.com/janiokq/Useless-blog/srv/user"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd *cobra.Command
	rootCmd = &cobra.Command{
		Use:   "janiokq",
		Short: "Useless-blog User module",
		Long:  `Useless-blog distributed User module`,
		Run: func(cmd *cobra.Command, args []string) {
			//  Do Stuff Here
		},
	}
	cmd.Execute(rootCmd)
	user.Run()
}
