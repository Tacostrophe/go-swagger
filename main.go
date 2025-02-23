package main

import (
	"fmt"
	"os"

	"github.com/Tacostrophe/go-swagger/entities"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := entities.NewModel()
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Printf("error occured: %v", err)
		os.Exit(1)
	}

}
