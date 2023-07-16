package cmd

import (
	"github.com/spf13/cobra"
)

var DBConnection *string

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse the PSGC .xlsx file",
	Long:  ``,
	Run:   parse,
}

func parse(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	filePath := args[0]
}

func init() {
	rootCmd.AddCommand(parseCmd)

	DBConnection = parseCmd.Flags().String("connection", "", "PostgreSQL connection string")
}
