/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/purpleclay/shake/internal/parser"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
)

type PromptsModel struct {
	focusIndex int
	Questions  []textinput.Model
}

func Prompts(cfg []parser.TextPrompt) *PromptsModel {
	// TODO: Style using lip gloss

	pm := PromptsModel{
		Questions: make([]textinput.Model, 0, len(cfg)),
	}
	for _, tp := range cfg {
		ti := textinput.New()
		ti.Placeholder = tp.Question
		ti.CharLimit = 32
		ti.CursorStyle = cursorStyle

		pm.Questions = append(pm.Questions, ti)
	}

	pm.Questions[0].Focus()
	pm.Questions[0].PromptStyle = focusedStyle
	pm.Questions[0].TextStyle = focusedStyle

	return &pm
}

func (m PromptsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PromptsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.Questions) {
				return m, tea.Quit
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.Questions) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.Questions)
			}

			cmds := make([]tea.Cmd, len(m.Questions))
			for i := 0; i <= len(m.Questions)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.Questions[i].Focus()
					m.Questions[i].PromptStyle = focusedStyle
					m.Questions[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.Questions[i].Blur()
				m.Questions[i].PromptStyle = noStyle
				m.Questions[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *PromptsModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Questions))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for q := range m.Questions {
		m.Questions[q], cmds[q] = m.Questions[q].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m PromptsModel) View() string {
	var b strings.Builder

	for i := range m.Questions {
		b.WriteString(m.Questions[i].View())
		if i < len(m.Questions)-1 {
			b.WriteRune('\n')
		}
	}

	// button := &blurredButton
	// if m.focusIndex == len(m.inputs) {
	// 	button = &focusedButton
	// }
	// fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	// b.WriteString(helpStyle.Render("cursor mode is "))
	// b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	// b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}
