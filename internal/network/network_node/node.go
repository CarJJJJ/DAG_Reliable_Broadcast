package networknode

import (
	"net"
	"sync"
)

type Node struct {
	NodeType string
	Id       int
	Conn     map[string]net.Conn
	mu       sync.Mutex // 添加一个互斥锁
}

func NewNode(nodeType string, id int) *Node {
	return &Node{
		NodeType: nodeType,
		Id:       id,
		Conn:     make(map[string]net.Conn),
	}
}
