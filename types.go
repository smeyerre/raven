package main

type MessageType string
type ThreadType string

const (
	GENERIC MessageType = "Generic"
	SHARE   MessageType = "Share"
	CALL    MessageType = "Call"
)

const (
	REGULAR       ThreadType = "Regular"
	REGULAR_GROUP ThreadType = "RegularGroup"
)

type ConfigFile struct {
	Username            string       `json:"username"`
	MessageFileType     string       `json:"messageFileType"`
	ConvoDirectoryNames []string     `json:"convoDirectoryNames"`
}

type MessageFile struct {
	Participants       []Participant `json:"participants"`
	Messages           []Message     `json:"messages"`
	Title              string        `json:"title"`
	IsStillParticipant bool          `json:"is_still_participant"`
	ThreadType         ThreadType    `json:"thread_type"`
	ThreadPath         string        `json:"thread_path"`
}

type FlourishMessageFile struct {
	Participants []Participant `json:"participants"`
	Messages     []FlourishMessage     `json:"messages"`
	Title        string        `json:"title"`
	ThreadType   ThreadType    `json:"thread_type"`
}

type Participant struct {
	Name string `json:"name"`
}

type Message struct {
	SenderName   string      `json:"sender_name"`
	Timestamp    int64       `json:"timestamp_ms"`
	Content      string      `json:"content"`
	Photos       []Photo     `json:"photos"`
	Reactions    []Reaction  `json:"reactions"`
	CallDuration int64       `json:"content"`
	MessageType  MessageType `json:"type"`
	Missed       bool        `json:"missed"`
}

type FlourishMessage struct {
	SenderName string `json:"sender_name"`
	Timestamp  int64  `json:"timestamp_ms"`
}

type Photo struct {
	URI               string `json:"uri"`
	CreationTimestamp int64  `json:"creation_timestamp"`
}

type Reaction struct {
	Reaction string `json:"reaction"`
	Actor    string `json:"actor"`
}
