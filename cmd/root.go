package cmd

import (
	"fmt"
	"os"

	"github.com/inabajunmr/anyconf/config"
	_ "github.com/inabajunmr/anyconf/statik"
	"github.com/inabajunmr/anyconf/vim"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "anyconf",
		Short: "anyconf open any config file of any tools.",
		Long:  `anyconf open any config file of any tools.`,
		Run: func(cmd *cobra.Command, args []string) {
			c, _ := config.ReadConfig()
			for _, v := range args {
				n, err := c.Read(v)
				if err != nil {
					fmt.Println("No config matched.")
					os.Exit(1)
				}
				c = n
			}

			if c.TargetConfigPath != "" {
				vim.LaunchVim(c.TargetConfigPath)
			} else {
				// TODO search more
			}
		},
	}
)

// Execute is just root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
