package calm

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func PrintTree(out io.Writer, format string, a *AST) {
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
	pad     string
	printer treePrinter
}

func (v *treePrintVisitor) Visit(n *AST) Visitor {
	pad := strings.Repeat("    ", v.dep)
	v.pad = pad
	v.printer.WriteOpen(v, n)
	if len(n.Args) > 0 {
		v.printer.WriteArgsOpen(v)
		v.dep++
		for i, a := range n.Args {
			if i > 0 {
				v.printer.WriteArgsSep(v)
			}
			Walk(v, a)
		}
		v.dep--
		v.pad = pad
		v.printer.WriteArgsClose(v)
	}
	v.printer.WriteClose(v)
	return nil
}

type treePrinter interface {
	WriteOpen(*treePrintVisitor, *AST)
	WriteClose(*treePrintVisitor)
	WriteArgsOpen(*treePrintVisitor)
	WriteArgsClose(*treePrintVisitor)
	WriteArgsSep(*treePrintVisitor)
}

type treeShortPrint struct{}

func (treeShortPrint) WriteOpen(v *treePrintVisitor, n *AST) {
	fmt.Fprint(v.out, v.pad)
	fmt.Fprintf(v.out, "%s", n.Type)
	if n.Name.Text != "" {
		fmt.Fprintf(v.out, " %s", n.Name.Text)
	}
}

func (treeShortPrint) WriteArgsOpen(v *treePrintVisitor) {
	fmt.Fprint(v.out, " [")
	fmt.Fprint(v.out, "\n")
}

func (treeShortPrint) WriteArgsClose(v *treePrintVisitor) {
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, "]")
}

func (treeShortPrint) WriteClose(v *treePrintVisitor) {
	if v.dep > 0 {
		fmt.Fprint(v.out, "\n")
	}
}

func (treeShortPrint) WriteArgsSep(v *treePrintVisitor) {
	fmt.Fprint(v.out, "")
}

type treeShortInlinePrint struct{}

func (treeShortInlinePrint) WriteOpen(v *treePrintVisitor, n *AST) {
	fmt.Fprintf(v.out, "%s", n.Type)
	if n.Name.Text != "" {
		fmt.Fprintf(v.out, " %s", n.Name.Text)
	}
}

func (treeShortInlinePrint) WriteArgsOpen(v *treePrintVisitor) {
	fmt.Fprint(v.out, " [ ")
}

func (treeShortInlinePrint) WriteArgsClose(v *treePrintVisitor) {
	fmt.Fprint(v.out, " ]")
}

func (treeShortInlinePrint) WriteClose(v *treePrintVisitor) {
	fmt.Fprint(v.out, "")
}

func (treeShortInlinePrint) WriteArgsSep(v *treePrintVisitor) {
	fmt.Fprint(v.out, ", ")
}

type treeNicePrint struct{ treeShortInlinePrint }

func (treeNicePrint) WriteOpen(v *treePrintVisitor, n *AST) {
	if n.Name.Text != "" {
		fmt.Fprintf(v.out, "%s", n.Name.Text)
	} else {
		fmt.Fprintf(v.out, "%s", n.Type)
	}
}

type treeJsonPrint struct{}

func (treeJsonPrint) WriteOpen(v *treePrintVisitor, n *AST) {
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, "{")
	fmt.Fprint(v.out, "\n")
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, v.pad)
	fmt.Fprintf(v.out, `    "type": "%s"`, n.Type)
	if n.Name.Text != "" {
		b, _ := json.Marshal(n.Name.Text)
		fmt.Fprint(v.out, ",")
		fmt.Fprint(v.out, "\n")
		fmt.Fprint(v.out, v.pad)
		fmt.Fprint(v.out, v.pad)
		fmt.Fprintf(v.out, `    "name": %s`, b)
	}
}

func (treeJsonPrint) WriteArgsOpen(v *treePrintVisitor) {
	fmt.Fprint(v.out, ",")
	fmt.Fprint(v.out, "\n")
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, `    "args": [`)
	fmt.Fprint(v.out, "\n")
}

func (treeJsonPrint) WriteArgsClose(v *treePrintVisitor) {
	fmt.Fprint(v.out, "\n")
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, "    ]")
}

func (treeJsonPrint) WriteClose(v *treePrintVisitor) {
	fmt.Fprint(v.out, "\n")
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, v.pad)
	fmt.Fprint(v.out, "}")
}

func (treeJsonPrint) WriteArgsSep(v *treePrintVisitor) {
	fmt.Fprint(v.out, ",")
	fmt.Fprint(v.out, "\n")
}

type treeJsonInlinePrint struct{}

func (treeJsonInlinePrint) WriteOpen(v *treePrintVisitor, n *AST) {
	fmt.Fprint(v.out, "{")
	fmt.Fprintf(v.out, ` "type": "%s"`, n.Type)
	if n.Name.Text != "" {
		b, _ := json.Marshal(n.Name.Text)
		fmt.Fprintf(v.out, `, "name": %s`, b)
	}
}

func (treeJsonInlinePrint) WriteArgsOpen(v *treePrintVisitor) {
	fmt.Fprint(v.out, `, "args": [`)
}

func (treeJsonInlinePrint) WriteArgsClose(v *treePrintVisitor) {
	fmt.Fprint(v.out, "]")
}

func (treeJsonInlinePrint) WriteClose(v *treePrintVisitor) {
	fmt.Fprint(v.out, " }")
}

func (treeJsonInlinePrint) WriteArgsSep(v *treePrintVisitor) {
	fmt.Fprint(v.out, ", ")
}
