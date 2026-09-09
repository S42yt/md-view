package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var gutterStyle = lipgloss.NewStyle().Faint(true)

func (m model) wrapWidth() int {
	if m.userWidth > 0 {
		return m.userWidth
	}
	return min(m.width, 120)
}

func (m model) renderMarkdown(src string) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath(m.style),
		glamour.WithWordWrap(m.wrapWidth()),
		glamour.WithEmoji(),
	)
	if err != nil {
		return "  " + err.Error()
	}
	out, err := r.Render(src)
	if err != nil {
		return "  " + err.Error()
	}
	return strings.Trim(out, "\n")
}

func (m model) renderRaw() string {
	src := strings.TrimRight(string(m.source), "\n")
	src = strings.ReplaceAll(src, "\t", "    ")
	ls := strings.Split(src, "\n")
	pad := len(strconv.Itoa(len(ls)))
	var b strings.Builder
	for i, l := range ls {
		b.WriteString(gutterStyle.Render(fmt.Sprintf(" %*d │ ", pad, i+1)))
		b.WriteString(strings.TrimRight(l, "\r"))
		b.WriteByte('\n')
	}
	return strings.TrimRight(b.String(), "\n")
}

func (m *model) render() {
	if m.raw {
		m.rendered = m.renderRaw()
	} else {
		m.rendered = m.renderMarkdown(string(m.source))
	}
	m.lines = strings.Split(ansi.Strip(m.rendered), "\n")
	m.vp.SetContent(m.rendered)
	if m.query != "" {
		offset := m.vp.YOffset
		m.search(m.query)
		m.vp.SetYOffset(offset)
		m.status = ""
	}
}
