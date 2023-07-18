package cmd

import (
	"fmt"
	"github.com/ej-agas/psgc-publication-parser/psgc"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
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
	file, err := excelize.OpenFile(filePath)

	if err != nil {
		fmt.Println(fmt.Errorf("error opening file: %s", err))
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(fmt.Errorf("error closing file: %s", err))
		}
	}()

	rows, err := file.GetRows("PSGC")
	if err != nil {
		fmt.Println(fmt.Errorf("error getting rows: %s", err))
		return
	}

	rowCount := 1
	for _, item := range rows {
		if rowCount == 1 {
			rowCount++
			continue
		}

		if rowCount == 6 {
			return
		}

		fmt.Printf("%#v\n", psgc.NewRow(item))
		rowCount++
	}
}

func init() {
	rootCmd.AddCommand(parseCmd)

	DBConnection = parseCmd.Flags().String("connection", "", "PostgreSQL connection string")
}
