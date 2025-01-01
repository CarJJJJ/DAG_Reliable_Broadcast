package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

type GraphNode struct {
	ResponseHash    string   `json:"ResponseHash"`
	EchoResponse    []string `json:"echoResponse"`
	ReadyResponse   []string `json:"readyResponse"`
	NodeDescription string   `json:"nodeDescription"`
}

func TestGenerateDAGGraph(t *testing.T) {
	// 读取 JSON 文件
	data, err := os.ReadFile("DAG_Graph.json")
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}

	// 解析 JSON
	var nodes []GraphNode
	if err := json.Unmarshal(data, &nodes); err != nil {
		t.Fatalf("解析 JSON 失败: %v", err)
	}

	// 创建输出文件
	outFile, err := os.Create("DAG_Graph.dot")
	if err != nil {
		t.Fatalf("创建文件失败: %v", err)
	}
	defer outFile.Close()

	// 生成 DOT 语言代码
	fmt.Fprintln(outFile, "digraph DAG {")
	fmt.Fprintln(outFile, "    // 设置图形属性")
	fmt.Fprintln(outFile, "    rankdir=LR;")       // 从左到右的布局
	fmt.Fprintln(outFile, "    node [shape=box];") // 节点使用方框形状

	// 创建节点映射
	nodeMap := make(map[string]string)
	for _, node := range nodes {
		nodeMap[node.ResponseHash] = node.NodeDescription
	}

	// 添加节点
	for _, node := range nodes {
		fmt.Fprintf(outFile, "    \"%s\" [label=\"%s\"];\n",
			node.NodeDescription,
			node.NodeDescription)
	}

	// 添加边
	for _, node := range nodes {
		// Echo 引用 - 绿色边
		for _, echo := range node.EchoResponse {
			if targetDesc, ok := nodeMap[echo]; ok {
				fmt.Fprintf(outFile, "    \"%s\" -> \"%s\" [color=green];\n",
					node.NodeDescription,
					targetDesc)
			}
		}

		// Ready 引用 - 红色边
		for _, ready := range node.ReadyResponse {
			if targetDesc, ok := nodeMap[ready]; ok {
				fmt.Fprintf(outFile, "    \"%s\" -> \"%s\" [color=red];\n",
					node.NodeDescription,
					targetDesc)
			}
		}
	}

	fmt.Fprintln(outFile, "}")

	t.Log("已生成 DAG 图的 DOT 文件：DAG_Graph.dot")
	t.Log("请使用以下命令生成图片：")
	t.Log("dot -Tpng DAG_Graph.dot -o DAG_Graph.png")
}
