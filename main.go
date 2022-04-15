package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vindecodex/wreckle/inputbox"
	"github.com/vindecodex/wreckle/selectionbox"
	"github.com/vindecodex/wreckle/utils"
	"golang.org/x/crypto/ssh/terminal"
)

type model struct {
	red   string
	blue  string
	green string

	status string

	row int
	col int

	width  int
	height int

	inputBoxes   [utils.ROW][utils.COL]inputbox.Model
	selectionBox selectionbox.Model
}

func initialModel() model {
	width, height, _ := terminal.GetSize(0)
	return model{
		red:   "",
		blue:  "",
		green: "",

		status: utils.SET_WORDS,

		row: 0,
		col: 0,

		width:  width,
		height: height,

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
			if m.status == utils.SET_WORDS {
				maxCol := len(m.inputBoxes[m.row]) - 1
				if m.row < len(m.inputBoxes)-1 && m.inputBoxes[m.row][maxCol].Value != 32 {
					m.row++
					m.col = 0
					m.status = utils.SET_COLOR
				}
			}
			if m.status == utils.SET_COLOR {
				m.inputBoxes[m.row][m.col].Color = m.selectionBox.Selected
				m.inputBoxes[m.row][m.col] = m.inputBoxes[m.row][m.col]
			}
			return m, nil

		default:
			if m.status == utils.SET_WORDS {
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

	if m.status == utils.SET_WORDS {

		for row := 0; row < len(m.inputBoxes); row++ {
			for col := 0; col < len(m.inputBoxes[m.row]); col++ {
				rowContainer = append(rowContainer, m.inputBoxes[row][col].View())
			}

			grid += alignHorizontal(rowContainer...) + "\n"
			rowContainer = nil
		}
	}

	if m.status == utils.SET_COLOR {
		for col := 0; col < len(m.inputBoxes[m.row]); col++ {
			rowContainer = append(rowContainer, m.inputBoxes[m.row-1][col].View())
		}
		grid += alignHorizontal(rowContainer...) + "\n"
		colorSelection = m.selectionBox.View()
		rowContainer = nil
	}

	help := "\nPress ctrl + c to quit.\n"
	help += "Press ctrl + r to restart.\n"

	container := alignVertical(logo, grid, colorSelection, help)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, container)
}

func (m model) reset() model {
	width, height, _ := terminal.GetSize(0)
	m.inputBoxes = generateInputBoxes()
	m.row = 0
	m.col = 0
	m.red = ""
	m.green = ""
	m.blue = ""
	m.status = utils.SET_WORDS
	m.width = width
	m.height = height
	return m
}

func alignHorizontal(boxes ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom, boxes...)
}

func alignVertical(components ...string) string {
	return lipgloss.JoinVertical(lipgloss.Center, components...)
}

func generateInputBoxes() [utils.ROW][utils.COL]inputbox.Model {
	ib := [utils.ROW][utils.COL]inputbox.Model{}
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
