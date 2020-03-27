package comet

import "net/http"

type ConnectionBuilder struct {
	Id int64
	Request *http.Request
}