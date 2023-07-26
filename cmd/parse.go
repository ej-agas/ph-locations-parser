package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ej-agas/ph-locations/postgresql"
	"github.com/ej-agas/psgc-publication-parser/psgc"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var (
	host     *string
	port     *string
	user     *string
	password *string
	db       *string
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse the PSGC .xlsx file",
	Long:  ``,
	Run:   process,
}

func process(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	dbPort, err := strconv.Atoi(*port)
	if err != nil {
		fmt.Println(fmt.Errorf("invalid port"))
		return
	}

	dbConfig := postgresql.Config{
		Host:         *host,
		Port:         dbPort,
		User:         *user,
		Password:     *password,
		DatabaseName: *db,
	}

	connection, err := postgresql.NewConnection(dbConfig)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to connect to PostgreSQL: %s", err))
		return
	}

	fmt.Println("Getting rows from file...")
	rows, err := psgc.GetRowsFromFile(args[0])

	if err != nil {
		fmt.Println(err)
		return
	}

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	s.Spinner = spinner.Points

	parser := psgc.Parser{
		State: *psgc.NewState(),
		Store: *psgc.NewStore(connection),
		Rows:  rows,
	}

	p := tea.NewProgram(Model{
		sub:     make(chan struct{}),
		parser:  parser,
		spinner: s,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(parseCmd)

	host = parseCmd.Flags().String("host", "127.0.0.1", "PostgreSQL Host")
	port = parseCmd.Flags().String("port", "5173", "PostgreSQL Port")
	user = parseCmd.Flags().String("user", "ph_locations_user", "PostgreSQL User")
	password = parseCmd.Flags().String("password", "", "PostgreSQL Password")
	db = parseCmd.Flags().String("db", "ph_locations_db", "PostgreSQL Database Name")
}
