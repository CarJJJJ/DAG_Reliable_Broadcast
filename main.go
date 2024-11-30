package main

import (
	"flag"
	"log"
	"os"

	networknode "DAG_Reliable_Broadcast/internal/network/network_node"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	nodeType := flag.String("type", "", "节点类型 (0:客户端, 1:服务端)")
	host := flag.String("host", "0.0.0.0", "服务器主机地址")
	port := flag.String("port", "8080", "服务器端口")
	configPath := flag.String("config", "config/host_config.json", "配置文件路径")
	broadcastType := flag.Int("broadcastType", 0, "广播类型")
	id := flag.Int("id", 1, "节点id")
	flag.Parse()

	switch *nodeType {
	case "0":
		networknode.StartClient(*configPath, *broadcastType, *id)
	case "1":
		networknode.StartServer(*host, *port, *broadcastType, *id)
	default:
		log.Println("[INFO] 使用方法: -type {0|1} [-config config.json] [-host host] [-port port]")
		log.Println("[INFO] 0 - 客户端")
		log.Println("[INFO] 1 - 服务端")
		os.Exit(1)
	}
}
