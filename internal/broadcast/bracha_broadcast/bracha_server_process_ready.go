package brachabroadcast

import (
	"DAG_Reliable_Broadcast/internal/broadcast/util"
	"log"
)

func (instance *NodeExtention) ProcessReady() {
	select {
	case msg := <-instance.ReadyPool:
		// log.Printf("[INFO] 收到Ready消息: %+v", msg) // 记录Ready

		// 计算initial的hash
		hash, err := util.CalculateHash(msg.InitialMessage)
		if err != nil {
			log.Printf("[ERROR] 计算哈希值失败: %v", err) // 记录错误
			return
		}

		// 统计收到当前initial的ready量级
		if _, exists := instance.ReadyCount.Get(hash); !exists {
			instance.ReadyCount.Set(hash, 1)
		} else {
			readyCount, _ := instance.ReadyCount.Get(hash)
			instance.ReadyCount.Set(hash, readyCount+1)
		}

		// 当对该initial的ready量级到达t+1
		readyCount, _ := instance.ReadyCount.Get(hash)
		if readyCount >= instance.T+1 {
			// 没echo就echo
			if _, exists := instance.HadEchoInitial.Get(hash); !exists {
				// 封装echo
				echo := EchoMessage{
					Type:           echo_type,
					InitialMessage: msg.InitialMessage,
					NodeID:         instance.Node.Id,
				}

				// 标记该initial被echo
				instance.HadEchoInitial.Set(hash, 1)

				// 广播echo
				instance.BroadcastEchoToServers(echo)
			}

			// 没ready就ready
			if _, exists := instance.HadReadyInitial.Get(hash); !exists {
				// 封装echo
				ready := ReadyMessage{
					Type:           ready_type,
					InitialMessage: msg.InitialMessage,
					NodeID:         instance.Node.Id,
				}

				// 标记该initial被ready
				instance.HadReadyInitial.Set(hash, 1)

				// 广播ready
				instance.BroadcastReadyToServers(ready)
			}
		}

		// 当对该initial的ready量级到达2t+1
		if readyCount >= 2*instance.T+1 {
			// log.Printf("消息已被可靠广播, 消息:%v", msg.InitialMessage)
			instance.ReliableBroadcastCount++
		}

	default:
		log.Println("[INFO] 当前没有Ready消息可处理") // 如果没有消息，记录日志
	}
}
