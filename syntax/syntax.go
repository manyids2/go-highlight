package syntax

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/manyids2/go-tree-sitter-with-markdown"
	"github.com/manyids2/go-tree-sitter-with-markdown/golang"
	"github.com/manyids2/go-tree-sitter-with-markdown/javascript"
	"github.com/manyids2/go-tree-sitter-with-markdown/markdown"
	tsquery "github.com/manyids2/go-tree-sitter-with-markdown/query"
)

type Syntax struct {
	Path    string
	Ext     string
	Parser  *sitter.Parser
	Lang    *sitter.Language
	Content *[]byte
	Tree    *sitter.Tree
	Root    *sitter.Node
}

func (s Syntax) String() string {
	line := ""
	WalkNamedChildren(s.Root, 0, func(n *sitter.Node, indent int) {
		line += fmt.Sprintf("%s %s ", strings.Repeat(" ", indent), n.Type())
		line += fmt.Sprintf("[%d, %d] - [%d, %d] (%d, %d)\n",
			n.StartPoint().Row, n.StartPoint().Column,
			n.EndPoint().Row, n.EndPoint().Column,
			n.StartByte(), n.EndByte())
	})
	return line
}

func LoadSyntax(path string) (*Syntax, error) {
	// Initialize parser
	parser := sitter.NewParser()

	// Set language
	ext := filepath.Ext(path)
	var lang *sitter.Language
	switch ext {
	case ".scm":
		lang = tsquery.GetLanguage()
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
	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, err
	}

	f := &Syntax{
		Path:    path,
		Ext:     ext,
		Parser:  parser,
		Lang:    lang,
		Content: &content,
		Tree:    tree,
		Root:    tree.RootNode(),
	}
	return f, nil
}

// WalkNamedChildren Walk tree-sitter node (root, then named children)
func WalkNamedChildren(n *sitter.Node, indent int, callback func(n *sitter.Node, indent int)) {
	callback(n, indent)
	for i := 0; i < int(n.NamedChildCount()); i++ {
		child := n.NamedChild(i)
		WalkNamedChildren(child, indent+2, callback)
	}
}

// LoadQueries Load queries from file
func LoadQueries(path string) ([]byte, error) {
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

	return content, nil
}

// Query Execute query on root node.
func (s *Syntax) Query(query []byte) {
	q, err := sitter.NewQuery([]byte(query), s.Lang)
	if err != nil {
		log.Fatalln("Invalid query:", query)
	}
	qc := sitter.NewQueryCursor()
	qc.Exec(q, s.Root)
	// Iterate over query results
	count := 0
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, *s.Content)
		for i, c := range m.Captures {
			fmt.Println(i,
				"\nc.Node.Content(*s.Content),", c.Node.Content(*s.Content),
				"\nc.Node.Type(),             ", c.Node.Type(),
				"\nc.Index,                   ", c.Index,
				"\nm.ID,                      ", m.ID,
				"\nm.PatternIndex)            ", m.PatternIndex,
				"\n", c.Node.String())

			// m.ID -> global index over matches

			count += 1
		}
	}
	fmt.Println("count", count)
}
