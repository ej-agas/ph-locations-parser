package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ej-agas/psgc-publication-parser/psgc"
	"os"
)

type responseMsg struct{}

func listenForActivity(m Model) tea.Cmd {

	return func() tea.Msg {

		if err := m.parser.Run(m.sub); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		m.done = true

		return nil
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

type Model struct {
	sub     chan struct{}
	parser  psgc.Parser
	spinner spinner.Model
	done    bool
	count   int
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m),   // generate activity
		waitForActivity(m.sub), // wait for activity
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.done = true
		return m, tea.Quit
	case responseMsg:
		if m.count >= len(m.parser.Rows)-1 {
			m.count++
			m.done = true
			return m, tea.Quit
		}

		m.count++
		return m, waitForActivity(m.sub)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m Model) View() string {
	totalCount := len(m.parser.Rows)

	s := fmt.Sprintf("\n %s Saved: %d/%d\n\n", m.spinner.View(), m.count, totalCount)
	if m.done {
		s += "\n"
	}
	return s
}
