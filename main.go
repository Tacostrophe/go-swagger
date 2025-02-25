package main

import (
	"fmt"
	"os"

	"github.com/Tacostrophe/go-swagger/pages"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := pages.NewInitialPage()
	program := tea.NewProgram(model)
	if _, err := program.Run(); err != nil {
		fmt.Printf("error occured: %v", err)
		os.Exit(1)
	}
}
