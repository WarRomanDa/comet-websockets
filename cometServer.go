package comet

import (
	"messenger-application/db/models"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

// package comet
type CometServer struct {
	Handlers           map[string]EventHandler
	ConnectionsStorage *ConnectionsBuffer
	Upgrader           websocket.Upgrader
	OnConnect          OnConnectDelegate
	OnDisconnect       OnDisconnectDelegate
	Storage            *redis.Client
}

func NewCometServer(redisUrl string, redisPassword string) *CometServer {
	return &CometServer{
		Handlers:           make(map[string]EventHandler),
		ConnectionsStorage: NewConnectionsBuffer(),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  0,
			WriteBufferSize: 0,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		OnConnect:    DefaultOnConnectDelegate(),
		OnDisconnect: DefaultOnDisconnectDelegate(),
		Storage:      NewRedisStorage(redisUrl, redisPassword),
	}
}

func (comet *CometServer) On(event string, handler EventHandler) {
	comet.Handlers[event] = handler
}
func (comet *CometServer) HandleEvent(conn *CometConnection, tx *Tx) {
	if handler := comet.Handlers[tx.Event]; handler != nil {
		handler(conn, tx.Data)
	}
}
func (comet *CometServer) GetServer() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		wsConn, err := comet.Upgrader.Upgrade(res, req, nil)
		if err != nil {
			wsConn.Close()
			return
		}
		builder := ConnectionBuilder{
			Id:      -1,
			Request: req,
		}
		if comet.OnConnect(&builder) {
			if conn := comet.NewConnection(builder.Id, wsConn); conn != nil {
				comet.ConnectionsStorage.AddConnection(conn, comet)
				return
			}
		}
		wsConn.Close()
	}
}
func (comet *CometServer) JoinChats(userId int64, chats *[]*models.Chat) {
	if chats == nil {
		return
	}
	comet.ConnectionsStorage.mutex.Lock()
	conn := comet.ConnectionsStorage.Connections[userId]
	comet.ConnectionsStorage.mutex.Unlock()
	if conn == nil {
		return
	}
	for _, chat := range *chats {
		if chat != nil {
			conn.Emit("joined_chat", chat.Id)
			conn.JoinChat(chat.Id)
		}
	}
}
func (comet *CometServer) JoinChat(userId int64, chatId int64) {
	comet.ConnectionsStorage.mutex.Lock()
	conn := comet.ConnectionsStorage.Connections[userId]
	comet.ConnectionsStorage.mutex.Unlock()
	if conn == nil {
		return
	}
	conn.JoinChat(chatId)
}
func (comet *CometServer) Emit(id int64, event string, data interface{}) {
	comet.ConnectionsStorage.Emit(id, event, data)
}
