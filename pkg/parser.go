package pkg

import (
	"errors"

	"github.com/drew-mcl/logz"
)

func Parse(nodes []*Node) (*Node, error) {
	root := &Node{Type: RootNode, Children: []*Node{}}

	for len(nodes) > 0 {
		var node *Node
		node, nodes = nodes[0], nodes[1:]

		if node.Type == VariableLine {
			logz.Debug("Parser found VariableLine, adding node with value of:", node.Value, "to root")

			root.Children = append(root.Children, node)
			continue
		}

		if node.Type == TargetLine {
			target := node
			target.Children = []*Node{}

			logz.Trace("Parsing target node with value of:", node.Key)

			// Commands follow their targets
			for len(nodes) > 0 && nodes[0].Type == CommandLine {
				logz.Debug("Assigning command:", nodes[0].Value, "to target:", target.Key)

				target.Children = append(target.Children, nodes[0])
				nodes = nodes[1:]
			}

			root.Children = append(root.Children, target)
			continue
		}

		return nil, errors.New("Invalid node type in parse: " + node.Type.String())
	}

	return root, nil
}
