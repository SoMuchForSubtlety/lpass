package ui

import (
	"fmt"
	"io"
	"os"

	"github.com/SoMuchForSubtlety/lpass/pkg/store"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

// TODO: move to util
type CompactEntryList struct {
	List     list.Model
	Items    []store.Entry
	choice   *store.Entry
	quitting bool
}

func (m *CompactEntryList) Init() tea.Cmd {
	return nil
}

func (m *CompactEntryList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.List.SelectedItem().(store.Entry)
			if ok {
				m.choice = &i
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m *CompactEntryList) View() string {
	if m.choice != nil {
		return ""
	}
	if m.quitting {
		return ""
	}
	return "\n" + m.List.View()
}

type conpactListEntry struct{}

func (d conpactListEntry) Height() int                               { return 1 }
func (d conpactListEntry) Spacing() int                              { return 0 }
func (d conpactListEntry) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d conpactListEntry) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	entry, ok := listItem.(store.Entry)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + entry.Name)
		}
	}

	fmt.Fprintf(w, fn(entry.Name))
}

func Select(entries []store.Entry) *store.Entry {
	var items []list.Item
	for _, e := range entries {
		items = append(items, e)
	}

	const defaultWidth = 20
	const listHeight = 14

	entryList := list.New(items, conpactListEntry{}, defaultWidth, listHeight)
	entryList.SetShowStatusBar(false)
	entryList.SetFilteringEnabled(false)
	entryList.SetShowTitle(false)
	entryList.SetShowHelp(false)
	entryList.SetShowPagination(false)
	m := &CompactEntryList{List: entryList}
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	return m.choice
}
