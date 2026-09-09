# md-view

A read-only markdown viewer for the terminal that feels like [nano](https://nano-editor.org).
Rendering is done by [glamour](https://github.com/charmbracelet/glamour), the engine behind
[glow](https://github.com/charmbracelet/glow). Nothing you press will ever change the file.

## Install

```sh
go install github.com/S42yt/md-view@latest
```

## Usage

```sh
md-view README.md
md-view -s dracula README.md
curl -s https://raw.githubusercontent.com/charmbracelet/glow/master/README.md | md-view
```

| Flag | Description |
|------|-------------|
| `-s` | style: `auto`, `dark`, `light`, `dracula`, `tokyo-night`, `pink`, `notty`, `ascii`, or a glamour JSON file |
| `-w` | wrap width, defaults to the terminal width capped at 120 |
| `-v` | print version |

## Keys

| Key | Alternative | Action |
|-----|-------------|--------|
| `^G` | `?` | Help |
| `^X` | `q`, `Esc` | Exit |
| `^W` | `/` | Search |
| `^N` / `^P` | `n` / `N` | Next / previous match |
| `^Y` / `^V` | `PgUp` / `PgDn` | Page up / down |
| `^T` | | Go to line |
| `^S` | | Cycle style |
| `^R` | | Toggle raw source |
| `^L` | | Reload file |
| `^C` | | Show position |
| `Home` / `End` | `g` / `G` | Top / bottom |

The mouse wheel scrolls as well.

## License

MIT
