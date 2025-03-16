package pages

import (
	"fmt"
	"strings"

	"github.com/Tacostrophe/go-swagger/usecases"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type swaggerPage struct {
	usecase usecases.SwaggerUsecase

	pathes       []usecases.PathMethod
	chosenPathes map[usecases.PathMethod]bool
	pathIdx      int

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
	chosenPathes := make(map[usecases.PathMethod]bool, len(pathes))
	tip := ""

	return swaggerPage{
		usecase: usecase,

		pathes:       pathes,
		chosenPathes: chosenPathes,
		pathIdx:      0,

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
		case tea.KeyCtrlN:
			if p.pathIdx >= len(p.pathes)-1 {
				p.pathIdx = 0
				return p, nil
			}
			p.pathIdx++
			return p, nil
		case tea.KeyCtrlP:
			if p.pathIdx <= 0 {
				p.pathIdx = len(p.pathes) - 1
				return p, nil
			}
			p.pathIdx--
			return p, nil
		case tea.KeyCtrlL:
			path := p.pathes[p.pathIdx]
			if isChosen, ok := p.chosenPathes[path]; isChosen && ok {
				p.chosenPathes[path] = false
			} else {
				p.chosenPathes[path] = true
			}
			return p, nil
		}

	// We handle errors just like any other message
	case error:
		p.tip = msg.Error()
		return p, nil
	}

	p.textInput, cmd = p.textInput.Update(msg)
	p.pathes = p.usecase.GetFilteredPathes(p.textInput.Value())
	if len(p.pathes) == 0 {
		p.pathIdx = 0
	} else if p.pathIdx >= len(p.pathes) {
		p.pathIdx = len(p.pathes) - 1
	}

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
			pointingMark := " "
			checkMark := "[ ]"

			if p.pathIdx == i+pageStartIdx {
				pointingMark = ">"
			}
			if isChosen, ok := p.chosenPathes[path]; isChosen && ok {
				checkMark = "[X]"
			}
			pathesRows += fmt.Sprintf("%s %s %s %s\n", pointingMark, checkMark, path.Method, path.Path)
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
