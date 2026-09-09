package main

const helpText = `# md-view

A read-only markdown viewer that feels like nano. Content is rendered with glamour, the engine behind glow. Nothing you press will ever modify the file.

## Keys

| Key | Alternative | Action |
|-----|-------------|--------|
| ^G | ? | Show or close this help |
| ^X | q, Esc | Exit md-view |
| ^W | / | Search (case-insensitive) |
| ^N | n | Jump to next match |
| ^P | N | Jump to previous match |
| ^Y | PgUp, b | Scroll up one page |
| ^V | PgDn, Space | Scroll down one page |
| ^U | u | Scroll up half a page |
| ^D | d | Scroll down half a page |
| ^T | | Go to a line number |
| ^S | | Cycle through styles |
| ^R | | Toggle raw source / rendered view |
| ^L | | Reload the file from disk |
| ^C | | Show current position |
| Up, Down | k, j | Scroll one line |
| Home, End | g, G | Jump to top or bottom |

The mouse wheel scrolls too.

## Styles

github, github-light, dark, light, dracula, tokyo-night, pink, notty, ascii

Pick one at start with ` + "`md-view -s dracula FILE`" + ` or pass the path to your own glamour JSON style file. Use ` + "`-w N`" + ` to force a wrap width.

## Reading from a pipe

` + "```sh" + `
curl -s https://raw.githubusercontent.com/charmbracelet/glow/master/README.md | md-view
` + "```" + `
`
