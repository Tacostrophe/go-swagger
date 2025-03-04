package pages

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/Tacostrophe/go-swagger/usecases"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type initialPage struct {
	usecase       usecases.SwaggerUsecase
	textInput     textinput.Model
	tip           string
	suggestionIdx *int
	lastInput     string
}

func NewInitialPage(usecase usecases.SwaggerUsecase) initialPage {
	ti := textinput.New()
	ti.Placeholder = "path/to/swagger.json"
	// ti.ShowSuggestions = true
	ti.Focus()
	// ti.CharLimit = 20
	ti.Width = 32
	tip := ""
	suggestions, err := getSuggestions("")
	if err != nil {
		tip = err.Error()
	} else {
		ti.SetSuggestions(suggestions)
	}

	return initialPage{
		usecase:       usecase,
		textInput:     ti,
		tip:           tip,
		suggestionIdx: nil,
		lastInput:     "",
	}
}

func (m initialPage) Init() tea.Cmd {
	return textinput.Blink
}

func getSuggestions(input string) ([]string, error) {
	pwd := "./"

	pathRegexp := regexp.MustCompile(`^(?P<dir>.+/)?(?P<file>[^/]*)$`)
	pathMatches := pathRegexp.FindStringSubmatch(input)
	dirIndex := pathRegexp.SubexpIndex("dir")
	fileIndex := pathRegexp.SubexpIndex("file")

	dirPath := pathMatches[dirIndex]
	fileName := pathMatches[fileIndex]

	if len(dirPath) != 0 && !strings.HasSuffix(dirPath, "/") {
		dirPath += "/"
	}
	path := fmt.Sprintf("%s%s", pwd, dirPath)
	f, err := os.Open(path)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()
	files, err := f.ReadDir(0)
	if err != nil {
		return []string{}, err
	}

	suggestions := make([]string, 0, len(files))
	for _, v := range files {
		if strings.HasPrefix(v.Name(), fileName) {
			suggestion := fmt.Sprintf("%s%s", dirPath, v.Name())
			suggestions = append(suggestions, suggestion)
		}
	}
	return suggestions, nil
}

func (m initialPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			suggestions := m.textInput.AvailableSuggestions()
			if len(suggestions) == 0 {
				m.tip = "no file suits provided path"
				return m, nil
			}

			if err := m.usecase.Init(m.textInput.Value()); err != nil {
				m.textInput.Reset()

				m.tip = fmt.Sprintf("error: can't parse file: %s", err.Error())
				return m, nil
			}

			return NewSwaggerPage(m.usecase), nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyCtrlP:
			suggestions := m.textInput.AvailableSuggestions()
			if len(suggestions) == 0 {
				m.tip = "no file suits provided path"
				return m, nil
			}

			if m.suggestionIdx == nil {
				m.suggestionIdx = new(int)
				*m.suggestionIdx = len(suggestions) - 1
				m.textInput.SetValue(suggestions[*m.suggestionIdx])
				m.textInput.CursorEnd()

				m.tip = ""
				return m, nil
			}

			if *m.suggestionIdx == 0 {
				m.suggestionIdx = nil
				m.textInput.SetValue(m.lastInput)
				m.textInput.CursorEnd()

				m.tip = ""
				return m, nil
			}

			*m.suggestionIdx--
			m.textInput.SetValue(suggestions[*m.suggestionIdx])
			m.textInput.CursorEnd()

			m.tip = ""
			return m, nil

		case tea.KeyCtrlN:
			suggestions := m.textInput.AvailableSuggestions()
			if len(suggestions) == 0 {
				m.tip = "no file suits provided path"

				return m, nil
			}

			if m.suggestionIdx == nil {
				m.suggestionIdx = new(int)
				*m.suggestionIdx = 0
				m.textInput.SetValue(suggestions[*m.suggestionIdx])
				m.textInput.CursorEnd()

				m.tip = ""
				return m, nil
			}

			if *m.suggestionIdx == len(suggestions)-1 {
				m.suggestionIdx = nil
				m.textInput.SetValue(m.lastInput)
				m.textInput.CursorEnd()

				m.tip = ""
				return m, nil
			}

			*m.suggestionIdx++
			m.textInput.SetValue(suggestions[*m.suggestionIdx])
			m.textInput.CursorEnd()

			m.tip = ""
			return m, nil

		case tea.KeyCtrlL:
			if m.suggestionIdx == nil {
				return m, nil
			}

			suggestions, err := getSuggestions(m.textInput.Value())
			if err != nil {
				m.tip = err.Error()
				return m, nil
			}
			m.textInput.SetSuggestions(suggestions)
			m.suggestionIdx = nil

			return m, nil

		default:
			m.textInput, cmd = m.textInput.Update(msg)
			suggestions, err := getSuggestions(m.textInput.Value())
			if err != nil {
				m.tip = err.Error()
				return m, nil
			}

			if !slices.Contains(suggestions, m.textInput.Value()) {
				m.textInput.SetSuggestions(suggestions)
				m.suggestionIdx = nil
			}

			if m.lastInput != m.textInput.Value() && m.suggestionIdx == nil {
				m.lastInput = m.textInput.Value()
				m.tip = ""
			}

			return m, cmd
		}

	// We handle errors just like any other message
	case error:
		m.tip = msg.Error()
		return m, nil
	}

	return m, cmd
}

