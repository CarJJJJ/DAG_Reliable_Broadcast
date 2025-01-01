package dagctbroadcast

import (
	"log"
	"sort"
	"sync"
)

func (node *NodeExtention) ProcessSend() {
	// 直接从通道接收消息
	msg := <-node.SendPool
	uniqueIndex := msg.UniqueIndex
	log.Printf("[INFO] 收到Send消息: uniqueIndex: %s, from:%v", uniqueIndex, msg.NodeID)

	// 1. 选择下两个要发送消息的节点
	selectedNodes := node.SelectNextNodes()

	// 2. 选择NodeID下未的UniqueIndex未被Echo的消息
	echoReferences := node.chooseEchoReferences(selectedNodes) // 至多2条Echo,至少0条

	// 附加：如果echoReferences为空，那么重选
	for len(echoReferences) == 0 {
		selectedNodes = node.SelectNextNodes()
		echoReferences = node.chooseEchoReferences(selectedNodes) // 至多2条Echo,至少0条
	}

	// 3. 选择Ready的UniqueIndex
	readyReferences := node.chooseReadyReferences() // 只有1条Ready

	// 4. 发送消息
	responseMessage := ResponseMessage{
		Type:            1,
		SendMessage:     msg,
		EchoReferences:  echoReferences,
		ReadyReferences: readyReferences,
		NodeID:          node.Node.Id,
		UniqueIndex:     uniqueIndex,
	}
	node.BroadcastResposneToServers(&responseMessage)
}

// SelectNextNodes 使用 RobinRound 算法选择下两个要发送消息的节点
func (node *NodeExtention) SelectNextNodes() []int {
	// 1.1 获取当前节点的 ID
	currentNodeID := node.Node.Id

	// 1.2 计算可用的节点（排除自身）
	availableNodes := make([]int, 0, len(node.RobinRound)-1)
	for _, id := range node.RobinRound {
		if id != currentNodeID {
			availableNodes = append(availableNodes, id)
		}
	}

	// 1.3 选择两个节点
	selectedNodes := make([]int, 0, 2)
	for i := 0; i < 2; i++ {
		selectedNode := availableNodes[(node.CurrentIndex+i)%len(availableNodes)]
		selectedNodes = append(selectedNodes, selectedNode)
	}

	// 1.4 更新 currentIndex 为下次选择的起始位置
	node.CurrentIndex = (node.CurrentIndex + 1) % len(availableNodes)

	return selectedNodes
}

// 选择Echo引用
func (node *NodeExtention) chooseEchoReferences(selectedNodes []int) []string {
	echoReferences := []string{}
	for _, nodeID := range selectedNodes {
		if uniqueIndexToResponseHash, ok := node.NodeIDToUniqueIndexToResponseHashSlice.Load(nodeID); ok {
			uniqueIndexToResponseHashMap := uniqueIndexToResponseHash.(*sync.Map)

			// 收集所有 uniqueIndex
			var uniqueIndexes []string
			uniqueIndexToResponseHashMap.Range(func(key, value interface{}) bool {
				uniqueIndexes = append(uniqueIndexes, key.(string))
				return true
			})

			// 排序 uniqueIndex
			sort.Strings(uniqueIndexes)

			// 按顺序处理
			for _, uniqueIndex := range uniqueIndexes {
				value, ok := uniqueIndexToResponseHashMap.Load(uniqueIndex)
				if !ok {
					continue
				}
				responseHashes := value.(*[]string)

				if _, ok := node.HadEchoUniqueIndex.Load(uniqueIndex); ok {
					// 如果当前的UniqueIndex已经被Echo过了,那就看下一个
					continue
				}

				// 如果没有Echo过，那就直接选靠前加入的那个
				echoReferences = append(echoReferences, (*responseHashes)[0])

				// 标记当前节点已经Echo过了这个UniqueIndex
				node.HadEchoUniqueIndex.Store(uniqueIndex, 1)

				// 找到一个就退出
				break
			}
		}
	}
	return echoReferences
}

func (node *NodeExtention) chooseReadyReferences() []string {
	readyReferences := []string{}

	// 收集所有 uniqueIndex
	var uniqueIndexes []string
	node.VK.Range(func(key, value interface{}) bool {
		uniqueIndexes = append(uniqueIndexes, key.(string))
		return true
	})

	// 排序 uniqueIndex
	sort.Strings(uniqueIndexes)

	// 按顺序处理每个 uniqueIndex
	for _, uniqueIndex := range uniqueIndexes {
		value, ok := node.VK.Load(uniqueIndex)
		if !ok {
			continue
		}
		responseHashSlice := value.(*[]string)

		// 看看这个UniqueIndex在不在HadReadyUniqueIndex中
		if _, ok := node.HadReadyUniqueIndex.Load(uniqueIndex); ok {
			// 如果在,那么就看下一个uniqueIndex
			continue
		}

		// 如果不在,那么就看这个UniqueIndex的Response消息是否被不同的服务节点Echo 2t+1次或ready t+1次
		echoNodeIds := NewThreadSafeSet()
		readyNodeIds := NewThreadSafeSet()
		for _, responseHash := range *responseHashSlice {
			voteValue, ok := node.Vote.Load(responseHash)
			if ok {
				voteValue := voteValue.(*VoteValue)
				IDS_ECHO := voteValue.IDS_ECHO
				IDS_READY := voteValue.IDS_READY
				for _, nodeId := range IDS_ECHO.ToSlice() {
					echoNodeIds.Add(nodeId)
				}
				for _, nodeId := range IDS_READY.ToSlice() {
					readyNodeIds.Add(nodeId)
				}
			}
		}

		// 选Ready的时候不要选和当前nodeid一样的节点哈希, 如果有nodeid一样的哈希, 那么就去看下一个uniqueIndex了
		checkNextUniqueIndex := false
		for _, responseHash := range *responseHashSlice {
			responseMessage, ok := node.HashToResponse.Load(responseHash)
			if ok {
				responseMessage := responseMessage.(*ResponseMessage)
				if responseMessage.NodeID == node.Node.Id {
					checkNextUniqueIndex = true
					break
				}
			}
		}
		if checkNextUniqueIndex {
			continue
		}

		// 没有问题，那么就选Ready
		if echoNodeIds.Length() >= 2*node.T+1 || readyNodeIds.Length() >= node.T+1 {
			hadChooseReady := false
			for _, responseHash := range *responseHashSlice {
				responseMessage, ok := node.HashToResponse.Load(responseHash)
				if ok {
					responseMessage := responseMessage.(*ResponseMessage)
					if responseMessage.NodeID != node.Node.Id {
						readyReferences = append(readyReferences, responseHash)
						hadChooseReady = true
						// 在HadReadyUniqueIndex中添加键值对{UniqueIndex: 1}, 表明当前节点已经对该UniqueIndex的Response消息Ready过了
						node.HadReadyUniqueIndex.Store(uniqueIndex, 1)
						// break会退出当前uniqueIndex哈希切片的循环,接着会去看下一个uniqueIndex
						break
					}
				}
			}

			if hadChooseReady {
				// Ready只选一个
				break
			}

			if !hadChooseReady {
				// 这个uniqueIndex的Ready选不了,因为节点Id一样，所以就看下一个UniqueIndex
				continue
			}
		}
	}
	return readyReferences
}
