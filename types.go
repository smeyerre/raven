package main

import {
	"encoding/json"
}

type MessageType int

const (
	GENERIC MessageType = iota
	SHARE MessageType = iota
	CALL MessageType = iota
)

type Message struct {
	senderName string `json:"sender_name"`
	timestamp int `json:"timestamp_ms"`
	content string `json:"content"`
	photos []Photo `json:"photoes"`
	callDuration int `json:"content"`
	messageType MessageType `json:"type"`
	missed bool `json:"missed"`
}

type Photo struct {
	uri string `json:"uri"`
	creationTimestamp int `json:"creation_timestamp"`
}
