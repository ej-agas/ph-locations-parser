package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var vFlag *bool
var Version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "internal-pub-parser",
	Short: "internal-pub-parser",
	Long:  `Parses PSGC publication file and stores it to ph-locations' database.`,
	Run: func(cmd *cobra.Command, args []string) {
		if *vFlag {
			println(Version)
			return
		}

		cmd.Help()
	},
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	vFlag = rootCmd.Flags().BoolP("version", "v", false, "Show program version")
}
