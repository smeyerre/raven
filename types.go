package main

import {
	"encoding/json"
}

type MessageType int
type ThreadType int

const (
	GENERIC MessageType = iota
	SHARE MessageType = iota
	CALL MessageType = iota
)

const (
	REGULAR ThreadType = iota
	REGULAR_GROUP ThreadType = iota
)

type MessageFile struct {
	participants []string `json:"participants"`
	messages []Message `json:"messages"`
	title string `json:"title"`
	isStillParticipant bool `json:"is_still_participant"`
	threadType ThreadType `json:"thread_type"`
	threadPath string `json:thread_path"`
}

type Message struct {
	senderName string `json:"sender_name"`
	timestamp int `json:"timestamp_ms"`
	content string `json:"content"`
	photos []Photo `json:"photoes"`
	reactions []Reaction `json:"reactions"`
	callDuration int `json:"content"`
	messageType MessageType `json:"type"`
	missed bool `json:"missed"`
}

type Photo struct {
	uri string `json:"uri"`
	creationTimestamp int `json:"creation_timestamp"`
}

type Reaction struct {
	reaction string `json:"reaction"`
	actor string `json:"actor"`
}
