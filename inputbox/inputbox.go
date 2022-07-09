package inputbox

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vindecodex/wreckle/c"
)

type Model struct {
	Value   string
	MetaTag string
}

func New() Model {
	return Model{
		Value:   " ",
		MetaTag: c.WHITE,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			m.Value = strings.ToUpper(msg.String())
			return m, nil
		}
	case Erase:
		m.Value = " "
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	return lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(2).
		Bold(true).
		Foreground(c.GetColor(m.MetaTag)).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(c.GetColor(m.MetaTag)).
		Render(m.Value)
}

type Erase struct{}
