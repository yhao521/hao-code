package backend

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源（开发环境）
	},
}

// TerminalWebSocketHandler 处理终端 WebSocket 连接
func TerminalWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// 创建一个新的终端会话
	session, err := NewTerminalSession()
	if err != nil {
		log.Printf("Failed to create terminal session: %v", err)
		return
	}
	defer session.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// 协程 1: 从 WebSocket 读取数据并写入 PTY (前端 -> 后端)
	go func() {
		defer wg.Done()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			session.Write(message)
		}
	}()

	// 协程 2: 从 PTY 读取数据并发送到 WebSocket (后端 -> 前端)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 4096)
		for {
			n, err := session.Pty.Read(buffer)
			if err != nil {
				break
			}
			if n > 0 {
				err := conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
				if err != nil {
					break
				}
			}
		}
	}()

	wg.Wait()
}
