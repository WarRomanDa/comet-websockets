package comet

import "fmt"

type OnDisconnectDelegate func(conn *CometConnection)

func DefaultOnDisconnectDelegate() OnDisconnectDelegate {
	return func(conn *CometConnection) {
		fmt.Println("Disconnected: ", conn.Id)
	}
}
