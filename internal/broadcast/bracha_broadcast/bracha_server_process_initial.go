package brachabroadcast

import (
	"log"
	"regexp"
)

func (instance *NodeExtention) ProcessInitial() {
	select {
	case msg := <-instance.InitialPool:
		// log.Printf("[INFO] 收到初始消息: %+v", msg) // 记录收到的初始消息

		// 假设你有一个接收到的消息
		receivedMessage := msg.Message

		// 使用正则表达式提取count的值
		re := regexp.MustCompile(`uniqueIndex:(\d+)`)
		matches := re.FindStringSubmatch(receivedMessage)
		uniqueIndex := matches[1]

		log.Printf("[INFO] uniqueIndex:%s", uniqueIndex)

		// 如果唯一键不存在，放入 map 中并处理消息
		instance.HadEchoInitial.Set(uniqueIndex, 1)

		// 封装Echo
		echo := EchoMessage{
			Type:           echo_type,
			InitialMessage: msg,
			NodeID:         instance.Node.Id,
		}

		// 广播echo
		instance.BroadcastEchoToServers(echo)

	default:
		log.Println("[INFO] 当前没有初始消息可处理") // 如果没有消息，记录日志
	}
}
