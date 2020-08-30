package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Execute(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