func (m initialPage) View() string {
	availableSuggestions := m.textInput.AvailableSuggestions()

	suggestionsBlock := ""
	if len(availableSuggestions) > 0 {
		suggestionsString := ""
		suggestionPagination := ""

		if m.suggestionIdx != nil {
			suggestionPagination = "  " + strings.Repeat(".", *m.suggestionIdx) + "x" + strings.Repeat(".", len(availableSuggestions)-*m.suggestionIdx-1)
		} else {
			suggestionPagination = "  " + strings.Repeat(".", len(availableSuggestions))
		}
		suggestionsPerPage := 3
		suggestions := make([]string, suggestionsPerPage)
		var pageStartIdx int
		if m.suggestionIdx == nil {
			pageStartIdx = 0
		} else {
			pageStartIdx = *m.suggestionIdx / suggestionsPerPage * suggestionsPerPage
		}
		var pageEndIdx int = pageStartIdx + suggestionsPerPage
		if len(availableSuggestions) < pageEndIdx {
			pageEndIdx = len(availableSuggestions)
		}

		pageWithSuggestion := m.textInput.AvailableSuggestions()[pageStartIdx:pageEndIdx]

		for i, suggestion := range pageWithSuggestion {
			checkMark := " "
			if m.suggestionIdx != nil && *m.suggestionIdx == i+pageStartIdx {
				checkMark = ">"
			}
			pathRegexp := regexp.MustCompile(`^(?P<dir>.+/)?(?P<file>[^/]*)$`)
			pathMatches := pathRegexp.FindStringSubmatch(m.textInput.Value())
			dirIndex := pathRegexp.SubexpIndex("dir")
			dirPath := pathMatches[dirIndex]

			suggestionFile := suggestion
			if dirPath != "" {
				suggestionFile = strings.Replace(suggestionFile, dirPath, "", 1)
			}

			suggestionString := fmt.Sprintf("%s %s", checkMark, suggestionFile)
			suggestions[i] = suggestionString
		}

		suggestionsString = strings.Join(suggestions, "\n")
		suggestionsManual := `next/previous - ctrl+n/p, choose - ctrl+l`
		suggestionsBlock = fmt.Sprintf(
			"suggestions(%s):\n%s\n%s\n",
			suggestionsManual,
			suggestionsString,
			suggestionPagination,
		)

	}

	tip := ""
	if m.tip != "" {
		tip = fmt.Sprintf("tip: %s\n", m.tip)
	}

	return fmt.Sprintf(
		"Enter path to swagger:\n\n%s\n\n%s%s%s\n",
		m.textInput.View(),
		suggestionsBlock,
		tip,
		"(esc to quit)",
	)
}
