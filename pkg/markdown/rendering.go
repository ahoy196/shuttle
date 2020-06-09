package markdown

import (
	"fmt"
	"html"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	md "github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/samfoo/ansi"
)

// Index corresponds to the heading level (e.g. h1, h2, h3...)
var headerStyles = [...]string{
	ansi.ColorCode("green+bhu"),
	ansi.ColorCode("green+bh"),
	ansi.ColorCode("green"),
	ansi.ColorCode("green"),
}

var emphasisStyles = [...]string{
	ansi.ColorCode("cyan+bh"),
	ansi.ColorCode("cyan+bhu"),
	ansi.ColorCode("cyan+bhi"),
}

var linkStyle = ansi.ColorCode("015+u")

const (
	UNORDERED = 1 << iota
	ORDERED
)

type list struct {
	kind  int
	index int
}

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

type Console struct {
	lists        []*list
	listLevel    int
	isBlockQuote bool
	columns      int
	rows         int
	rootDir      string
}

func NewConsoleRenderer(rootDir string) *Console {

	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return &Console{
		columns: int(ws.Col),
		rows:    int(ws.Row),
		rootDir: rootDir,
	}
}

func (r *Console) BlockCode(w io.Writer, text []byte, lang string) {

}

func (r *Console) BlockQuote(w io.Writer, text []byte) {

}

func (r *Console) BlockHtml(w io.Writer, text []byte) {
	w.Write(text)
}

func (r *Console) HRule(w io.Writer) {
	r.out(w, "\n\u2015\u2015\u2015\u2015\u2015\n\n")
}

func (r *Console) List(w io.Writer, text func() bool, flags int) {

}

func (r *Console) ListItem(w io.Writer, text []byte, flags int) {
	current := r.lists[len(r.lists)-1]

	for i := 0; i < len(r.lists); i++ {
		r.out(w, "  ")
	}

	if current.kind == ORDERED {
		r.out(w, fmt.Sprintf("%d. ", current.index))
		current.index += 1
	} else {
		r.out(w, ansi.ColorCode("red+bh"))
		r.out(w, "* ")
		r.out(w, ansi.ColorCode("reset"))
	}

	w.Write(text)
	r.out(w, "\n\n")
}

func (r *Console) Paragraph(w io.Writer, text func() bool) {
	r.out(w, "\n\n")
}

func (r *Console) Table(w io.Writer, header []byte, body []byte, columnData []int) {}
func (r *Console) TableRow(w io.Writer, text []byte)                               {}
func (r *Console) TableHeaderCell(w io.Writer, text []byte, flags int)             {}
func (r *Console) TableCell(w io.Writer, text []byte, flags int)                   {}
func (r *Console) Footnotes(w io.Writer, text func() bool)                         {}
func (r *Console) FootnoteItem(w io.Writer, name, text []byte, flags int)          {}

func (r *Console) TitleBlock(w io.Writer, text []byte) {
	r.out(w, "\n")
	r.out(w, headerStyles[0])
	w.Write(text)
	r.out(w, ansi.ColorCode("reset"))
	r.out(w, "\n\n")
}

func (r *Console) AutoLink(w io.Writer, link []byte, kind int) {
	r.out(w, linkStyle)
	w.Write(link)
	r.out(w, ansi.ColorCode("reset"))
}

func (r *Console) CodeSpan(w io.Writer, text []byte) {
	r.out(w, ansi.ColorCode("015+b"))
	w.Write(text)
	r.out(w, ansi.ColorCode("reset"))
}

func (r *Console) DoubleEmphasis(w io.Writer, text []byte) {
	r.out(w, emphasisStyles[1])
	w.Write(text)
	r.out(w, ansi.ColorCode("reset"))
}

func (r *Console) Emphasis(w io.Writer, text []byte) {
	r.out(w, emphasisStyles[0])
	w.Write(text)
	r.out(w, ansi.ColorCode("reset"))
}

func (r *Console) Image(w io.Writer, link []byte, title []byte, alt []byte) {
	r.out(w, " [ image ] ")
}

func (r *Console) LineBreak(w io.Writer) {
	r.out(w, "\n")
}

func (r *Console) Link(w io.Writer, link []byte, title []byte, content []byte) {
	w.Write(content)
	r.out(w, " (")
	r.out(w, linkStyle)
	w.Write(link)
	r.out(w, ansi.ColorCode("reset"))
	r.out(w, ")")
}

func (r *Console) RawHtmlTag(w io.Writer, tag []byte) {
	r.out(w, ansi.ColorCode("magenta"))
	w.Write(tag)
	r.out(w, ansi.ColorCode("reset"))
}

func (r *Console) TripleEmphasis(w io.Writer, text []byte) {
	r.out(w, emphasisStyles[2])
	w.Write(text)
	r.out(w, ansi.ColorCode("reset"))
}

func (r *Console) StrikeThrough(w io.Writer, text []byte) {
	r.out(w, ansi.ColorCode("008+s"))
	r.out(w, "\u2015")
	w.Write(text)
	r.out(w, "\u2015")
	r.out(w, ansi.ColorCode("reset"))
}

func (r *Console) FootnoteRef(w io.Writer, ref []byte, id int) {
}

func (r *Console) Entity(w io.Writer, entity []byte) {
	r.out(w, html.UnescapeString(string(entity)))
}

