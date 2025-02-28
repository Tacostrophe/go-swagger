package pages

import (
	"fmt"

	EP "github.com/Tacostrophe/go-swagger/extract_pathes"
	RS "github.com/Tacostrophe/go-swagger/read_swagger"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type initialPage struct {
	textInput textinput.Model
	tip       string
}

func NewInitialPage() initialPage {
	ti := textinput.New()
	ti.Placeholder = "path/to/swagger.json"
	ti.Focus()
	// ti.CharLimit = 20
	ti.Width = 32
	return initialPage{
		textInput: ti,
		tip:       "",
	}
}

func (m initialPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m initialPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// m.tip = fmt.Sprintf("Looking for %q", m.textInput.Value())
			swaggerPathes, _, err := RS.ReadSwagger(m.textInput.Value())
			if err != nil {
				m.tip = fmt.Sprintf("error: can't read swagger: %s", err.Error())
				m.textInput.Reset()
				return m, nil
			}

			// extract pathes from swagger
			pathesMethodes, err := EP.ExtractPathes(swaggerPathes)
			if err != nil {
				m.tip = fmt.Sprintf("error: can't extract pathes: %s", err.Error())
				m.textInput.Reset()
				return m, nil
			}

			// transform array of pathes to human readable list with idx
			// pathes, err := TS.TransformPathesToString(pathesMethodes)
			// if err != nil {
			// 	m.tip = fmt.Sprintf("error: can't transform pathes to string: %s", err.Error())
			// 	m.textInput.Reset()
			// 	return m, nil
			// }

			return NewSwaggerPage(pathesMethodes), nil
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
}

func (m initialPage) View() string {
	return fmt.Sprintf(
		"Enter path to swagger\n\n%s\n\n%s\n%s\n",
		m.textInput.View(),
		m.tip,
		"(esc to quit)",
	)
}
