package markdown

import (
	"io/ioutil"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

func RenderFile(file string, rootDir string) (string, error) {
	// extensions := 0 |
	// 	bf.NoIntraEmphasis |
	// 	bf.FencedCode |
	// 	bf.Autolink |
	// 	bf.Strikethrough |
	// 	bf.SpaceHeadings |
	// 	bf.HeadingIDs |
	// 	bf.BackslashLineBreak |
	// 	bf.DefinitionLists

	ast, err := ParseFile(file, rootDir)
	if err != nil {
		return "", err
	}
	output := md.Render(ast, NewConsoleRenderer(rootDir))
	return string(output), nil
}

func ParseFile(file string, rootDir string) (ast.Node, error) {
	// extensions := 0 |
	// 	bf.NoIntraEmphasis |
	// 	bf.FencedCode |
	// 	bf.Autolink |
	// 	bf.Strikethrough |
	// 	bf.SpaceHeadings |
	// 	bf.HeadingIDs |
	// 	bf.BackslashLineBreak |
	// 	bf.DefinitionLists

	input, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	node := md.Parse(input, parser.New())
	return node, nil
}
