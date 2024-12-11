package signbroadcast

import (
	"DAG_Reliable_Broadcast/internal/config"
	"fmt"
	"log"
	"net"
	"os"

	"math/rand"

	"github.com/CarJJJJ/go-bls"
	cmap "github.com/orcaman/concurrent-map/v2"
)

const (
	ConfigPath = "config/host_config.json"
)

const (
	BCBSendType        = 0
	BCBRepType         = 1
	BCBFinalType       = 2
	BCBDisperseType    = 3
	BCBReconstructType = 4
	BCBReadyType       = 5
)

type BCBSendMessage struct {
	NodeID      int    `json:"node_id"`
	Type        int    `json:"type"`
	Message     []byte `json:"message"`
	UniqueIndex string `json:"unique_index"`
}

type BCBRepMessage struct {
	NodeID      int    `json:"node_id"`
	Type        int    `json:"type"`
	Message     []byte `json:"message"`
	SigmaFrom   []byte `json:"sigma_from"`
	UniqueIndex string `json:"unique_index"`
}

type BCBFinalMessage struct {
	NodeID       int    `json:"node_id"`
	Type         int    `json:"type"`
	Message      []byte `json:"message"`
	SigmaCombine []byte `json:"sigma_combine"`
	UniqueIndex  string `json:"unique_index"`
}

type DisperseMessage struct {
	NodeID      int    `json:"node_id"`
	Type        int    `json:"type"`
	DataFrom    []byte `json:"data_from"`
	UniqueIndex string `json:"unique_index"`
}

type ReconstructMessage struct {
	NodeID      int    `json:"node_id"`
	Type        int    `json:"type"`
	DataFrom    []byte `json:"data_from"`
	UniqueIndex string `json:"unique_index"`
}

type ReadyMessage struct {
	NodeID      int    `json:"node_id"`
	Type        int    `json:"type"`
	Message     []byte `json:"message"`
	UniqueIndex string `json:"unique_index"`
}

var Instance *NodeExtention

type Node struct {
	NodeType string
	Id       int
	Conn     map[string]net.Conn
}

func NewNode(nodeType string, id int) *Node {
	return &Node{
		NodeType: nodeType,
		Id:       id,
	}
}

type NodeExtention struct {
	// 节点信息
	Node Node

	// 路由配置
	Config config.Config

	// 缓存池
	BCBSendPool        chan BCBSendMessage
	BCBRepPool         chan BCBRepMessage
	BCBFinalPool       chan BCBFinalMessage
	BCBDispersePool    chan DisperseMessage
	BCBReconstructPool chan ReconstructMessage

	// 公私密钥
	GroupKey      bls.PublicKey
	MemberKeys    []bls.PublicKey
	GroupSecret   bls.PrivateKey
	MemberSecrets []bls.PrivateKey

	// 用于记录发送过的消息类型的map
	HadRepUniqueIndex         cmap.ConcurrentMap[string, int] // 组成{Gethash(str(initial)): 1}
	HadFinalUniqueIndex       cmap.ConcurrentMap[string, int] // 组成{Gethash(str(initial)): 1}
	HadDisperseUniqueIndex    cmap.ConcurrentMap[string, int] // 组成{Gethash(str(initial)): 1}
	HadReconstructUniqueIndex cmap.ConcurrentMap[string, int] // 组成{Gethash(str(initial)): 1}
	HadReadyUniqueIndex       cmap.ConcurrentMap[string, int] // 组成{Gethash(str(initial)): 1}

	// 签名需要用到的系统参数
	Pairing bls.Pairing
	System  bls.System

	// 拜占庭阈值
	T int
	N int
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

	// Generate key shares.
	const paramsFile = "bls_params.txt"
	var params bls.Params
	if _, err := os.Stat(paramsFile); os.IsNotExist(err) {
		params = bls.GenParamsTypeF(256)
		paramsBytes, _ := params.ToBytes()
		if err := os.WriteFile(paramsFile, paramsBytes, 0600); err != nil {
			log.Printf("[ERROR] 保存参数失败: %v", err)
		}
	} else {
		data, err := os.ReadFile(paramsFile)
		if err != nil {
			log.Printf("[ERROR] 读取参数文件失败: %v", err)
		}
		params, err = bls.ParamsFromBytes(data)
		if err != nil {
			log.Printf("[ERROR] 加载参数失败: %v", err)
		}
		fmt.Printf("Params content: %s\n", string(data))
	}

	// 生成pairing
	pairing := bls.GenPairing(params)
	const systemFile = "bls_system.bin"
	var system bls.System

	if _, err := os.Stat(systemFile); os.IsNotExist(err) {
		system, err = bls.GenSystem(pairing)
		if err != nil {
			log.Printf("[ERROR] 生成系统失败: %v", err)
		}

		systemBytes := system.ToBytes()
		if err := os.WriteFile(systemFile, systemBytes, 0600); err != nil {
			log.Printf("[ERROR] 保存系统失败: %v", err)
		}
	} else {
		systemData, err := os.ReadFile(systemFile)
		if err != nil {
			log.Printf("[ERROR] 读取系统文件失败: %v", err)
		}
		system, err = bls.SystemFromBytes(pairing, systemData)
		if err != nil {
			log.Printf("[ERROR] 加载系统失败: %v", err)
		}
	}

	// 使用固定的随机种子
	rand.Seed(int64(1))

	// 打印当前的 System 的字节表示
	systemBytes := system.ToBytes()
	log.Printf("本次的 System 字节表示: %x\n", systemBytes)

	// 从生成密钥
	groupKey, memberKeys, groupSecret, memberSecrets, err := bls.GenKeyShares(T, N, system)
	if err != nil {
		log.Printf("[ERROR] 生成密钥失败: %v", err)
	}

	return &NodeExtention{
		Node:                      node,
		Config:                    *config,
		BCBSendPool:               make(chan BCBSendMessage, 999999),     // 设置缓冲区大小
		BCBRepPool:                make(chan BCBRepMessage, 999999),      // 设置缓冲区大小
		BCBFinalPool:              make(chan BCBFinalMessage, 999999),    // 设置缓冲区大小
		BCBDispersePool:           make(chan DisperseMessage, 999999),    // 设置缓冲区大小
		BCBReconstructPool:        make(chan ReconstructMessage, 999999), // 设置缓冲区大小
		T:                         T,
		N:                         N,
		GroupKey:                  groupKey,
		MemberKeys:                memberKeys,
		GroupSecret:               groupSecret,
		MemberSecrets:             memberSecrets,
		HadRepUniqueIndex:         cmap.New[int](),
		HadFinalUniqueIndex:       cmap.New[int](),
		HadDisperseUniqueIndex:    cmap.New[int](),
		HadReconstructUniqueIndex: cmap.New[int](),
		HadReadyUniqueIndex:       cmap.New[int](),
		Pairing:                   pairing,
		System:                    system,
	}
}