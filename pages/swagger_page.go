package pages

import (
	"fmt"
	"strings"

	"github.com/Tacostrophe/go-swagger/usecases"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type swaggerPage struct {
	usecase   usecases.SwaggerUsecase
	pathes    []usecases.PathMethod
	textInput textinput.Model
	tip       string
	pathIdx   int
}

func NewSwaggerPage(usecase usecases.SwaggerUsecase) swaggerPage {
	ti := textinput.New()
	// ti.Placeholder = "path/to/swagger.json"
	ti.Focus()
	// ti.CharLimit = 20
	ti.Width = 64
	pathes := usecase.GetFilteredPathes("")
	tip := ""

	return swaggerPage{
		usecase:   usecase,
		pathes:    pathes,
		textInput: ti,
		tip:       tip,
		pathIdx:   0,
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
	filterBlock := fmt.Sprintf("Filter: %s", p.textInput.View())
	tipsBlock := fmt.Sprintf("%s\n(esq to quit)\n", p.tip)

	pathesBlock := ""
	if len(p.pathes) == 0 {
		pathesBlock = "no pathes that suits filter found"
	} else {
		pathes := p.pathes

		pathesRows := ""

		pathesPagination := "  " + strings.Repeat(".", p.pathIdx) + "x" + strings.Repeat(".", len(pathes)-p.pathIdx-1)
		pathesPerPage := 10

		var pageStartIdx int
		pageStartIdx = p.pathIdx / pathesPerPage * pathesPerPage
		var pageEndIdx int = pageStartIdx + pathesPerPage
		if len(pathes) < pageEndIdx {
			pageEndIdx = len(pathes)
		}

		pageWithPathes := pathes[pageStartIdx:pageEndIdx]

		for i, path := range pageWithPathes {
			checkMark := " "
			if p.pathIdx == i+pageStartIdx {
				checkMark = ">"
			}
			pathesRows += fmt.Sprintf("%s [ ] %s %s\n", checkMark, path.Method, path.Path)
		}

		pathesBlock = fmt.Sprintf(
			"%s\n%s",
			pathesRows,
			pathesPagination,
		)
	}

	return fmt.Sprintf(
		"%s\n%s\n%s",
		filterBlock,
		pathesBlock,
		tipsBlock,
	)
}
