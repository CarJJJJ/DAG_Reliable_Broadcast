package dagbroadcast

import (
	"log"
	"strconv"
	"sync"

	"DAG_Reliable_Broadcast/internal/broadcast/util"
)

func (node *NodeExtention) ProcessResponse() {
	select {
	case msg := <-node.ResponsePool:
		// log.Printf("[INFO] 收到Response消息: uniqueIndex: %s, from:%v", msg.UniqueIndex, msg.NodeID)

		responseMessagehash, err := util.CalculateHash(msg)
		if err != nil {
			log.Printf("[ERROR] 计算哈希值失败: %v", err)
			return
		}

		if node.V.Contains(responseMessagehash) {
			log.Printf("[INFO] 已经收到过该Response消息: %s", responseMessagehash)
			return
		}

		// InBuffer删掉这个Send消息
		node.InBuffer.Delete(msg.UniqueIndex)

		// 检查一下这个uniqueIndex可靠接受了没, 因为Inbuffer, 可能会再收到一遍
		if _, ok := node.ReliableBroadcastUniqueIndex.Load(msg.UniqueIndex); ok {
			return
		}

		// 添加节点
		node.V.Add(responseMessagehash)

		// 绘图 默认不开启
		// writeResponseToJSON(node.Node.Id, responseMessagehash, msg)

		// 添加边
		for _, echoReference := range msg.EchoReferences {
			node.E.Add(echoReference)
		}
		for _, readyReference := range msg.ReadyReferences {
			node.E.Add(readyReference)
		}

		// 添加HashToResponse
		node.HashToResponse.Store(responseMessagehash, &msg)

		// 添加 NodeIDToUniqueIndexToResponseHashSlice
		uniqueIndexToResponseHashSliceMap, ok := node.NodeIDToUniqueIndexToResponseHashSlice.Load(msg.NodeID)
		if ok {
			// 拿到对应的NodeID下的Map
			uniqueIndexToResponseHashSliceMap := uniqueIndexToResponseHashSliceMap.(*sync.Map)
			// 拿到对应的UniqueIndex下的哈希切片
			responseHashSlice, ok := uniqueIndexToResponseHashSliceMap.Load(msg.UniqueIndex)
			if !ok {
				responseHashSlice = &[]string{responseMessagehash}
			} else {
				responseHashSlice = append(*responseHashSlice.(*[]string), responseMessagehash)
			}
			uniqueIndexToResponseHashSliceMap.Store(msg.UniqueIndex, responseHashSlice)
		}

		// 添加VK
		responseHashSlice, ok := node.VK.Load(msg.UniqueIndex)
		if ok {
			responseHashSlice := responseHashSlice.(*[]string)
			*responseHashSlice = append(*responseHashSlice, responseMessagehash)
		} else {
			responseHashSlice := []string{responseMessagehash}
			node.VK.Store(msg.UniqueIndex, &responseHashSlice)
		}

		// 更新Vote
		node.UpdateVote(msg, responseMessagehash)

		// 统计Ready
		node.CountReady()

	default:
		log.Println("[INFO] 当前没有Response消息可处理") // 如果没有消息，记录日志
	}
}

func (node *NodeExtention) UpdateVote(responseMessage ResponseMessage, responseMessageHash string) {
	node.Vote.Range(func(key, value interface{}) bool {
		voteValue := value.(*VoteValue)
		echoContains := false
		readyContains := false
		for _, responseHash := range responseMessage.EchoReferences {
			if voteValue.References.Contains(responseHash) {
				voteValue.IDS_ECHO.Add(strconv.Itoa(responseMessage.NodeID))
				echoContains = true
			}
		}

		for _, responseHash := range responseMessage.ReadyReferences {
			if voteValue.References.Contains(responseHash) {
				voteValue.IDS_READY.Add(strconv.Itoa(responseMessage.NodeID))
				readyContains = true
			}
		}

		if echoContains || readyContains {
			voteValue.References.Add(responseMessageHash)
		}

		return true
	})
	reference := NewThreadSafeSet()
	reference.Add(responseMessageHash)
	idsEcho := NewThreadSafeSet()
	idsEcho.Add(strconv.Itoa(responseMessage.NodeID))
	idsReady := NewThreadSafeSet()
	idsReady.Add(strconv.Itoa(responseMessage.NodeID))
	voteValue := &VoteValue{
		References: reference,
		IDS_ECHO:   idsEcho,
		IDS_READY:  idsReady,
	}
	node.Vote.Store(responseMessageHash, voteValue)
}

func (node *NodeExtention) CountReady() {
	node.VK.Range(func(key, value interface{}) bool {
		uniqueIndex := key.(string)
		responseHashSlice := value.(*[]string)
		readyNodeIds := NewThreadSafeSet()
		for _, responseHash := range *responseHashSlice {
			voteValue, ok := node.Vote.Load(responseHash)
			if ok {
				voteValue := voteValue.(*VoteValue)
				IDS_READY := voteValue.IDS_READY
				for _, nodeId := range IDS_READY.ToSlice() {
					readyNodeIds.Add(nodeId)
				}
			}
		}

		if readyNodeIds.Length() >= 2*node.T+1 {
			// log.Printf("[INFO] 可靠广播消息: uniqueIndex: %s, readyNodeIds: %v", uniqueIndex, readyNodeIds.ToSlice())
			_, loaded := node.ReliableBroadcastUniqueIndex.LoadOrStore(uniqueIndex, 1)
			if !loaded {
				node.ReliableBroadcastCount++
			}

			// 删除Vote
			for _, responseHash := range *responseHashSlice {
				node.Vote.Delete(responseHash)
			}
			// 删除VK
			node.VK.Delete(uniqueIndex)

			// 删除NodeIDToUniqueIndexToResponseHashSlice
			node.NodeIDToUniqueIndexToResponseHashSlice.Range(func(key, value interface{}) bool {
				uniqueIndexToResponseHashSliceMap := value.(*sync.Map)
				uniqueIndexToResponseHashSliceMap.Delete(uniqueIndex)
				return true
			})
		}
		return true
	})
}
