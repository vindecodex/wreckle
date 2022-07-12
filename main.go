package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vindecodex/wreckle/c"
	"github.com/vindecodex/wreckle/inputbox"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	MAXCOL = 5
	MAXROW = 6
)

type State interface {
	onBackspace()
	onEnter()
	handleInput(msg string)
	handleColor(msg string)
	handlePick(msg string)
	render() string
}

type InputState struct {
	m *model
}

type ColorState struct {
	m *model
}

type PickState struct {
	m *model
}

func (h InputState) onBackspace() {
	if h.m.col > 0 {
		h.m.inputBoxes[h.m.row][h.m.col-1].Value = " "
		h.m.col--
	}
}

func (h InputState) onEnter() {
	if h.m.col == MAXCOL && h.m.row < MAXROW-1 {
		h.m.row++
		h.m.col = 0
		h.m.setState(h.m.colorState)
	}
}

func (h InputState) handleInput(msg string) {
	if h.m.col < MAXCOL {
		h.m.inputBoxes[h.m.row][h.m.col].Value = strings.ToUpper(msg)
		h.m.col++
	}
}

func (h InputState) handleColor(msg string) {}

func (h InputState) handlePick(msg string) {}

func (h InputState) render() string {
	var rowContainer []string
	var gridContainer string
	logo := logo()
	for r := 0; r < MAXROW; r++ {
		for c := 0; c < MAXCOL; c++ {
			rowContainer = append(rowContainer, h.m.inputBoxes[r][c].View())
		}
		gridContainer += alignHorizontal(rowContainer...) + "\n"
		rowContainer = nil
	}

	help := "Press [ctrl+c] to exit \n"
	help += "Press [ctrl+r] to restart \n"

	app := alignVertical(logo, gridContainer, help)
	return lipgloss.Place(h.m.width, h.m.height, lipgloss.Center, lipgloss.Center, app)
}

func (h ColorState) onBackspace() {
	if h.m.col > 0 {
		h.m.inputBoxes[h.m.row-1][h.m.col-1].MetaTag = c.WHITE
		h.m.col--
	}
}

func (h ColorState) onEnter() {
	if h.m.col == MAXCOL && h.m.row < MAXROW-1 {
		h.m.col = 0
		h.m.setState(h.m.inputState)
	}
}

func (h ColorState) handleInput(msg string) {
	if h.m.col < MAXCOL {
		switch msg {
		case "x":
			h.m.inputBoxes[h.m.row-1][h.m.col].MetaTag = c.GREY
			h.m.col++
		case "g":
			h.m.inputBoxes[h.m.row-1][h.m.col].MetaTag = c.GREEN
			h.m.col++
		case "y":
			h.m.inputBoxes[h.m.row-1][h.m.col].MetaTag = c.YELLOW
			h.m.col++

		}
	}

}

func (h ColorState) handleColor(msg string) {}

func (h ColorState) handlePick(msg string) {}

func (h ColorState) render() string {
	var rowContainer []string
	var gridContainer string
	logo := logo()
	guide := "x = gray, g = green, y = yellow"
	for c := 0; c < MAXCOL; c++ {
		rowContainer = append(rowContainer, h.m.inputBoxes[h.m.row-1][c].View())
	}
	gridContainer += lipgloss.NewStyle().Render(alignHorizontal(rowContainer...))

	help := "Press [ctrl+c] to exit \n"
	help += "Press [ctrl+r] to restart \n"

	app := alignVertical(logo, guide, gridContainer, help)
	return lipgloss.Place(h.m.width, h.m.height, lipgloss.Center, lipgloss.Center, app)
}

func (h PickState) onBackspace() {}

func (h PickState) onEnter() {}

func (h PickState) handleInput(msg string) {}

func (h PickState) handleColor(msg string) {}

func (h PickState) handlePick(msg string) {}

func (h PickState) render() string {
	return ""
}

type model struct {
	inputBoxes [MAXROW][MAXCOL]inputbox.Model

	col int
	row int

	width  int
	height int

	state      State
	inputState State
	colorState State
	pickState  State
}

func initialModel() *model {
	width, height, _ := terminal.GetSize(0)
	m := &model{
		inputBoxes: generateInputBoxes(),

		row: 0,
		col: 0,

		width:  width,
		height: height,
	}

	inputState := &InputState{m}
	colorState := &ColorState{m}
	pickState := &PickState{m}

	m.state = inputState
	m.inputState = inputState
	m.colorState = colorState
	m.pickState = pickState
	return m
}

func (m *model) setState(state State) {
	// clearing the terminal to fix rendering issue
	fmt.Printf("\x1bc")
	m.state = state
}

func generateInputBoxes() [MAXROW][MAXCOL]inputbox.Model {
	inputBoxes := [MAXROW][MAXCOL]inputbox.Model{}
	for r := 0; r < MAXROW; r++ {
		for c := 0; c < MAXCOL; c++ {
			inputBoxes[r][c] = inputbox.New()
		}
	}
	return inputBoxes
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+r":
			return initialModel(), nil
		case "enter":
			m.state.onEnter()
			return m, nil
		case "backspace":
			m.state.onBackspace()
			return m, nil
		default:
			m.state.handleInput(msg.String())
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	return m.state.render()
}

func alignHorizontal(components ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Bottom, components...)
}

func alignVertical(components ...string) string {
	return lipgloss.JoinVertical(lipgloss.Center, components...)
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
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}
