package networknode

import (
	"net"
)

type Node struct {
	nodeType string
	id       string
	conn     map[string]net.Conn
}

func NewNode(nodeType, id string) *Node {
	return &Node{
		nodeType: nodeType,
		id:       id,
		conn:     make(map[string]net.Conn),
	}
}
