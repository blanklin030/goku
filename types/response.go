package types

import (
	"encoding/json"
)

// APIResponseContext 响应内容
type APIResponseContext struct {
	RequestID string      `json:"request_id"`
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Error     string      `json:"error"`
}

func (a *APIResponseContext) String() string {
	b, _ := json.MarshalIndent(a, " ", " ")
	return string(b)
}
