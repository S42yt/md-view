package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const version = "0.1.0"

func main() {
	style := flag.String("s", "auto", "style: auto, dark, light, dracula, tokyo-night, pink, notty, ascii, or a glamour JSON file")
	width := flag.Int("w", 0, "wrap width (0 = terminal width, capped at 120)")
	showVersion := flag.Bool("v", false, "print version and exit")
	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Println("md-view " + version)
		return
	}

	var (
		src  []byte
		name string
		path string
		err  error
		opts []tea.ProgramOption
	)

	switch {
	case flag.NArg() >= 1 && flag.Arg(0) != "-":
		path = flag.Arg(0)
		name = filepath.Base(path)
		src, err = os.ReadFile(path)
	case stdinPiped():
		name = "stdin"
		src, err = io.ReadAll(os.Stdin)
		opts = append(opts, tea.WithInputTTY())
	default:
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "md-view: %v\n", err)
		os.Exit(1)
	}

	if *style == "auto" {
		if lipgloss.HasDarkBackground() {
			*style = "dark"
		} else {
			*style = "light"
		}
	}

	m := newModel(path, name, src, *style, *width)
	opts = append(opts, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := tea.NewProgram(m, opts...).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "md-view: %v\n", err)
		os.Exit(1)
	}
}

func stdinPiped() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}

func usage() {
	fmt.Fprintf(os.Stderr, "md-view %s - view markdown in the terminal, nano style\n\n", version)
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintln(os.Stderr, "  md-view [options] FILE")
	fmt.Fprintln(os.Stderr, "  cat FILE.md | md-view [options]")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}
