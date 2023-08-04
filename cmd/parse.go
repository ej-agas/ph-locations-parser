package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ej-agas/ph-locations-parser/internal"
	"github.com/ej-agas/ph-locations/postgresql"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strconv"
)

var (
	host         *string
	port         *string
	user         *string
	password     *string
	db           *string
	passwordFlag bool
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

	if passwordFlag == false {
		fmt.Println("no database password provided")
		os.Exit(1)
	}

	if err := promptPassword(); err != nil {
		fmt.Println(fmt.Errorf("error parsing password: %w", err))
		os.Exit(1)
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

	fmt.Println("\nGetting rows from file...")
	rows, err := internal.GetRowsFromFile(args[0])

	if err != nil {
		fmt.Println(err)
		return
	}

	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	s.Spinner = spinner.Points

	parser := internal.Parser{
		State: *internal.NewState(),
		Store: *internal.NewStore(connection),
		Rows:  rows,
	}

	p := tea.NewProgram(Model{
		sub:     make(chan struct{}),
		parser:  parser,
		spinner: s,
	})

	if _, err := p.Run(); err != nil {
		fmt.Println(fmt.Errorf("could not start program: %w", err))
		os.Exit(1)
	}

	style := lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.Color("#50C878"))

	fmt.Println(style.Render("All done!"))
}

func init() {
	rootCmd.AddCommand(parseCmd)

	host = parseCmd.Flags().String("host", "127.0.0.1", "PostgreSQL Host")
	port = parseCmd.Flags().String("port", "5432", "PostgreSQL Port")
	user = parseCmd.Flags().String("user", "ph_locations_user", "PostgreSQL User")
	parseCmd.Flags().BoolVar(&passwordFlag, "password", false, "PostgreSQL password")
	db = parseCmd.Flags().String("db", "ph_locations_db", "PostgreSQL Database Name")
}

func promptPassword() error {
	fmt.Print("Enter password: ")

	passwordBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		return err
	}

	pass := string(passwordBytes)
	password = &pass

	return nil
}
