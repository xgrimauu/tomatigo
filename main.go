package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xgrimauu/models"
	"os"
)

func main() {
	program := tea.NewProgram(models.InitModel(), tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		os.Exit(1)
	}
}
