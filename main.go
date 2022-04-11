package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vindecodex/wreckle/inputbox"
)

type model struct {
	red   string
	blue  string
	green string

	row int
	col int

	inputBoxes [6][5]inputbox.Model
}

func initialModel() model {
	return model{
		red:   "",
		blue:  "",
		green: "",

		row: 0,
		col: 0,

		inputBoxes: generateInputBoxes(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "backspace":
			if m.col > 0 {
				m.inputBoxes[m.row][m.col-1].Value = 32
				m.col--
			}
			return m, nil

		case "enter":
			maxCol := len(m.inputBoxes[m.row]) - 1
			if m.row < len(m.inputBoxes)-1 && m.inputBoxes[m.row][maxCol].Value != 32 {
				m.row++
				m.col = 0
			}
			return m, nil

		default:
			if m.col < len(m.inputBoxes[m.row]) {
				mdl, _ := m.inputBoxes[m.row][m.col].Update(msg)
				m.inputBoxes[m.row][m.col] = mdl
				m.col++
			}
			return m, nil
		}

	}

	return m, nil
}

func (m model) View() string {
	var rowContainer []string
	s := "Input your wordle guess here! \n\n"

	for row := 0; row < len(m.inputBoxes); row++ {
		for col := 0; col < len(m.inputBoxes[m.row]); col++ {
			rowContainer = append(rowContainer, m.inputBoxes[row][col].View())
		}

		s += alignVertical(rowContainer) + "\n"
		rowContainer = nil
	}

	s += "\nPress ctrl + c to quit.\n"

	return lipgloss.JoinVertical(lipgloss.Center, s)
}

func alignVertical(box []string) string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom, box...)
}

func generateInputBoxes() [6][5]inputbox.Model {
	ib := [6][5]inputbox.Model{}
	for row := 0; row < len(ib); row++ {
		for col := 0; col < len(ib[row]); col++ {
			ib[row][col] = inputbox.New()
		}
	}
	return ib
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Whoops! something went wrong", err)
		os.Exit(1)
	}
}
