package brachabroadcast

import (
	"DAG_Reliable_Broadcast/internal/broadcast/util"
	"log"
)

func (instance *NodeExtention) ProcessEcho() {
	select {
	case msg := <-instance.EchoPool:
		// log.Printf("[INFO] 收到Echo消息: %+v", msg) // 记录收到的Echo消息

		// 计算哈希值
		hash, err := util.CalculateHash(msg.InitialMessage)
		if err != nil {
			log.Printf("[ERROR] 计算哈希值失败: %v", err) // 记录错误
			return
		}

		// 统计对该initial的echo量级
		if _, exists := instance.EchoCount.Get(hash); !exists {
			instance.EchoCount.Set(hash, 1)
		} else {
			echoCount, _ := instance.EchoCount.Get(hash)
			instance.EchoCount.Set(hash, echoCount+1)
		}

		// 如果对该initial的echo量级达到(n+t)/2
		echoCount, _ := instance.EchoCount.Get(hash)
		if echoCount >= (instance.N+instance.T)/2 {
			// 如果没echo的就echo
			if _, exists := instance.HadEchoInitial.Get(hash); !exists {
				// 封装echo
				echo := EchoMessage{
					Type:           echo_type,
					InitialMessage: msg.InitialMessage,
					NodeID:         instance.Node.Id,
				}
				// 标记已echo
				instance.HadEchoInitial.Set(hash, 1)

				// 广播echo
				instance.BroadcastEchoToServers(echo)
			}

			// 如果没ready的就ready
			if _, exists := instance.HadReadyInitial.Get(hash); !exists {
				// 封装ready
				ready := ReadyMessage{
					Type:           ready_type,
					InitialMessage: msg.InitialMessage,
					NodeID:         instance.Node.Id,
				}
				// 标记已ready
				instance.HadReadyInitial.Set(hash, 1)

				// 广播ready
				instance.BroadcastReadyToServers(ready)
			}
		}

	default:
		log.Println("[INFO] 当前没有Echo消息可处理") // 如果没有消息，记录日志
	}
}
