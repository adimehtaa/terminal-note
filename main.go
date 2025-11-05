package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	vaultDir    string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error Getting Home directory.", err)
	}

	vaultDir = fmt.Sprint("%s/.terminal-note", homeDir)
}

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil

		case "enter":
			return m, nil
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	return m, cmd
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
	if m.createFileInputVisible {
		view = m.newFileInput.View()
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n", welcome, view, help)
}

func initialModel() model {

	err := os.MkdirAll(vaultDir, 0750)
	if err != nil {
		log.Fatal(err)
	}

	ti := textinput.New()
	ti.Placeholder = "Create Your Notes"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Cursor.Style = cursorStyle
	ti.PromptStyle = cursorStyle
	ti.TextStyle = cursorStyle

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
	}
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
