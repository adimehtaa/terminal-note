package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	msg string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {

	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingLeft(2).
		PaddingRight(2)

	welcome := style.Render("welcome to terminal-note 🧠")

	var styleHelp = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("7")).
		PaddingLeft(2).
		PaddingRight(2)

	help := styleHelp.Render("Ctrl+N: new file · Ctrl+L: list · Esc: back/save · Ctrl+S: save · Ctrl+Q: quit")

	view := ""

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n", welcome, view, help)
}

func initialModel() model {
	return model{
		msg: "hhllh",
	}
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
