package comet

import (
	"fmt"
	"sync"
)

type ConnectionExamples struct {
	Examples map[*CometConnection]bool
	Id       int64
	mutex    sync.Mutex
	Server   *CometServer
}

func (conns *ConnectionExamples) Emit(event string, data interface{}) {
	conns.mutex.Lock()
	fmt.Println("sending started...")

	for conn, _ := range conns.Examples {
		fmt.Println("sending to user: ", conn.Id)
		if conn == nil {
			delete(conns.Examples, conn)
			continue
		}
		go conn.Emit(event, data)
	}
	conns.mutex.Unlock()
}
func (conns *ConnectionExamples) Add(conn *CometConnection) {
	conns.mutex.Lock()
	defer conns.mutex.Unlock()
	conns.Examples[conn] = true
}
func (conns *ConnectionExamples) Delete(conn *CometConnection) {
	conns.mutex.Lock()
	defer conns.mutex.Unlock()
	delete(conns.Examples, conn)
}
func NewConnectionExamples(id int64, server *CometServer) *ConnectionExamples {
	return &ConnectionExamples{
		Examples: make(map[*CometConnection]bool),
		mutex:    sync.Mutex{},
		Id:       id,
		Server:   server,
	}
}
