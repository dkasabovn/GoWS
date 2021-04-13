package ws

import (
	"encoding/json"
	"log"
)

// Generic Messages
const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const RoomJoinedAction = "room-joined"

// Game Messages
const ReadyUp = "read-up"
const StartGame = "start-game"
const EndGame = "end-game"
const NextQuestion = "next-question"
const QuestionSubmitted = "question-submitted"
const SendAnswer = "send-answer"

type Message struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
	Sender *Client                `json:"user"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	message.Sender = &msg.Sender
	return nil
}
