package brachabroadcast

import (
	"log"
	"regexp"
)

func (instance *NodeExtention) ProcessReady() {
	select {
	case msg := <-instance.ReadyPool:
		// log.Printf("[INFO] 收到Ready消息: %+v", msg) // 记录Ready

		// 假设你有一个接收到的消息
		receivedMessage := msg.InitialMessage.Message

		// 使用正则表达式提取count的值
		re := regexp.MustCompile(`uniqueIndex:(\d+)`)
		matches := re.FindStringSubmatch(receivedMessage)
		uniqueIndex := matches[1]

		// 统计收到当前initial的ready量级
		if _, exists := instance.ReadyCount.Get(uniqueIndex); !exists {
			instance.ReadyCount.Set(uniqueIndex, 1)
		} else {
			readyCount, _ := instance.ReadyCount.Get(uniqueIndex)
			instance.ReadyCount.Set(uniqueIndex, readyCount+1)
		}

		// 当对该initial的ready量级到达t+1
		readyCount, _ := instance.ReadyCount.Get(uniqueIndex)
		if readyCount >= instance.T+1 {
			// 没echo就echo
			if _, exists := instance.HadEchoInitial.Get(uniqueIndex); !exists {
				// 封装echo
				echo := EchoMessage{
					Type:           echo_type,
					InitialMessage: msg.InitialMessage,
					NodeID:         instance.Node.Id,
				}

				// 标记该initial被echo
				instance.HadEchoInitial.Set(uniqueIndex, 1)

				// 广播echo
				instance.BroadcastEchoToServers(echo)
			}

			// 没ready就ready
			if _, exists := instance.HadReadyInitial.Get(uniqueIndex); !exists {
				// 封装echo
				ready := ReadyMessage{
					Type:           ready_type,
					InitialMessage: msg.InitialMessage,
					NodeID:         instance.Node.Id,
				}

				// 标记该initial被ready
				instance.HadReadyInitial.Set(uniqueIndex, 1)

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
