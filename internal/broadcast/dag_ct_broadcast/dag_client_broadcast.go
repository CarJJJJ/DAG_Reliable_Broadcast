package dagctbroadcast

import (
	"DAG_Reliable_Broadcast/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/klauspost/reedsolomon"
)

func BroadcastToServers(node Node) {
	config, err := config.LoadConfig(ConfigPath)
	if err != nil {
		log.Printf("[ERROR] 加载配置失败: %v", err)
		return
	}

	N := config.N
	T := config.T

	encoder, err := reedsolomon.New(N-T, T)
	if err != nil {
		log.Printf("[ERROR] 纠删码编码器创建失败: %v", err)
	}
	ticker := time.NewTicker(1 * time.Second) // 每秒触发一次
	defer ticker.Stop()                       // 确保在函数结束时停止计时器
	// 定义n的大小
	const n = 1024 // 例如，n = 1024
	// 生成n*n字节的随机消息内容
	messageContent := make([]byte, n*n) // 创建n*n B的字节切片
	for i := range messageContent {
		messageContent[i] = byte('A' + rand.Intn(26))
	}

	// 确保messageContent的长度是n*n
	if len(messageContent) < n*n {
		messageContent = append(messageContent, make([]byte, n*n-len(messageContent))...)
	} else if len(messageContent) > n*n {
		messageContent = messageContent[:n*n]
	}
	time.Sleep(10 * time.Second) // 等待10s连接建立完再发消息

	count := 1
	for range ticker.C {
		if count >= 10 {
			break
		}
		for i := 0; i < 2; i++ { // 每秒发送5条消息
			// 纠删码对消息进行分片

			// 计算每个分片的大小
			shardSize := (len(messageContent) + N - 1) / N

			// 创建数据分片
			data := make([][]byte, N)
			for i := range data {
				data[i] = make([]byte, shardSize)
			}

			// 将消息填充到数据分片中
			for i := range messageContent {
				data[i%N][i/N] = messageContent[i]
			}

			// 纠删码对数据分片进行编码
			err = encoder.Encode(data)
			if err != nil {
				log.Printf("[ERROR] 纠删码分片失败: %v", err)
				continue
			}

			for i := 0; i < N; i++ {
				sendMessage := SendMessage{
					Type:        0,                              // 设置消息类型
					Message:     string(data[i]),                // 分片消息
					NodeID:      node.Id,                        // 设置节点ID
					UniqueIndex: fmt.Sprintf("%d_%d", count, i), // 设置唯一索引
				}

				message, err := json.Marshal(sendMessage) // 序列化为JSON
				if err != nil {
					log.Printf("[ERROR] 序列化消息失败: %v", err)
					continue
				}

				message = append(message, '\n')
				serverAddr := config.Servers[i]
				Host := serverAddr.Host
				Port := serverAddr.Port
				serverAddrStr := fmt.Sprintf("%s:%s", Host, Port)
				conn, ok := node.Conn[serverAddrStr]
				if !ok {
					log.Printf("[ERROR] 连接到 %s 失败", serverAddrStr)
					continue
				}
				conn.Write(message)
				// log.Print("[INFO] 广播消息, uniqueIndex: ", count)
			}
			count++
		}
	}
}
