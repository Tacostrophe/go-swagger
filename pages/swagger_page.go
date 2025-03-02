package pages

import (
	"fmt"

	"github.com/Tacostrophe/go-swagger/usecases"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type swaggerPage struct {
	usecase   usecases.SwaggerUsecase
	pathes    []usecases.PathMethod
	textInput textinput.Model
	tip       string
}

func NewSwaggerPage(usecase usecases.SwaggerUsecase) swaggerPage {
	ti := textinput.New()
	// ti.Placeholder = "path/to/swagger.json"
	ti.Focus()
	// ti.CharLimit = 20
	ti.Width = 64
	pathes := usecase.GetFilteredPathes("")
	tip := fmt.Sprintf("found pathes %d", len(pathes))

	return swaggerPage{
		usecase:   usecase,
		pathes:    pathes,
		textInput: ti,
		tip:       tip,
	}
}

func (p swaggerPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p swaggerPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return p, tea.Quit
			// case tea.KeyRunes:
			// 	p.pathes = p.usecase.GetFilteredPathes(p.textInput.Value())
			// 	return p, nil
		}

	// We handle errors just like any other message
	case error:
		p.tip = msg.Error()
		return p, nil
	}

	p.textInput, cmd = p.textInput.Update(msg)
	p.pathes = p.usecase.GetFilteredPathes(p.textInput.Value())
	return p, cmd
}

func (p swaggerPage) View() string {
	header := fmt.Sprintf("Filter: %s", p.textInput.View())
	footer := fmt.Sprintf("%s\n(esq to quit)\n", p.tip)

	body := ""
	if len(p.pathes) == 0 {
		body = "no pathes that suits filter found"
	} else {
		pathes := p.pathes
		if len(pathes) > 10 {
			pathes = pathes[:10]
		}
		for _, path := range pathes {
			body += fmt.Sprintf("[ ] %s %s\n", path.Method, path.Path)
		}
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		header,
		body,
		footer,
	)
}
