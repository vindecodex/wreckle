package inputbox

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Value byte
	Color string
}

func initialModel() Model {
	return Model{
		Value: 32,
		Color: "white",
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		default:
			m.Value = []byte(strings.ToUpper(msg.String()))[0]
			return m, nil
		}
	}

	return m, nil
}

func (m Model) View() string {

	return m.inputBox()
}

func (m Model) inputBox() string {
	return lipgloss.NewStyle().
		Bold(true).
		PaddingLeft(2).
		PaddingRight(2).
		Foreground(lipgloss.Color("#FFFFFF")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FFFFFF")).
		Render(string(m.Value))
}

func New() Model {
	return Model{
		Value: 32,
		Color: "white",
	}
}
