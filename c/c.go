package c

import "github.com/charmbracelet/lipgloss"

const (
	WHITE  string = "white"
	GREY   string = "grey"
	GREEN  string = "green"
	YELLOW string = "yellow"
)

func GetColor(color string) lipgloss.Color {
	switch color {
	case GREY:
		return lipgloss.Color("#636363")
	case GREEN:
		return lipgloss.Color("#00FF33")
	case YELLOW:
		return lipgloss.Color("#EDDA05")
	default:
		return lipgloss.Color("#FFFFFF")

	}
}
