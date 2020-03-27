package comet

import (
	"fmt"
	"time"
)

// package comet
func ChatKey(id int64) string {
	return fmt.Sprintf("chat%d", id)
}
func ClientChatsKey(id int64) string {
	return fmt.Sprintf("client%d_chats", id)
}
func (comet CometServer) EmitToChat(id int64, event string, data interface{}) error {
	var ids = []int64{}
	begin := time.Now()
	err := comet.Storage.SMembers(ChatKey(id)).ScanSlice(&ids)
	go fmt.Println("Scanned after: ", time.Since(begin))
	if err != nil {
		fmt.Println("error: ", err.Error())
		return err
	}
	for _, id := range ids {
		go comet.ConnectionsStorage.Emit(id, event, data)
	}
	fmt.Println("Time to send to chats: ", time.Since(begin))
	return nil
}
func (conn *CometConnection) JoinChat(chatId int64) {
	if conn == nil {
		return
	}
	if conn.Id < 1 || conn.Examples == nil {
		return
	}
	if conn.Examples.Server == nil {
		return
	}
	go conn.Examples.Server.Storage.SAdd(ClientChatsKey(conn.Id), chatId)
	go conn.Examples.Server.Storage.SAdd(ChatKey(chatId), conn.Id)
}
func (conns *ConnectionExamples) JoinChat(chatId int64) {
	if conns == nil {
		fmt.Println("err:0")
		return
	}
	if conns.Id < 1 || conns.Server == nil {
		fmt.Println("err:1")
		return
	}
	go conns.Server.Storage.SAdd(ClientChatsKey(conns.Id), chatId)
	go conns.Server.Storage.SAdd(ChatKey(chatId), conns.Id)
}

func (conn *CometConnection) LeaveChat(chatId int64) {
	if conn == nil {
		return
	}
	if conn.Id < 1 || conn.Examples == nil {
		return
	}
	if conn.Examples.Server == nil {
		go conn.Examples.Server.Storage.SRem(ChatKey(chatId), conn.Id)
	}

}
func (conn *CometConnection) GetChats() *[]int64 {
	if conn == nil {
		return nil
	}
	if conn.Id < 1 || conn.Examples == nil {
		return nil
	}
	if conn.Examples.Server == nil {
		return nil
	}
	var ids []int64
	err := conn.Examples.Server.Storage.SMembers(ClientChatsKey(conn.Id)).ScanSlice(&ids)
	if err != nil {
		return nil
	}
	return &ids
}
func (conn *CometConnection) LeaveChats() {
	if conn == nil {
		return
	}
	connId := conn.Id
	if connId < 1 || conn.Examples == nil {
		return
	}
	chats := conn.GetChats()
	if chats == nil || conn.Examples.Server == nil {
		return
	}
	for _, id := range *chats {
		go conn.Examples.Server.Storage.SRem(ChatKey(id), connId)
	}
}
