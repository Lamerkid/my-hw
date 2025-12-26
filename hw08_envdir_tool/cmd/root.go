/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-envdir",
	Short: "Runs another program with environment modified according to files in a specified directory",
	Long: ` d is a single argument.	

envdir sets various environment variables as specified by files in the directory named d.  It then runs child.

If d contains a file named s whose first line is t, envdir removes an environment variable named s if one exists, and then adds an environ-
ment variable named s with value t.  The name s must not contain =. Spaces and tabs at the end of t are removed. Nulls in t are changed	to
newlines in the environment variable.

If  the file s is completely empty (0 bytes long), envdir removes an environment variable named s if one exists, without adding a new vari-
able.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hw08_envdir_tool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
