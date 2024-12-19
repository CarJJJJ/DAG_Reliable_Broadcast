package dagbroadcast

import (
	"DAG_Reliable_Broadcast/internal/config"
	"log"
	"net"
	"sync"
)

const (
	ConfigPath = "config/host_config.json"
)

const (
	SendType     = 0
	ResponseType = 1
	SyncType     = 2
)

var Instance *NodeExtention

type SendMessage struct {
	Type        int    `json:"type"`
	Message     string `json:"message"`
	NodeID      int    `json:"node_id"`
	UniqueIndex string `json:"unique_index"`
}

type ResponseMessage struct {
	Type        int         `json:"type"`
	SendMessage SendMessage `json:"send_message"`
	References  []string    `json:"references"` // 图的边, 为Response的哈希值
	NodeID      int         `json:"node_id"`
	UniqueIndex string      `json:"unique_index"`
}

type SyncMessage struct {
	Type   int               `json:"type"`
	Buffer []ResponseMessage `json:"buffer"`
	NodeID int               `json:"node_id"`
}

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

// 论文中表示为votes_i中的value
type VoteValue struct {
	References *ThreadSafeSet // 论文中表示为refs, 存储ResponseHash
	IDS_ECHO   *ThreadSafeSet // 论文中表示为ids_e, 存储ReponseHash, 但只来自于ref1和ref2
	IDS_READY  *ThreadSafeSet // 论文中表示为ids_r, 存储UniqueIndex, 但只来自于ref3
}

type NodeExtention struct {
	// 节点信息
	Node Node

	// 拜占庭容错参数
	T int
	N int

	// 配置
	Config config.Config

	// 缓存池
	SendPool     chan SendMessage
	ResponsePool chan ResponseMessage
	SyncPool     chan SyncMessage

	// 论文中表示为G_i, 组成{MessageHash}, 表示为图的节点
	V *ThreadSafeSet

	// 论文中表示为E_i, 组成{ResponseHash}, 表示为图的边
	E *ThreadSafeSet

	// Robin-Round Slice
	RobinRound []int

	// had echo or ready for uniqueIndex
	HadEchoUniqueIndex  sync.Map // 组成{UniqueIndex: 1}, 论文当中表示为Echo_i
	HadReadyUniqueIndex sync.Map // 组成{UniqueIndex: 1}, 论文当中表示为Ready_i

	// 论文中表示为votes_i
	Vote sync.Map // 组成{ResponseHash: VoteValue},表示这个节点被哪些节点引用，被哪些节点echo，被哪些节点ready

	// 论文中表示VK
	VK sync.Map // 组成{UniqueIndex: ResponseHash}, 表示这个UniqueIndex和哪些ResponseHash有关系
}

func NewNodeExtentions(node Node) *NodeExtention {
	config, err := config.LoadConfig(ConfigPath)
	if err != nil {
		log.Printf("[ERROR] 加载配置失败: %v", err)
		return nil
	}
	T := config.T
	N := config.N
	log.Printf("[INFO] 加载配置成功: T=%d, N=%d", T, N) // 添加日志记录 T 和 N 的值
	return &NodeExtention{
		Node:         node,
		Config:       *config,
		T:            T,
		N:            N,
		SendPool:     make(chan SendMessage, 999999),     // 设置缓冲区大小
		ResponsePool: make(chan ResponseMessage, 999999), // 设置缓冲区大小
		SyncPool:     make(chan SyncMessage, 999999),     // 设置缓冲区大小
		V:            NewThreadSafeSet(),
		E:            NewThreadSafeSet(),
		RobinRound:   make([]int, N),
		// SyncMap不需要显示初始化
	}
}
