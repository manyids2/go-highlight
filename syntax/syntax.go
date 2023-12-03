package syntax

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	sitter "github.com/manyids2/go-tree-sitter-with-markdown"
	"github.com/manyids2/go-tree-sitter-with-markdown/golang"
	"github.com/manyids2/go-tree-sitter-with-markdown/javascript"
	"github.com/manyids2/go-tree-sitter-with-markdown/markdown"
)

func LoadTree(path string) (*sitter.Tree, error) {
	// Initialize parser
	parser := sitter.NewParser()

	// Set language
	ext := filepath.Ext(path)
	var lang *sitter.Language
	switch ext {
	case ".go":
		lang = golang.GetLanguage()
	case ".js":
		lang = javascript.GetLanguage()
	case ".md", ".MD", ".markdown":
		lang = markdown.GetLanguage()
	default:
		err := errors.New(fmt.Sprintf("Unsupported language: %s\n", ext))
		return nil, err
	}
	parser.SetLanguage(lang)

	// Make sure file exists
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read to []byte
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Parse
	tree := parser.Parse(nil, content)
	return tree, nil
}
