package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	vaultDir    string
	docStyle    = lipgloss.NewStyle().Margin(1, 2)
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error Getting Home directory.", err)
	}

	vaultDir = fmt.Sprintf("%s/.terminal-note", homeDir)
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	newFileInput           textinput.Model
	createFileInputVisible bool
	currentFile            *os.File
	noteTextArea           textarea.Model
	statusMsg              string
	list                   list.Model
	showingList            bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+n":
			m.createFileInputVisible = true
			return m, nil

		case "ctrl+l":
			noteList := listFile()
			m.list.SetItems(noteList)
			m.showingList = true
			return m, nil

		case "ctrl+s":

			if m.currentFile == nil {
				break
			}

			if err := m.currentFile.Truncate(0); err != nil {
				fmt.Println("Can not save the file :( ")
				return m, nil
			}

			if _, err := m.currentFile.Seek(0, 0); err != nil {
				fmt.Println("Can not save the file :( ")
				return m, nil
			}

			if _, err := m.currentFile.WriteString(m.noteTextArea.Value()); err != nil {
				fmt.Println("Can not save the file :( ")
				return m, nil
			}

			if err := m.currentFile.Close(); err != nil {
				fmt.Println("can not close the file.")
			}

			m.currentFile = nil
			m.noteTextArea.SetValue("")

			return m, nil

		case "enter":
			if m.currentFile != nil {
				break
			}

			if m.showingList {
				selectedItem, ok := m.list.SelectedItem().(item)
				if ok {
					filepath := fmt.Sprintf("%s/%s", vaultDir, selectedItem.title)
					content, err := os.ReadFile(filepath)
					if err != nil {
						log.Printf("Error while reading file %v", err)
						return m, nil
					}

					m.noteTextArea.SetValue(string(content))

					file, err := os.OpenFile(filepath, os.O_RDWR, 0644)
					if err != nil {
						log.Printf("Error while reading file %v", err)
						return m, nil
					}

					m.currentFile = file
					m.showingList = false
				}
				return m, nil
			}

			if m.createFileInputVisible {
				filename := m.newFileInput.Value()
				if filename != "" {
					filepath := filepath.Join(vaultDir, fmt.Sprintf("%s.md", filename))

					if _, err := os.Stat(filepath); err == nil {
						// File exists — show user feedback
						m.statusMsg = fmt.Sprintf("⚠️ File '%s.md' already exists!", filename)
						return m, nil
					}

					file, err := os.Create(filepath)
					if err != nil {
						m.statusMsg = fmt.Sprintf("❌ Error creating file: %v", err)
						return m, nil
					}

					m.currentFile = file
					m.createFileInputVisible = false
					m.newFileInput.SetValue("")
					m.statusMsg = fmt.Sprintf("✅ Created new note: %s.md", filename)
				}
			}

			return m, nil
		}
	}

	if m.createFileInputVisible {
		m.newFileInput, cmd = m.newFileInput.Update(msg)
	}

	if m.currentFile != nil {
		m.noteTextArea, cmd = m.noteTextArea.Update(msg)
	}

	if m.showingList {
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	welcome := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("205")).
		PaddingLeft(2).
		PaddingRight(2).
		Render("welcome to terminal-note 🧠")

	help := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("16")).
		Background(lipgloss.Color("7")).
		PaddingLeft(2).
		PaddingRight(2).
		Render("Ctrl+N: new file • Ctrl+L: list • Esc: back/save • Ctrl+S: save • Ctrl+Q: quit")

	content := ""
	if m.createFileInputVisible {
		content = m.newFileInput.View()
	} else if m.currentFile != nil {
		content = m.noteTextArea.View()
	} else if m.showingList {
		content = m.list.View()
	}

	statusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	status := ""
	if m.statusMsg != "" {
		status = statusStyle.Render(m.statusMsg)
	}

	return fmt.Sprintf("\n%s\n\n%s\n\n%s\n\n%s\n", welcome, content, help, status)
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

	//input noteTextArea
	ta := textarea.New()
	ta.Placeholder = "write your note"
	ta.Focus()
	ta.ShowLineNumbers = false

	//list
	noteList := listFile()
	finalList := list.New(noteList, list.NewDefaultDelegate(), 0, 0)
	finalList.Title = "All Notes"
	finalList.Styles.Title = lipgloss.NewStyle().Foreground(lipgloss.Color("16")).Background(lipgloss.Color("254")).Padding(0, 1)

	return model{
		newFileInput:           ti,
		createFileInputVisible: false,
		noteTextArea:           ta,
		list:                   finalList,
	}
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func listFile() []list.Item {
	items := make([]list.Item, 0)

	entries, err := os.ReadDir(vaultDir)
	if err != nil {
		log.Fatal("Error Reading notes")
	}

	for _, entry := range entries {

		if !entry.IsDir() {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			modTime := info.ModTime().Format("2005-01-02 15:04")

			items = append(items, item{
				title: entry.Name(),
				desc:  fmt.Sprintf("Modified: %s", modTime),
			})
		}
	}

	return items
}
