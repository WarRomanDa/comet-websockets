package comet

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type CometConnection struct {
	Conn     *websocket.Conn
	Examples *ConnectionExamples
	Id       int64
	reader   sync.Mutex
	writer   sync.Mutex
}

func (conn *CometConnection) Emit(event string, data interface{}) {
	conn.writer.Lock()
	defer conn.writer.Unlock()
	fmt.Println("Emit error: ", conn.Conn.WriteJSON(map[string]interface{}{
		"event": event,
		"data":  data,
	}))
}
func (conn *CometConnection) Close() {
	if conn.Examples != nil {
		if conn.Examples.Server != nil {
			conn.Examples.Server.OnDisconnect(conn)
			conn.Examples.Server.ConnectionsStorage.Delete(conn)
		}
	}
	conn.LeaveChats()
	if conn.Conn != nil {
		conn.Conn.Close()
	}
}
func (comet *CometServer) NewConnection(id int64, wsConn *websocket.Conn) *CometConnection {
	if wsConn == nil || id < 1 {
		return nil
	}
	return &CometConnection{
		reader:   sync.Mutex{},
		writer:   sync.Mutex{},
		Conn:     wsConn,
		Examples: nil,
		Id:       id,
	}
}

func (conn *CometConnection) Recieve() {
	for conn != nil {
		conn.reader.Lock()
		var data Tx
		if err := conn.Conn.ReadJSON(&data); err != nil {
			fmt.Println("error: ", err.Error())
			conn.reader.Unlock()
			conn.Close()
			return
		}
		conn.reader.Unlock()
		if conn.Examples != nil {
			if conn.Examples.Server != nil {
				conn.Examples.Server.HandleEvent(conn, &data)
			}
		}
	}
	conn.Close()
}
