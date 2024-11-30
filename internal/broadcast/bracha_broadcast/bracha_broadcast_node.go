package brachabroadcast

import (
	"net"
)

type Node struct {
	NodeType string
	Id       string
	Conn     map[string]net.Conn
}

func NewNode(nodeType, id string) *Node {
	return &Node{
		NodeType: nodeType,
		Id:       id,
		Conn:     make(map[string]net.Conn),
	}
}
