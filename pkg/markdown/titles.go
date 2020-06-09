package markdown

import (
	"strings"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

type Title struct {
	Name     string
	SubTitle []Title
	Location string
	Level    int
}

func GetTitles(file string, rootDir string) ([]Title, error) {
	node, err := ParseFile(file, rootDir)
	if err != nil {
		return nil, err
	}

	renderer := NewConsoleRenderer(rootDir)

	children := node.GetChildren()
	titles, _ := extractTitles(1, children, 0, file, renderer)
	return titles, nil
}

func extractTitles(level int, children []ast.Node, index int, file string, renderer md.Renderer) ([]Title, int) {
	var titles []Title
	for index <= len(children)-1 {
		child := children[index]
		childHeading, ok := child.(*ast.Heading)
		if ok {
			if childHeading.Level >= level {
				output := strings.Trim(string(md.Render(childHeading, renderer)), "\n\t ")
				//fmt.Printf("i%v %s (%v) (%v) C%v T%v\n", index, output, childHeading.Level, level, len(children), calcTitles(children))
				subTitles, newIndex := extractTitles(level+1, children, index+1, file, renderer)
				index = newIndex
				titles = append(titles, Title{
					Name:     output,
					SubTitle: subTitles,
					Location: file,
					Level:    childHeading.Level,
				})
				continue
			} else {
				return titles, index
			}
		}
		index++
	}
	return titles, index
}

func calcTitles(children []ast.Node) int {
	count := 0
	for _, child := range children {
		children = children[1:]
		_, ok := child.(*ast.Heading)
		if ok {
			count++
		}
	}
	return count
}
