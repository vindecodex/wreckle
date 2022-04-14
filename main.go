package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vindecodex/wreckle/inputbox"
	"github.com/vindecodex/wreckle/selectionbox"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	SET_WORDS string = "SET_WORDS"
	SET_COLOR        = "SET_COLOR"
)

type model struct {
	red   string
	blue  string
	green string

	status string

	row int
	col int

	inputBoxes   [6][5]inputbox.Model
	selectionBox selectionbox.Model
}

func initialModel() model {
	return model{
		red:   "",
		blue:  "",
		green: "",

		status: SET_WORDS,

		row: 0,
		col: 0,

		inputBoxes:   generateInputBoxes(),
		selectionBox: selectionbox.New(),
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

		case "ctrl+r":
			return m.reset(), nil

		case "backspace":
			if m.col > 0 {
				m.inputBoxes[m.row][m.col-1].Value = 32
				m.col--
			}
			return m, nil

		case "up", "k", "down", "j":
			m.selectionBox, _ = m.selectionBox.Update(msg)
			return m, nil

		case "enter":
			maxCol := len(m.inputBoxes[m.row]) - 1
			if m.row < len(m.inputBoxes)-1 && m.inputBoxes[m.row][maxCol].Value != 32 {
				m.row++
				m.col = 0
				m.status = SET_COLOR
			}
			return m, nil

		default:
			if m.status == SET_WORDS {
				if m.col < len(m.inputBoxes[m.row]) {
					mdl, _ := m.inputBoxes[m.row][m.col].Update(msg)
					m.inputBoxes[m.row][m.col] = mdl
					m.col++
				}
			}
			return m, nil
		}

	}

	return m, nil
}

func (m model) View() string {
	var rowContainer []string
	var grid string
	var colorSelection string

	logo := logo() + "\n\n"

	if m.status == SET_WORDS {

		for row := 0; row < len(m.inputBoxes); row++ {
			for col := 0; col < len(m.inputBoxes[m.row]); col++ {
				rowContainer = append(rowContainer, m.inputBoxes[row][col].View())
			}

			grid += alignHorizontal(rowContainer...) + "\n"
			rowContainer = nil
		}
	}

	if m.status == SET_COLOR {
		for col := 0; col < len(m.inputBoxes[m.row]); col++ {
			rowContainer = append(rowContainer, m.inputBoxes[m.row-1][col].View())
		}
		grid += alignHorizontal(rowContainer...) + "\n"
		colorSelection = m.selectionBox.View()
		rowContainer = nil
	}

	help := "\nPress ctrl + c to quit.\n"
	help += "Press ctrl + r to restart.\n"

	width, height, _ := terminal.GetSize(0)

	container := alignVertical(logo, grid, colorSelection, help)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, container)
}

func (m model) reset() model {
	m.inputBoxes = generateInputBoxes()
	m.row = 0
	m.col = 0
	m.red = ""
	m.green = ""
	m.blue = ""
	m.status = SET_WORDS
	return m
}

func alignHorizontal(boxes ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom, boxes...)
}

func alignVertical(components ...string) string {
	return lipgloss.JoinVertical(lipgloss.Center, components...)
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

func logo() string {
	return `
╔╦═╦╗─────╔╗
║║║║╠╦╦═╦═╣╠╦╗╔═╗
║║║║║╔╣╩╣═╣═╣╚╣╩╣
╚═╩═╩╝╚═╩═╩╩╩═╩═╝
	   v1.0.0
	`
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Whoops! something went wrong", err)
		os.Exit(1)
	}
}
