package calm

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func PrintTree(out io.Writer, format string, a *Ast) {
	switch format {
	case "json":
		Walk(&treePrintVisitor{out: out, printer: &treeJsonPrint{}}, a)
	case "json-inline":
		Walk(&treePrintVisitor{out: out, printer: &treeJsonInlinePrint{}}, a)
	case "short":
		Walk(&treePrintVisitor{out: out, printer: &treeShortPrint{}}, a)
	case "short-inline":
		Walk(&treePrintVisitor{out: out, printer: &treeShortInlinePrint{}}, a)
	case "nice":
		Walk(&treePrintVisitor{out: out, printer: &treeNicePrint{}}, a)
	default:
		panic("unknown format " + format)
	}
}

type treePrintVisitor struct {
	out     io.Writer
	dep     int
	printer treePrinter
}

func (v *treePrintVisitor) Visit(n *Ast) Visitor {
	pad := strings.Repeat("    ", v.dep)
	v.printer.WriteOpen(v.out, pad, n)
	if len(n.Args) > 0 {
		v.printer.WriteArgsOpen(v.out, pad)
		v.dep++
		for i, a := range n.Args {
			if i > 0 {
				v.printer.WriteArgsSep(v.out)
			}
			Walk(v, a)
		}
		v.dep--
		v.printer.WriteArgsClose(v.out, pad)
	}
	v.printer.WriteClose(v.out, pad)
	return nil
}

type treePrinter interface {
	WriteOpen(out io.Writer, pad string, n *Ast)
	WriteClose(out io.Writer, pad string)
	WriteArgsOpen(out io.Writer, pad string)
	WriteArgsClose(out io.Writer, pad string)
	WriteArgsSep(out io.Writer)
}

type treeShortPrint struct{}

func (treeShortPrint) WriteOpen(out io.Writer, pad string, n *Ast) {
	fmt.Fprint(out, pad)
	fmt.Fprintf(out, "%s", n.Type)
	if n.Name.Text != "" {
		fmt.Fprintf(out, " %s", n.Name.Text)
	}
}

func (treeShortPrint) WriteArgsOpen(out io.Writer, pad string) {
	fmt.Fprint(out, " [")
	fmt.Fprint(out, "\n")
}

func (treeShortPrint) WriteArgsClose(out io.Writer, pad string) {
	fmt.Fprint(out, pad)
	fmt.Fprint(out, "]")
}

func (treeShortPrint) WriteClose(out io.Writer, pad string) {
	fmt.Fprint(out, "\n")
}

func (treeShortPrint) WriteArgsSep(out io.Writer) {
	fmt.Fprint(out, "")
}

type treeShortInlinePrint struct{}

func (treeShortInlinePrint) WriteOpen(out io.Writer, pad string, n *Ast) {
	fmt.Fprintf(out, "%s", n.Type)
	if n.Name.Text != "" {
		fmt.Fprintf(out, " %s", n.Name.Text)
	}
}

func (treeShortInlinePrint) WriteArgsOpen(out io.Writer, pad string) {
	fmt.Fprint(out, " [ ")
}

func (treeShortInlinePrint) WriteArgsClose(out io.Writer, pad string) {
	fmt.Fprint(out, " ]")
}

func (treeShortInlinePrint) WriteClose(out io.Writer, pad string) {
	fmt.Fprint(out, "")
}

func (treeShortInlinePrint) WriteArgsSep(out io.Writer) {
	fmt.Fprint(out, ", ")
}

type treeNicePrint struct{ treeShortInlinePrint }

func (treeNicePrint) WriteOpen(out io.Writer, pad string, n *Ast) {
	if n.Name.Text != "" {
		fmt.Fprintf(out, "%s", n.Name.Text)
	} else {
		fmt.Fprintf(out, "%s", n.Type)
	}
}

type treeJsonPrint struct{}

func (treeJsonPrint) WriteOpen(out io.Writer, pad string, n *Ast) {
	fmt.Fprint(out, pad)
	fmt.Fprint(out, pad)
	fmt.Fprint(out, "{")
	fmt.Fprint(out, "\n")
	fmt.Fprint(out, pad)
	fmt.Fprint(out, pad)
	fmt.Fprintf(out, `    "type": "%s"`, n.Type)
	if n.Name.Text != "" {
		b, _ := json.Marshal(n.Name.Text)
		fmt.Fprint(out, ",")
		fmt.Fprint(out, "\n")
		fmt.Fprint(out, pad)
		fmt.Fprint(out, pad)
		fmt.Fprintf(out, `    "name": %s`, b)
	}
}

func (treeJsonPrint) WriteArgsOpen(out io.Writer, pad string) {
	fmt.Fprint(out, ",")
	fmt.Fprint(out, "\n")
	fmt.Fprint(out, pad)
	fmt.Fprint(out, pad)
	fmt.Fprint(out, `    "args": [`)
	fmt.Fprint(out, "\n")
}

func (treeJsonPrint) WriteArgsClose(out io.Writer, pad string) {
	fmt.Fprint(out, "\n")
	fmt.Fprint(out, pad)
	fmt.Fprint(out, pad)
	fmt.Fprint(out, "    ]")
}

func (treeJsonPrint) WriteClose(out io.Writer, pad string) {
	fmt.Fprint(out, "\n")
	fmt.Fprint(out, pad)
	fmt.Fprint(out, pad)
	fmt.Fprint(out, "}")
}

func (treeJsonPrint) WriteArgsSep(out io.Writer) {
	fmt.Fprint(out, ",")
	fmt.Fprint(out, "\n")
}

type treeJsonInlinePrint struct{}

func (treeJsonInlinePrint) WriteOpen(out io.Writer, pad string, n *Ast) {
	fmt.Fprint(out, "{")
	fmt.Fprintf(out, ` "type": "%s"`, n.Type)
	if n.Name.Text != "" {
		b, _ := json.Marshal(n.Name.Text)
		fmt.Fprintf(out, `, "name": %s`, b)
	}
}

func (treeJsonInlinePrint) WriteArgsOpen(out io.Writer, pad string) {
	fmt.Fprint(out, `, "args": [`)
}

func (treeJsonInlinePrint) WriteArgsClose(out io.Writer, pad string) {
	fmt.Fprint(out, "]")
}

func (treeJsonInlinePrint) WriteClose(out io.Writer, pad string) {
	fmt.Fprint(out, " }")
}

func (treeJsonInlinePrint) WriteArgsSep(out io.Writer) {
	fmt.Fprint(out, ", ")
}
