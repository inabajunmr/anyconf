package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/inabajunmr/anyconf/statik"

	"github.com/rakyll/statik/fs"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "anyconf",
		Short: "anyconf open any config file of any tools.",
		Long:  `anyconf open any config file of any tools.`,
		Run: func(cmd *cobra.Command, args []string) {
			statikFS, err := fs.New()
			if err != nil {
				log.Fatal(err)
			}

			r, err := statikFS.Open("/configs.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer r.Close()
			contents, err := ioutil.ReadAll(r)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(contents))
			os.Exit(1)
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
