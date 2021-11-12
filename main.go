package main

import (
	"fmt"

	command "github.com/neobaran/csac/cmd"
	"github.com/spf13/cobra"
)

var Version string

var (
	file  string
	debug bool
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "generate ssl and upload to cloud",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(file) == 0 {
			file = "config.yaml"
		}
		command.Generate(file, debug)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "查看版本",
	Long:  `查看版本`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().StringVarP(&file, "config", "c", "", "config file")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug mode")
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
