package comet

import (
	"fmt"
	"sync"
)

type ConnectionsBuffer struct {
	Connections map[int64]*ConnectionExamples
	mutex       sync.RWMutex
}

func (cb *ConnectionsBuffer) Emit(id int64, event string, data interface{}) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	if cb.Connections[id] != nil {
		cb.Connections[id].Emit(event, data)
	} else {
		fmt.Println("no such connection")
	}
}
func (cb *ConnectionsBuffer) AddConnection(conn *CometConnection, server *CometServer) {
	if server == nil {
		conn.Close()
		return
	}
	cb.mutex.Lock()
	id := conn.Id
	if cb.Connections[id] == nil {
		cb.Connections[id] = NewConnectionExamples(id, server)
	}
	cb.Connections[id].Add(conn)
	conn.Examples = cb.Connections[id]
	cb.mutex.Unlock()
	go conn.Recieve()
}
func (cb *ConnectionsBuffer) Delete(conn *CometConnection) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	id := conn.Id
	if cb.Connections[id] != nil {
		cb.Connections[id].Delete(conn)
	}
}
func NewConnectionsBuffer() *ConnectionsBuffer {
	return &ConnectionsBuffer{
		Connections: make(map[int64]*ConnectionExamples),
		mutex:       sync.RWMutex{},
	}
}
