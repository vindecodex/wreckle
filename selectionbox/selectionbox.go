package selectionbox

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	index    int
	selected string
	options  []string
}

func New() Model {
	return Model{
		index:    0,
		selected: "",
		options: []string{
			"Red",
			"Green",
			"Blue",
		},
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "up", "k":
			if m.index > 0 {
				m.index--
			}

		case "down", "j":
			if m.index < len(m.options)-1 {
				m.index++
			}

		case "enter", " ":
			m.selected = m.options[m.index]
		}
	}

	return m, nil
}

func (m Model) View() string {
	var options []string

	for i, option := range m.options {

		if m.index == i {
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
