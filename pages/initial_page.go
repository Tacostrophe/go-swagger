package pages

import (
	"fmt"
	"os"
	"regexp"
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
	ti.ShowSuggestions = true
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
			if err := m.usecase.Init(m.textInput.Value()); err != nil {
				m.tip = fmt.Sprintf("error: can't parse file: %s", err.Error())
				m.textInput.Reset()

				return m, nil
			}

			return NewSwaggerPage(m.usecase), nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlN:
			suggestions := m.textInput.AvailableSuggestions()
			if m.suggestionIdx == nil {
				m.suggestionIdx = new(int)
				*m.suggestionIdx = 0
				m.textInput.SetValue(suggestions[*m.suggestionIdx])
				m.textInput.CursorEnd()

				return m, nil
			}

			if *m.suggestionIdx == len(suggestions)-1 {
				m.suggestionIdx = nil
				m.textInput.SetValue(m.lastInput)
				m.textInput.CursorEnd()

				return m, nil
			}

			*m.suggestionIdx++
			m.textInput.SetValue(suggestions[*m.suggestionIdx])
			m.textInput.CursorEnd()

			return m, nil
		// case tea.KeyRunes:
		default:
			m.textInput, cmd = m.textInput.Update(msg)
			suggestions, err := getSuggestions(m.textInput.Value())
			if err != nil {
				m.tip = err.Error()
				return m, nil
			}
			m.textInput.SetSuggestions(suggestions)
			m.suggestionIdx = nil
			m.lastInput = m.textInput.Value()

			return m, cmd
		}

	// We handle errors just like any other message
	case error:
		// m.err = msg
		return m, nil
	}

	return m, cmd
}

func (m initialPage) View() string {
	suggestions := make([]string, len(m.textInput.AvailableSuggestions()))
	for i, suggestion := range m.textInput.AvailableSuggestions() {
		checkMark := "[ ]"
		if m.suggestionIdx != nil && *m.suggestionIdx == i {
			checkMark = "[X]"
		}
		suggestionString := fmt.Sprintf("%s %s", checkMark, suggestion)
		suggestions[i] = suggestionString
	}

	return fmt.Sprintf(
		"Enter path to swagger\n\n%s\n\n%s\n%s\n%s\n",
		m.textInput.View(),
		strings.Join(suggestions, ", "),
		m.tip,
		"(esc to quit)",
	)
}
