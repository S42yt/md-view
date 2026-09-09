package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	modeView mode = iota
	modeSearch
	modeGoto
	modeHelp
)

var styleNames = []string{"github", "github-light", "dark", "light", "dracula", "tokyo-night", "pink", "notty", "ascii"}

type model struct {
	path        string
	name        string
	source      []byte
	style       string
	styleIdx    int
	userWidth   int
	width       int
	height      int
	ready       bool
	vp          viewport.Model
	input       textinput.Model
	mode        mode
	raw         bool
	rendered    string
	lines       []string
	query       string
	matches     []int
	matchIdx    int
	status      string
	savedOffset int
}

func newModel(path, name string, src []byte, style string, width int) model {
	ti := textinput.New()
	ti.CharLimit = 256
	m := model{
		path:      path,
		name:      name,
		source:    src,
		style:     style,
		styleIdx:  -1,
		userWidth: width,
		input:     ti,
	}
	for i, s := range styleNames {
		if s == style {
			m.styleIdx = i
		}
	}
	m.status = fmt.Sprintf("Read %d lines", sourceLines(src))
	return m
}

func sourceLines(src []byte) int {
	if len(src) == 0 {
		return 0
	}
	n := bytes.Count(src, []byte("\n"))
	if src[len(src)-1] != '\n' {
		n++
	}
	return n
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		h := max(1, msg.Height-4)
		if !m.ready {
			m.vp = viewport.New(msg.Width, h)
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = h
		}
		m.render()
		if m.mode == modeHelp {
			m.vp.SetContent(m.renderMarkdown(helpText))
		}
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()
	switch m.mode {
	case modeSearch, modeGoto:
		return m.handlePromptKey(msg)
	case modeHelp:
		switch key {
		case "ctrl+x", "ctrl+g", "q", "esc":
			m.closeHelp()
			return m, nil
		}
		return m.scrollKey(msg)
	}

	m.status = ""
	switch key {
	case "ctrl+x", "q", "esc":
		return m, tea.Quit
	case "ctrl+g", "?":
		m.openHelp()
		return m, nil
	case "ctrl+w", "/":
		return m.openPrompt(modeSearch, "Search: ", m.query)
	case "ctrl+t":
		return m.openPrompt(modeGoto, "Enter line number: ", "")
	case "ctrl+n", "n":
		m.step(1)
		return m, nil
	case "ctrl+p", "N":
		m.step(-1)
		return m, nil
	case "ctrl+r":
		m.raw = !m.raw
		m.render()
		if m.raw {
			m.status = "Showing raw source"
		} else {
			m.status = "Showing rendered markdown"
		}
		return m, nil
	case "ctrl+s":
		if m.raw {
			m.status = "Styles apply to the rendered view only"
			return m, nil
		}
		m.styleIdx = (m.styleIdx + 1) % len(styleNames)
		m.style = styleNames[m.styleIdx]
		m.render()
		m.status = "Style: " + m.style
		return m, nil
	case "ctrl+l":
		m.reload()
		return m, nil
	case "ctrl+c":
		m.status = m.position()
		return m, nil
	}
	return m.scrollKey(msg)
}

func (m model) scrollKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+y":
		m.vp.PageUp()
		return m, nil
	case "ctrl+v":
		m.vp.PageDown()
		return m, nil
	case "g", "home", "ctrl+home":
		m.vp.GotoTop()
		return m, nil
	case "G", "end", "ctrl+end":
		m.vp.GotoBottom()
		return m, nil
	}
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m model) openPrompt(md mode, prompt, value string) (tea.Model, tea.Cmd) {
	m.mode = md
	m.input.Prompt = prompt
	m.input.SetValue(value)
	m.input.CursorEnd()
	m.input.Width = max(1, m.width-len(prompt)-1)
	return m, m.input.Focus()
}

func (m model) handlePromptKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c", "ctrl+x":
		m.mode = modeView
		m.input.Blur()
		m.status = "Cancelled"
		return m, nil
	case "enter":
		value := strings.TrimSpace(m.input.Value())
		md := m.mode
		m.mode = modeView
		m.input.Blur()
		m.status = ""
		switch md {
		case modeSearch:
			if value != "" {
				m.search(value)
			}
		case modeGoto:
			n, err := strconv.Atoi(value)
			if err != nil || n < 1 {
				m.status = "Invalid line number"
			} else {
				m.vp.SetYOffset(n - 1)
			}
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *model) search(q string) {
	m.query = q
	m.matches = m.matches[:0]
	lq := strings.ToLower(q)
	for i, l := range m.lines {
		if strings.Contains(strings.ToLower(l), lq) {
			m.matches = append(m.matches, i)
		}
	}
	if len(m.matches) == 0 {
		m.status = fmt.Sprintf("\"%s\" not found", q)
		return
	}
	m.matchIdx = 0
	for i, ln := range m.matches {
		if ln > m.vp.YOffset {
			m.matchIdx = i
			break
		}
	}
	m.showMatch()
}

func (m *model) step(d int) {
	if m.query == "" {
		m.status = "No search query, press ^W first"
		return
	}
	if len(m.matches) == 0 {
		m.status = fmt.Sprintf("\"%s\" not found", m.query)
		return
	}
	next := m.matchIdx + d
	wrapped := next < 0 || next >= len(m.matches)
	m.matchIdx = (next + len(m.matches)) % len(m.matches)
	m.showMatch()
	if wrapped {
		m.status = "Search wrapped, " + m.status
	}
}

func (m *model) showMatch() {
	ln := m.matches[m.matchIdx]
	m.vp.SetYOffset(ln)
	m.status = fmt.Sprintf("match %d of %d at line %d", m.matchIdx+1, len(m.matches), ln+1)
}

func (m *model) reload() {
	if m.path == "" {
		m.status = "Cannot reload stdin"
		return
	}
	src, err := os.ReadFile(m.path)
	if err != nil {
		m.status = err.Error()
		return
	}
	m.source = src
	m.render()
	m.status = fmt.Sprintf("Reloaded, %d lines", sourceLines(src))
}

func (m *model) openHelp() {
	m.mode = modeHelp
	m.savedOffset = m.vp.YOffset
	m.vp.SetContent(m.renderMarkdown(helpText))
	m.vp.GotoTop()
}

func (m *model) closeHelp() {
	m.mode = modeView
	m.vp.SetContent(m.rendered)
	m.vp.SetYOffset(m.savedOffset)
}

func (m model) position() string {
	total := len(m.lines)
	top := min(m.vp.YOffset+1, total)
	pct := 0
	if total > 0 {
		pct = top * 100 / total
	}
	return fmt.Sprintf("line %d/%d (%d%%), %d source lines, %d bytes", top, total, pct, sourceLines(m.source), len(m.source))
}
