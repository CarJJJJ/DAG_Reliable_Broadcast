package dagctbroadcast

import (
	"DAG_Reliable_Broadcast/internal/broadcast/util"
	"DAG_Reliable_Broadcast/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
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
	Type            int         `json:"type"`
	SendMessage     SendMessage `json:"send_message"`
	EchoReferences  []string    `json:"echo_references"`
	ReadyReferences []string    `json:"ready_references"`
	NodeID          int         `json:"node_id"`
	UniqueIndex     string      `json:"unique_index"`
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
	IDS_ECHO   *ThreadSafeSet // 论文中表示为ids_e, 存储节点号, 但只来自于ref1和ref2
	IDS_READY  *ThreadSafeSet // 论文中表示为ids_r, 存储节点号, 但只来自于ref3
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

	// 论文中表示为G_i, 组成{ResponseHash}, 表示为图的节点
	V *ThreadSafeSet

	// 论文中表示为E_i, 组成{ResponseHash}, 表示为图的边
	E *ThreadSafeSet

	// Robin-Round Slice
	RobinRound   []int
	CurrentIndex int // 记录当前轮询的起始位置

	// had echo or ready for uniqueIndex
	HadEchoUniqueIndex  sync.Map // 组成{UniqueIndex: 1}, 论文当中表示为Echo_i
	HadReadyUniqueIndex sync.Map // 组成{UniqueIndex: 1}, 论文当中表示为Ready_i

	// 论文中表示为votes_i
	Vote sync.Map // 组成{ResponseHash: VoteValue},表示这个节点被哪些节点引用，被哪些节点echo，被哪些节点ready

	// 论文中表示VK
	VK sync.Map // 组成{UniqueIndex: ResponseHash切片}, 表示这个UniqueIndex和哪些ResponseHash有关系

	InitialResponseHash string // 论文中的初始Response

	// 论文中没有提及但是需要的结构
	HashToResponse *sync.Map // 组成{ResponseHash: Response消息}

	// 维护一个论文里面没有提及的结构，用于选ECHO的
	// 在接受Response的时候，将Response取哈希后存入对应的NodeID的UniqueIndex的分组
	// 在可靠广播了某个UniqueIndex后, 遍历这个结构，删掉UniqueIndex下的ResponseHash
	NodeIDToUniqueIndexToResponseHashSlice *sync.Map // 组成{NodeID: &{UniqueIndex: &ResponseHash切片}}

	// 论文里面的inbuffer
	InBuffer *sync.Map // 组成{UniqueIndex: Send消息}
	Y        int

	// 可靠广播UniqueIndex对应的分片数量
	ReliableBroadcastUniqueIndexShardCount sync.Map // 组成{UniqueIndex: 分片数量}

	// 可靠广播量级
	ReliableBroadcastCount           int
	ReliableBroadcastCountLastSecond int

	// 可靠广播的UniqueIndex
	ReliableBroadcastUniqueIndex sync.Map // 组成{UniqueIndex: 1}, 表示这个UniqueIndex已经被可靠广播了

	// 可靠广播的UniqueIndex下的分片ID
	ReliableBroadcastUniqueIndexToShardIDSlice sync.Map // 组成{UniqueIndex: &{分片ID切片}}

	random *rand.Rand
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

	// 初始化RobinRound
	RobinRound := make([]int, N)
	for i := 0; i < N; i++ {
		RobinRound[i] = i
	}

	// 初始化Response
	initialResponse := ResponseMessage{
		Type: ResponseType,
		SendMessage: SendMessage{
			Type:        SendType,
			Message:     "initial",
			NodeID:      -1,
			UniqueIndex: "-1",
		},
		EchoReferences:  []string{},
		ReadyReferences: []string{},
		NodeID:          -1,
		UniqueIndex:     "-1",
	}
	initialResponseHash, err := util.CalculateHash(initialResponse)
	if err != nil {
		log.Printf("[ERROR] 计算初始Response哈希值失败: %v", err)
		return nil
	}

	writeResponseToJSON(node.Id, initialResponseHash, initialResponse)

	// V 初始化
	V := NewThreadSafeSet()
	V.Add(initialResponseHash)

	// HashToResponse
	hashToResponse := &sync.Map{}
	hashToResponse.Store(initialResponseHash, &initialResponse)

	// 初始化论文未提及的NodeIDToUniqueIndexToResponsehashes
	nodeIDToUniqueIndexToResponseHashSliceMap := &sync.Map{}
	for i := 0; i < N; i++ {
		// 初始化的时候，不分配内存
		uniqueIndexToResponseHashSliceMap := &sync.Map{}
		uniqueIndexToResponseHashSliceMap.Store("-1", &[]string{initialResponseHash})
		nodeIDToUniqueIndexToResponseHashSliceMap.Store(i, uniqueIndexToResponseHashSliceMap)
	}

	return &NodeExtention{
		Node:                                   node,
		Config:                                 *config,
		T:                                      T,
		N:                                      N,
		SendPool:                               make(chan SendMessage, 20000),     // 设置缓冲区大小
		ResponsePool:                           make(chan ResponseMessage, 20000), // 设置缓冲区大小
		SyncPool:                               make(chan SyncMessage, 20000),     // 设置缓冲区大小
		V:                                      V,
		E:                                      NewThreadSafeSet(),
		RobinRound:                             RobinRound,
		InitialResponseHash:                    initialResponseHash,
		NodeIDToUniqueIndexToResponseHashSlice: nodeIDToUniqueIndexToResponseHashSliceMap,
		HashToResponse:                         hashToResponse,
		random:                                 rand.New(rand.NewSource(time.Now().UnixNano() + int64(node.Id))),
		Y:                                      5,
		// SyncMap不需要显示初始化
		InBuffer: &sync.Map{}, //?
	}
}

