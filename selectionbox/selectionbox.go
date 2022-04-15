package selectionbox

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Index    int
	Selected string
	Options  []string
}

func New() Model {
	return Model{
		Index:    0,
		Selected: "",
		Options: []string{
			"red",
			"green",
			"blue",
		},
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "up", "k":
			if m.Index > 0 {
				m.Index--
			}

		case "down", "j":
			if m.Index < len(m.Options)-1 {
				m.Index++
			}

		case "enter", " ":
			m.Selected = m.Options[m.Index]
		}
	}

	return m, nil
}

func (m Model) View() string {
	var options []string

	for i, option := range m.Options {

		if m.Index == i {
			options = append(options, lipgloss.NewStyle().Align(lipgloss.Center).Foreground(lipgloss.Color("#000000")).Background(lipgloss.Color("#ffffff")).Width(50).Render(option+"\n"))
		} else {
			options = append(options, lipgloss.NewStyle().Align(lipgloss.Center).Width(50).Render(option+"\n"))
		}

	}

	return lipgloss.NewStyle().PaddingTop(5).PaddingBottom(5).Render(alignVertical(options...))
}

func alignVertical(components ...string) string {
	return lipgloss.JoinVertical(lipgloss.Center, components...)
}
