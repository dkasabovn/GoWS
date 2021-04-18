package game

import (
	"main/ws"
	"time"
)

type GameManager struct {
	questions       []MultipleChoice
	answers         map[string][]interface{}
	timer           *time.Timer
	currentQuestion int
	room            *ws.Room
}

func (gme *GameManager) readyUpStage() {
	startMsg := &ws.Message{
		Action: ws.SendMessageAction,
		Data: map[string]interface{}{
			"message": "Game started",
		},
	}
	gme.room.Broadcast <- startMsg
	// TODO Currently skips if one player readies up; Probably should be both
	gme.timer = time.NewTimer(1 * time.Second)
	for {
		select {
		case <-gme.timer.C:
			return
		case command := <-gme.room.Commands:
			if command.Action == ws.ReadyUp {
				return
			}
		}
	}
}

func (gme *GameManager) playGameStage() {
	// TODO do a better job of tracking if all users have answered; Skip if all answered
	gme.timer = time.NewTimer(1 * time.Second)
	for {
		select {
		case <-gme.timer.C:
			gme.sendQuestion()
			gme.currentQuestion += 1
			if gme.currentQuestion == len(gme.questions) {
				return
			}
			gme.timer.Reset(1 * time.Second)
			break
		}
	}
}

func (gme *GameManager) sendQuestion() {
	nextQuestionMsg := &ws.Message{
		Action: ws.NextQuestion,
		Data: map[string]interface{}{
			"message": "You really thought I connected the db. LOL!",
		},
	}
	gme.room.Broadcast <- nextQuestionMsg
}

func (gme *GameManager) endGameStage() {
	endGameMsg := &ws.Message{
		Action: ws.EndGame,
		Data: map[string]interface{}{
			"message": "game is over!",
		},
	}
	gme.room.Broadcast <- endGameMsg
}

func (gme *GameManager) Run() {
	// gme.readyUpStage()
	// log.Println("Finished Ready Up Stage")
	// gme.playGameStage()
	// log.Println("Finished Playing the Game")
	// gme.endGameStage()
	// log.Println("Game wrapped up; Destroying myself ~OwO~")
	gme.room.Terminate()
}

func NewGameManager(room *ws.Room) *GameManager {
	return &GameManager{
		questions:       []MultipleChoice{{qType: 1, answer: 1, question: "asdfasdf"}},
		currentQuestion: 0,
		room:            room,
	}
}
