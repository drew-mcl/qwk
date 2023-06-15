package pkg

import (
	"fmt"
)

func printNode(node *Node, indent string) {
	switch node.Type {
	case VariableLine:
		fmt.Printf("%sVariable: %s = %s\n", indent, node.Key, node.Value)
	case TargetLine:
		fmt.Printf("%sTarget: %s\n", indent, node.Key)
		if node.Children != nil {
			for _, child := range node.Children {
				printNode(child, indent+"\t")
			}
		}
	case CommandLine:
		fmt.Printf("%sCommand: %s\n", indent, node.Value)
	}
}

func PrintTree(root *Node) {
	for _, child := range root.Children {
		printNode(child, "")
	}
}
