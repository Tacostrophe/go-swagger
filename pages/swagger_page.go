package pages

import (
	"fmt"

	"github.com/Tacostrophe/go-swagger/structs"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type swaggerPage struct {
	pathes    []structs.PathMethod
	textInput textinput.Model
	tip       string
}

func NewSwaggerPage(pathes []structs.PathMethod) swaggerPage {
	ti := textinput.New()
	// ti.Placeholder = "path/to/swagger.json"
	ti.Focus()
	// ti.CharLimit = 20
	ti.Width = 64
	return swaggerPage{
		pathes:    pathes,
		textInput: ti,
		tip:       "",
	}
}

func (m swaggerPage) Init() tea.Cmd {
	return textinput.Blink
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
	return fmt.Sprintf(
		"Filter: %s\n\n%s\n%s\n",
		m.textInput.View(),
		m.tip,
		"(esc to quit)",
	)
}
