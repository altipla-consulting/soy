package parsepasses

import (
	"github.com/altipla-consulting/soy/ast"
	"github.com/altipla-consulting/soy/soymsg"
	"github.com/altipla-consulting/soy/template"
)

// ProcessMessages calculates the message ids and placeholder names for {msg}
// nodes and sets that information on the node.
func ProcessMessages(reg template.Registry) {
	for _, t := range reg.Templates {
		processTemplateMsgs(t.Node)
	}
}

func processTemplateMsgs(node ast.Node) {
	switch node := node.(type) {
	case *ast.MsgNode:
		soymsg.SetPlaceholdersAndID(node)
	default:
		if parent, ok := node.(ast.ParentNode); ok {
			for _, child := range parent.Children() {
				processTemplateMsgs(child)
			}
		}
	}
}
