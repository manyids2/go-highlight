package highlights

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/muesli/termenv"
)

// Highlights is dictionary to hold styles
type Highlights map[string]termenv.Style

// LoadHighlights reads file by line and parse into to list of styles
func LoadHighlights(path string, useCterm bool) (Highlights, error) {
	// Make sure file exists
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Whether to use cterm or gui
	colorPrefix := "gui"
	if useCterm {
		colorPrefix = "cterm"
	}

	// Use just one memory for clear style
	clearStyle := termenv.Style{}

	// Keep track of links, use later
	links := map[string]string{}

	// Scan all the highlights
	highlights := Highlights{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse nvim formatting
		parts := strings.Split(line, "xxx")

		// Ignore empty lines
		if len(parts) != 2 {
			continue
		}

		// Get name first
		name := strings.Trim(parts[0], " ")
		var hi termenv.Style

		// Then attributes
		attrs := strings.Split(parts[1], " ")
		switch attrs[0] {
		case "cleared":
			hi = clearStyle
		case "links":
			links[name] = attrs[2]
		default:
			hi = termenv.Style{}
			for _, attr := range attrs[1:] {
				parts := strings.Split(attr, "=")
				if len(parts) != 2 {
					continue
				}
				switch parts[0] {
				case colorPrefix:
					for _, t := range strings.Split(parts[1], ",") {
						switch t {
						case "reverse":
							hi = hi.Reverse()
						case "bold":
							hi = hi.Bold()
						case "italic":
							hi = hi.Italic()
						case "underline":
							hi = hi.Underline()
						case "undercurl":
							hi = hi.Underline()
						case "strikethrough":
							hi = hi.CrossOut()
						case "blink":
							hi = hi.Blink()
						}
					}
				case colorPrefix + "fg":
					hi = hi.Foreground(termenv.ANSI.Color(parts[1]))
				case colorPrefix + "bg":
					hi = hi.Background(termenv.ANSI.Color(parts[1]))
				case colorPrefix + "sp":
					// Do nothing
				default:
				}
			}
		}

		// Put into map
		highlights[name] = hi
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Fix linked styles
	for name, src := range links {
		var ok bool
		hi, ok := highlights[src]
		if !ok {
			return nil, err
		}
		highlights[name] = hi
	}

	return highlights, nil
}

// String prints styles in a table
func (h Highlights) Print() {
	// p := termenv.ColorProfile()
	for name, hi := range h {
		fmt.Printf("%s\n", hi.Styled(name))
	}
}