func (r *Console) NormalText(w io.Writer, text []byte) {
	s := string(text)
	reg, _ := regexp.Compile("\\s+")

	r.out(w, reg.ReplaceAllString(s, " "))
}

func (r *Console) out(w io.Writer, text string) {
	w.Write([]byte(text))
}

// RenderNode is the main rendering method. It will be called once for
// every leaf node and twice for every non-leaf node (first with
// entering=true, then with entering=false). The method should write its
// rendition of the node to the supplied writer w.
func (r *Console) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if entering && r.isBlockQuote {
		r.out(w, "\n  | ")
	}

	switch node := node.(type) {
	case *ast.Text:
		w.Write(node.Literal)
	case *ast.Softbreak:
		r.out(w, "\n")
	case *ast.Hardbreak:
		r.out(w, "\n")
	case *ast.Emph:
		r.out(w, emphasisStyles[0])
		w.Write(node.Literal)
		r.out(w, ansi.ColorCode("reset"))
	case *ast.Strong:
		r.out(w, ansi.ColorCode("blue"))
		w.Write(node.Literal)
		r.out(w, ansi.ColorCode("reset"))
	case *ast.Del:

	case *ast.HTMLSpan:

	case *ast.Link:

	case *ast.Image:
		if os.Getenv("ITERM_SESSION_ID") != "" {
			image, err := renderItermImage(string(node.Destination), r.rootDir)
			if err != nil {
				r.out(w, fmt.Sprintf("Error showing: ![%s](%s)\nError: %s\n", string(node.Title), string(node.Destination), err))
			}
			r.out(w, image)
		} else {
			r.out(w, fmt.Sprintf("![%s](%s)\n", string(node.Title), string(node.Destination)))
		}
	case *ast.Code:
		r.out(w, ansi.ColorCode("red"))
		w.Write(node.Literal)
		r.out(w, ansi.ColorCode("reset"))
	case *ast.Document:

	case *ast.Paragraph:
		if entering {
			r.out(w, "\n")
		} else {
			//r.out(w, "\n")
			r.out(w, "\n")
		}
	case *ast.BlockQuote:
		if entering {
			//r.out(w, "\n  | ")
			r.isBlockQuote = true
		} else {
			r.isBlockQuote = false
			r.out(w, "\n\n")
		}

	case *ast.HTMLBlock:

	case *ast.Heading:
		if entering {
			if node.Parent.GetChildren()[0] != node {
				r.out(w, "\n\n")
			}
			r.out(w, headerStyles[node.Level-1])
		} else {
			r.out(w, ansi.ColorCode("reset"))
			r.out(w, "\n\n")
		}
	case *ast.HorizontalRule:
		r.out(w, strings.Repeat("─", r.columns))
	case *ast.List:
		spaceCount := 2
		if _, ok := node.Parent.(*ast.ListItem); ok {
			spaceCount = 1
		}
		if entering {
			r.listLevel++
			r.out(w, strings.Repeat("\n", spaceCount))
		} else {
			r.listLevel--
			r.out(w, strings.Repeat("\n", spaceCount))
		}

	case *ast.ListItem:
		if entering {
			strings.Repeat(" ", r.listLevel*2)
			r.out(w, fmt.Sprintf("%s• ", strings.Repeat(" ", r.listLevel*2)))
		} else {

		}
	case *ast.CodeBlock:
		s := string(node.Literal)

		lines := strings.Split(s, "\n")

		r.out(w, "\n")
		r.out(w, "  ┌"+strings.Repeat("─", r.columns-5)+"┐\n")
		for i, line := range lines {
			if i == len(lines)-1 && strings.Trim(line, " ") == "" {
				break
			}
			//r.out(w, ansi.ColorCode("red"))
			if line == "" {
				r.out(w, fmt.Sprintf("  │ %-"+strconv.Itoa(r.columns-7)+"v │\n", line))
				continue
			}
			restline := line
			for len(restline) > 0 {
				curLine := restline
				if len(curLine) > r.columns-7 {
					curLine = restline[0:(r.columns - 8)]
					restline = restline[(r.columns - 8 + 1):]
					r.out(w, fmt.Sprintf("  │ %-"+strconv.Itoa(r.columns-8)+"v↵ │\n", curLine))
				} else {
					r.out(w, fmt.Sprintf("  │ %-"+strconv.Itoa(r.columns-7)+"v │\n", curLine))
					restline = ""
				}
			}

			//r.out(w, ansi.ColorCode("reset"))
		}
		r.out(w, "  └"+strings.Repeat("─", r.columns-5)+"┘\n")

		//r.out(w, reg.ReplaceAllString(s, "\n    "))
		//r.out(w, "\n")
	case *ast.Table:
		r.out(w, "TABLE")
	case *ast.TableCell:
		r.out(w, "TABLE CELL")
	case *ast.TableHeader:
		r.out(w, "TABLE HEADER")
	case *ast.TableBody:
		r.out(w, "TABLE BODY")
	case *ast.TableRow:
		r.out(w, "TABLE ROW")
	default:
		panic(fmt.Sprintf("Unknown node type %T", node))
	}
	return ast.GoToNext
}
func (r *Console) RenderHeader(w io.Writer, ast ast.Node) {
}
func (r *Console) RenderFooter(w io.Writer, ast ast.Node) {
}

var _ md.Renderer = &Console{}