// 定义一个全局互斥锁来保护文件写入
var jsonFileMutex sync.Mutex

func writeResponseToJSON(nodeId int, responseHash string, msg ResponseMessage) {
	if nodeId != 0 {
		// log.Printf("[INFO] 节点%d不写入初始Response", nodeId)
		return
	}

	// 获取文件锁
	jsonFileMutex.Lock()
	defer jsonFileMutex.Unlock()

	// 打开文件，需要读写权限
	file, err := os.OpenFile("DAG_Graph.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("[ERROR] 无法打开文件: %v", err)
		return
	}
	defer file.Close()

	// 创建新的节点数据
	newNode := struct {
		ResponseHash    string   `json:"ResponseHash"`
		EchoResponse    []string `json:"echoResponse"`
		ReadyResponse   []string `json:"readyResponse"`
		NodeDescription string   `json:"nodeDescription"`
	}{
		ResponseHash:    responseHash,
		EchoResponse:    msg.EchoReferences,
		ReadyResponse:   msg.ReadyReferences,
		NodeDescription: fmt.Sprintf("R(%d_%s)", msg.NodeID, msg.UniqueIndex),
	}

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("[ERROR] 获取文件信息失败: %v", err)
		return
	}

	if fileInfo.Size() == 0 {
		// 如果文件为空，写入数组开始
		if _, err := file.WriteString("[\n  "); err != nil {
			log.Printf("[ERROR] 写入失败: %v", err)
			return
		}
	} else {
		// 如果文件不为空，删除最后的 "\n]" 并添加逗号
		if _, err := file.Seek(-2, 2); err != nil {
			log.Printf("[ERROR] 文件定位失败: %v", err)
			return
		}
		if _, err := file.WriteString(",\n  "); err != nil {
			log.Printf("[ERROR] 写入失败: %v", err)
			return
		}
	}

	// 将新节点转换为JSON并写入
	jsonData, err := json.Marshal(newNode)
	if err != nil {
		log.Printf("[ERROR] JSON 编码失败: %v", err)
		return
	}

	// 写入新节点并关闭数组
	if _, err := file.WriteString(string(jsonData) + "\n]"); err != nil {
		log.Printf("[ERROR] 写入失败: %v", err)
		return
	}
}
