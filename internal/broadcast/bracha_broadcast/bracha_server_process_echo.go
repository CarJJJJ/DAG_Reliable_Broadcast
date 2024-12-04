package brachabroadcast

import (
	"log"
	"regexp"
)

func (instance *NodeExtention) ProcessEcho() {
	select {
	case msg := <-instance.EchoPool:
		// log.Printf("[INFO] 收到Echo消息: %+v", msg) // 记录收到的Echo消息

		// 假设你有一个接收到的消息
		receivedMessage := msg.InitialMessage.Message

		// 使用正则表达式提取count的值
		re := regexp.MustCompile(`uniqueIndex:(\d+)`)
		matches := re.FindStringSubmatch(receivedMessage)
		uniqueIndex := matches[1]

		// 统计对该initial的echo量级
		if _, exists := instance.EchoCount.Get(uniqueIndex); !exists {
			instance.EchoCount.Set(uniqueIndex, 1)
		} else {
			echoCount, _ := instance.EchoCount.Get(uniqueIndex)
			instance.EchoCount.Set(uniqueIndex, echoCount+1)
		}

		// 如果对该initial的echo量级达到(n+t)/2
		echoCount, _ := instance.EchoCount.Get(uniqueIndex)
		if echoCount >= (instance.N+instance.T)/2 {
			// 如果没echo的就echo
			if _, exists := instance.HadEchoInitial.Get(uniqueIndex); !exists {
				// 封装echo
				echo := EchoMessage{
					Type:           echo_type,
					InitialMessage: msg.InitialMessage,
					NodeID:         instance.Node.Id,
				}
				// 标记已echo
				instance.HadEchoInitial.Set(uniqueIndex, 1)

				// 广播echo
				instance.BroadcastEchoToServers(echo)
			}

			// 如果没ready的就ready
			if _, exists := instance.HadReadyInitial.Get(uniqueIndex); !exists {
				// 封装ready
				ready := ReadyMessage{
					Type:           ready_type,
					InitialMessage: msg.InitialMessage,
					NodeID:         instance.Node.Id,
				}
				// 标记已ready
				instance.HadReadyInitial.Set(uniqueIndex, 1)

				// 广播ready
				instance.BroadcastReadyToServers(ready)
			}
		}

	default:
		log.Println("[INFO] 当前没有Echo消息可处理")
	}
}
