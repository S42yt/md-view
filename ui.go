package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

var (
	barStyle    = lipgloss.NewStyle().Reverse(true).Bold(true)
	keyStyle    = lipgloss.NewStyle().Reverse(true)
	statusStyle = lipgloss.NewStyle().Reverse(true)
)

type shortcut struct {
	key   string
	label string
}

var (
	viewRow1   = []shortcut{{"^G", "Help"}, {"^W", "Where Is"}, {"^N", "Next"}, {"^Y", "Prev Page"}, {"^T", "Go To Line"}, {"^S", "Style"}}
	viewRow2   = []shortcut{{"^X", "Exit"}, {"^R", "Raw"}, {"^P", "Prev"}, {"^V", "Next Page"}, {"^L", "Reload"}, {"^C", "Position"}}
	helpRow1   = []shortcut{{"^X", "Close"}, {"^Y", "Prev Page"}}
	helpRow2   = []shortcut{{"^G", "Close"}, {"^V", "Next Page"}}
	promptRow1 = []shortcut{{"^C", "Cancel"}}
	promptRow2 = []shortcut{{"Enter", "Confirm"}}
)

func (m model) View() string {
	if !m.ready {
		return ""
	}
	return strings.Join([]string{m.titleBar(), m.vp.View(), m.statusLine(), m.shortcutBar()}, "\n")
}

func (m model) titleBar() string {
	w := m.width
	left := "  md-view " + version
	name := m.name
	if m.raw {
		name += " [raw]"
	}
	if m.mode == modeHelp {
		name = "Help"
	}
	right := fmt.Sprintf("%s  %d/%d  %3.0f%%  ", m.style, min(m.vp.YOffset+1, len(m.lines)), len(m.lines), m.vp.ScrollPercent()*100)

	lpad := w/2 - lipgloss.Width(name)/2 - lipgloss.Width(left)
	if lpad < 2 {
		lpad = 2
	}
	line := left + strings.Repeat(" ", lpad) + name
	rpad := w - lipgloss.Width(line) - lipgloss.Width(right)
	if rpad >= 2 {
		line += strings.Repeat(" ", rpad) + right
	}
	line = ansi.Truncate(line, w, "")
	if pad := w - lipgloss.Width(line); pad > 0 {
		line += strings.Repeat(" ", pad)
	}
	return barStyle.Render(line)
}

func (m model) statusLine() string {
	w := m.width
	switch m.mode {
	case modeSearch, modeGoto:
		return ansi.Truncate(m.input.View(), w, "")
	}
	if m.status == "" {
		return ""
	}
	msg := ansi.Truncate("[ "+m.status+" ]", w, "")
	pad := max(0, (w-lipgloss.Width(msg))/2)
	return strings.Repeat(" ", pad) + statusStyle.Render(msg)
}

func (m model) shortcutBar() string {
	row1, row2 := viewRow1, viewRow2
	switch m.mode {
	case modeHelp:
		row1, row2 = helpRow1, helpRow2
	case modeSearch, modeGoto:
		row1, row2 = promptRow1, promptRow2
	}
	cols := max(1, min(6, m.width/13))
	cellW := m.width / cols
	cell := lipgloss.NewStyle().Width(cellW).MaxWidth(cellW)
	render := func(row []shortcut) string {
		var b strings.Builder
		for i := 0; i < cols && i < len(row); i++ {
			b.WriteString(cell.Render(keyStyle.Render(row[i].key) + " " + row[i].label))
		}
		return b.String()
	}
	return render(row1) + "\n" + render(row2)
}
