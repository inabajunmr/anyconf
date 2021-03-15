package cmd

import (
	"fmt"
	"os"

	"github.com/inabajunmr/anyconf/config"
	_ "github.com/inabajunmr/anyconf/statik"
	"github.com/inabajunmr/anyconf/vim"
	"gopkg.in/AlecAivazis/survey.v1"

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
					fmt.Println("No config matched.1")
					os.Exit(1)
				}
				c = n
			}

			for {
				if c.TargetConfigPath != "" {
					vim.LaunchVim(c.TargetConfigPath)
					os.Exit(0)
				} else {

					answers := struct {
						Key string
					}{}
					qs := []*survey.Question{
						{
							Name: "Key",
							Prompt: &survey.Select{
								Message: "What's next key?",
								Options: c.NextKeys(),
							},
						},
					}
					survey.Ask(qs, &answers)
					fmt.Println(answers.Key)
					n, err := c.Read(answers.Key)
					if err != nil {
						fmt.Println("No config matched.2")
						os.Exit(1)
					}
					c = n
				}
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
