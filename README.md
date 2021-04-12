# GoWS

For now this library is made for managing multiple Rooms. Users join a room / 'game' and when they leave that room they get disconnected from the WS. However, adding functionality for switching rooms is half baked in.

### Features currently supported:
- Redis Pub/Sub
- Reading / Writing to a room

### Features that can be easily implemented:
- Switching rooms

### Planned:
- Adding Supabase for auth / user storage
- Making stuff more organized

# Docs for big boy Neo:

## HTTP

`/create`  
Creates a new room. Returns a JSON objet containing the room's unique ID
```js
{
  socket: "1234-1231231-23123123"
}
```

## WS

`/ws?name={NAME}&room={ROOM_ID}`  
Connects to the proper websocket room. Room must already be created (using `/create`) before joining.

### Message Object
```
{
  Action: string,
  Message: string | map[string]interface{},
  Sender: Client
}
```

##### Actions
```
// Generic Messages
SendMessageAction = "send-message"
JoinRoomAction = "join-room"
LeaveRoomAction = "leave-room"
UserJoinedAction = "user-join"
UserLeftAction = "user-left"
RoomJoinedAction = "room-joined"

// Game Messages
StartGame = "start-game"
EndGame = "end-game"
NextQuestion = "next-question"
QuestionSubmitted = "question-submitted"
SendAnswer = "send-answer"
```
  
