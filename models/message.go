package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FormId   int64  //发送者
	TargetID int64  // 接受者
	Type     int    // 传播类型 群聊 私聊 广播
	Media    int    // 消息类型 文字 图片 音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

// 表示一个客户端连接
type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte   // 消息发送队列，存储需要发送的数据。
	GroupSets set.Interface // 表示用户参与的群组集合。
}

// 映射关系
// 是一个全局映射，维护用户 ID 和对应连接节点 (Node) 的关系。
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁 用于保护 clientMap 的并发读写
var rwLocker sync.RWMutex

// 发送者ID 接受者ID 消息类型 发生的内容 发送类型
// 它的功能是接受用户的连接请求，建立 WebSocket 连接，并设置用户与服务端的交互机制。
func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1. 获取参数并校验合法性
	// token := query.Get("token")
	query := request.URL.Query()
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	isvalid := true

	// 使用 websocket.Upgrader 将普通 HTTP 请求升级为 WebSocket 连接。
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建一个 Node 对象表示当前用户的连接，
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 3. 用户关系
	// 4. userid 与node绑定并且加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5. 完成发送的逻辑
	go sendProc(node)
	// 6. 完成接受的逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎来到聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		// wait for data to be available on channel
		case data := <-node.DataQueue: // received and stored in the data variable
			// send the received data
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("接受到消息[ws] >>>>>", data)
	}
}

// queue messages that need to be broadcasted
var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

// automatically called when the package is initialized
// 只要这个包被导入到主程序中或直接作为主包运行，init() 都会被自动执行。
func init() {
	go udpSendProc()
	go udpRecvProc()
}

// Implementation for sending UDP message
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 122, 255),
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// Implementation for reciving UDP message
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// backend scheduling and processing logic
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: // private message
		sendMsg(msg.TargetID, data)
		// case 2:// bulk message
		// 	sendGroupMsg()
		// case 3:// broadcast
		// 	sendAllMsg()
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
