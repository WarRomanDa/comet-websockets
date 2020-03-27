package comet

// package comet
type OnConnectDelegate func(*ConnectionBuilder) bool

func DefaultOnConnectDelegate() OnConnectDelegate {
	return func(builder *ConnectionBuilder) bool {
		return false
	}
}
