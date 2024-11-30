package brachabroadcast

import (
	"DAG_Reliable_Broadcast/internal/config"
	"log"
	"net"
)

const (
	configPath = "config/host_config.json"
)

var instance *NodeExtention

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

type NodeExtention struct {
	// 节点信息
	Node Node

	// 缓存池
	InitialPool chan InitialMessage
	EchoPool    chan EchoMessage
	ReadyPool   chan ReadyMessage

	// 拜占庭阈值
	T int
	N int

	// 用于记录发送过的消息类型的map
	HadEchoInitial  map[string]int // 组成{Gethash(str(initial)): 1}
	HadReadyInitial map[string]int // 组成{Gethash(str(initial)): 1}

}

func NewNodeExtentions(node Node) *NodeExtention {
	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Printf("[ERROR] 加载配置失败: %v", err)
		return nil
	}
	T := config.T
	N := config.N
	log.Printf("[INFO] 加载配置成功: T=%d, N=%d", T, N) // 添加日志记录 T 和 N 的值
	return &NodeExtention{
		Node:            node,
		InitialPool:     make(chan InitialMessage, 100000), // 设置缓冲区大小
		EchoPool:        make(chan EchoMessage, 100000),    // 设置缓冲区大小
		ReadyPool:       make(chan ReadyMessage, 100000),   // 设置缓冲区大小
		T:               T,
		N:               N,
		HadEchoInitial:  make(map[string]int),
		HadReadyInitial: make(map[string]int),
	}
}
