package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/inabajunmr/anyconf/config"
	"github.com/inabajunmr/anyconf/editor"
	_ "github.com/inabajunmr/anyconf/statik"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "anyconf",
		Short: "anyconf open any config file of any tools.",
		Long:  `anyconf open any config file of any tools.`,
		Run: func(cmd *cobra.Command, args []string) {
			c, _ := config.ReadConfigs()
			viper.ReadInConfig()
			e := viper.GetString("editor")

			for _, v := range args {
				n, err := c.Read(v)
				if err != nil {
					fmt.Println(contributionAd(v))
					os.Exit(0)
				}
				c = n
			}

			for {
				if c.TargetConfigPath != "" {
					path := config.GetPath(c.TargetConfigPath)
					if _, err := os.Stat(path); os.IsNotExist(err) {
						// if there no file, show prompt
						create := false
						prompt := &survey.Confirm{
							Message: fmt.Sprintf("Create %v?", c.TargetConfigPath),
						}
						survey.AskOne(prompt, &create, nil)
						if !create {
							os.Exit(0)
						}

						// create new file
						os.MkdirAll(filepath.Dir(path), 0755)
						_, err := os.Create(path)
						if err != nil {
							fmt.Println("Failed to create new file.")
							os.Exit(1)
						}
					}
					editor.LaunchEditor(path, e)
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
						fmt.Println(contributionAd(answers.Key))
						os.Exit(0)
					}
					c = n
				}
			}
		},
	}
)

func init() {
	conf, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(conf, ".anyconf", "config.yml")
	viper.SetConfigFile(path)
	viper.SetDefault("editor", "vim")

}

func contributionAd(val string) string {
	return fmt.Sprintf("anyconf doesn't support %v yet. \nYou can support %v at https://github.com/inabajunmr.", val, val)
}

// Execute is just root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
