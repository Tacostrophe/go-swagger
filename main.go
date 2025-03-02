package main

import (
	"fmt"
	"os"

	"github.com/Tacostrophe/go-swagger/pages"
	"github.com/Tacostrophe/go-swagger/usecases"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	swaggerUsecase, err := usecases.NewSwaggerFromFileV1()
	if err != nil {
		panic(err)
	}
	page := pages.NewInitialPage(swaggerUsecase)
	program := tea.NewProgram(page)
	if _, err := program.Run(); err != nil {
		fmt.Printf("error occured: %v", err)
		os.Exit(1)
	}
}
