package main

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "r":
			rows, err := scanDirs()
			if err != nil {
				panic(err)
			}
			m.table.SetRows(rows)
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.table.Cursor()
			rows := m.table.Rows()
			rows[selected][2] = "yes"
			err := os.RemoveAll(d + "/" + rows[selected][0])
			if err != nil {
				panic(err)
			}
			m.table.SetRows(rows)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func scanDirs() ([]table.Row, error) {
	targets := []string{"target", "venv", "env", "node_modules"}
	rows := []table.Row{}

	return rows, filepath.WalkDir(d, func(path string, dir fs.DirEntry, err error) error {
		if slices.Contains(targets, dir.Name()) && dir.IsDir() {
			size, err := DirSize(path)
			if err != nil {
				return err
			}
			rows = append(rows, table.Row{
				strings.TrimPrefix(path, d+"/"),
				strconv.FormatInt(size>>20, 10) + "MB",
				"no",
			})
			return filepath.SkipDir
		}
		return nil
	})
}

var d = "."

func main() {
	if len(os.Args) > 1 {
		d = os.Args[1]
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	d = path.Join(wd, d)

	columns := []table.Column{
		{Title: "Name", Width: 60},
		{Title: "Size", Width: 15},
		{Title: "Deleted", Width: 8},
	}

	rows, err := scanDirs()
	if err != nil {
		panic(err)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
