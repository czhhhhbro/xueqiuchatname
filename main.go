package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"time"
)

// 消息结构
type Message struct {
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Self      bool      `json:"self"`
}

var (
	messages []Message
	mu       sync.Mutex
)

func main() {
	// 首页
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// 发送消息
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var msg Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		msg.Timestamp = time.Now()

		mu.Lock()
		messages = append(messages, msg)
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	})

	// 获取消息列表
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		list := messages
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list)
	})

	// 端口适配云平台
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	http.ListenAndServe(":"+port, nil)
}
