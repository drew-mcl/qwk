package pkg

import (
	"fmt"
)

type NodeType int

const (
	RootNode NodeType = iota
	VariableLine
	TargetLine
	CommandLine
)

func (nt NodeType) String() string {
	switch nt {
	case RootNode:
		return "RootNode"
	case VariableLine:
		return "VariableLine"
	case TargetLine:
		return "TargetLine"
	case CommandLine:
		return "CommandLine"
	default:
		return fmt.Sprintf("Unknown(%d)", nt)
	}
}

type Node struct {
	Type     NodeType
	Key      string
	Value    string
	Children []*Node
}
