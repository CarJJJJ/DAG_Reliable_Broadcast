package dagbroadcast

import (
	"net"
)

type Node struct {
	NodeType string
	Id       int
	Conn     map[string]net.Conn
}

func NewNode(nodeType string, id int) *Node {
	return &Node{
		NodeType: nodeType,
		Id:       id,
		Conn:     make(map[string]net.Conn),
	}
}
