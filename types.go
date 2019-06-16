package main

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
	Participants []Participant `json:"participants"`
	Messages []Message `json:"messages"`
	Title string `json:"title"`
	IsStillParticipant bool `json:"is_still_participant"`
	ThreadType ThreadType `json:"thread_type"`
	ThreadPath string `json:thread_path"`
}

type Participant struct {
	Name string `json:"name"`
}

type Message struct {
	SenderName string `json:"sender_name"`
	Timestamp int `json:"timestamp_ms"`
	Content string `json:"content"`
	Photos []Photo `json:"photoes"`
	Reactions []Reaction `json:"reactions"`
	CallDuration int `json:"content"`
	MessageType MessageType `json:"type"`
	Missed bool `json:"missed"`
}

type Photo struct {
	URI string `json:"uri"`
	CreationTimestamp int `json:"creation_timestamp"`
}

type Reaction struct {
	Reaction string `json:"reaction"`
	Actor string `json:"actor"`
}
