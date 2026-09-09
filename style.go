package main

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
)

func styleOption(name string) glamour.TermRendererOption {
	switch name {
	case "github":
		return glamour.WithStyles(githubStyle(styles.DarkStyleConfig, "75", "245", "245"))
	case "github-light":
		return glamour.WithStyles(githubStyle(styles.LightStyleConfig, "27", "243", "243"))
	}
	return glamour.WithStylePath(name)
}

func githubStyle(base ansi.StyleConfig, link, url, quote string) ansi.StyleConfig {
	t := true
	s := base
	s.Heading.Color = nil
	s.H1 = ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Bold: &t, Underline: &t}}
	s.H2 = ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Bold: &t, Underline: &t}}
	s.H3 = ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Bold: &t}}
	s.H4 = ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Bold: &t}}
	s.H5 = ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Bold: &t}}
	s.H6 = ansi.StyleBlock{StylePrimitive: ansi.StylePrimitive{Bold: &t, Faint: &t}}
	s.LinkText = ansi.StylePrimitive{Color: &link, Underline: &t}
	s.Link = ansi.StylePrimitive{Color: &url, Faint: &t}
	s.Code.Color = nil
	s.BlockQuote.Color = &quote
	s.HorizontalRule.Format = "\n────────────────────────────────────────\n"
	return s
}
