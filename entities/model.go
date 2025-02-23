package entities

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	textInput textinput.Model
	tip       string
	// cursor   int
	// rows     []string
	// selected map[int]struct{}
}

func NewModel() model {
	ti := textinput.New()
	ti.Placeholder = "path/to/swagger.json"
	ti.Focus()
	// ti.CharLimit = 20
	ti.Width = 32
	return model{
		textInput: ti,
		tip:       "",
	}
	// return model{
	// 	rows: []string{
	// 		"first one",
	// 		"second",
	// 		"",
	// 		"previous one is empty",
	// 		"and the last one",
	// 	},
	// 	selected: make(map[int]struct{}),
	// }
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.tip = fmt.Sprintf("Looking for %q", m.textInput.Value())
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		// m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd

	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// 	switch msg.String() {
	// 	case "ctrl+c", "q":
	// 		return m, tea.Quit
	// 	case "up", "k":
	// 		if m.cursor > 0 {
	// 			m.cursor--
	// 		}
	// 	case "down", "j":
	// 		if m.cursor < len(m.rows)-1 {
	// 			m.cursor++
	// 		}
	// 	case "enter", " ":
	// 		if _, ok := m.selected[m.cursor]; ok {
	// 			delete(m.selected, m.cursor)
	// 		} else {
	// 			m.selected[m.cursor] = struct{}{}
	// 		}
	// 	}
	// }
	// return m, nil
}

func (m model) View() string {
	return fmt.Sprintf(
		"Enter path to swagger\n\n%s\n\n%s\n%s\n",
		m.textInput.View(),
		m.tip,
		"(esc to quit)",
	)
	// s := "First line which i'll figure out later\n\n"

	// for i, row := range m.rows {
	// 	cursor := " "
	// 	if m.cursor == i {
	// 		cursor = ">"
	// 	}

	// 	checked := " "
	// 	if _, ok := m.selected[i]; ok {
	// 		checked = "x"
	// 	}

	// 	s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, row)
	// }

	// s += "\n q/ctrl+c - quit, spacebar/enter - choose"
	// return s
}
