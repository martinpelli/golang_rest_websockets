package models

type Websocketmessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
