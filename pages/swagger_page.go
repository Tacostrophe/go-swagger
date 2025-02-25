package pages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type swaggerPage struct {
	pathes string
}

func NewSwaggerPage(pathes string) swaggerPage {
	return swaggerPage{pathes: pathes}
}

func (m swaggerPage) Init() tea.Cmd {
	return nil
}

func (m swaggerPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		// m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m swaggerPage) View() string {
	return m.pathes
}
