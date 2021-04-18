package game

import (
	"log"
	"main/ws"
	"time"
)

type GameManager struct {
	repo  *QuestionRepo
	timer *time.Timer
	room  *ws.Room
}

func (gme *GameManager) readyUpStage() {
	// TODO Currently skips if one player readies up; Probably should be both
	gme.timer = time.NewTimer(120 * time.Second)
	for {
		select {
		case <-gme.timer.C:
			return
		case command := <-gme.room.Commands:
			if command.Action == ws.ReadyUp {
				gme.room.Broadcast <- createSimpleMessage(ws.StartGame, "Game Starting you mongoloid. Calm down. Take a breath.")
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
			q := gme.repo.nextQuestion()
			if q == nil {
				return
			}
			gme.sendQuestion(q)
			gme.timer.Reset(30 * time.Second)
			break
		}
	}
}

func (gme *GameManager) sendQuestion(q *Question) {
	payload, err := structToMap(q)
	if err != nil {
		log.Println("Error")
		return
	}
	nextQuestionMsg := &ws.Message{
		Action: ws.NextQuestion,
		Data:   payload,
	}
	gme.room.Broadcast <- nextQuestionMsg
}

func (gme *GameManager) endGameStage() {
	endGameMsg := &ws.Message{
		Action: ws.EndGame,
		Data: map[string]interface{}{
			"message": "game is over! Room blowing up in 10 seconds ~OwO~",
		},
	}
	gme.room.Broadcast <- endGameMsg
	gme.timer = time.NewTimer(10 * time.Second)
	for {
		select {
		case <-gme.timer.C:
			return
		}
	}
}

func (gme *GameManager) sendQuestionIntervalTest() {
	interval := time.NewTicker(time.Second * 3)
	q := gme.repo.nextQuestion()
	for {
		select {
		case <-interval.C:
			gme.sendQuestion(q)
		}
	}
}

func (gme *GameManager) Run() {
	gme.readyUpStage()
	log.Println("Finished Ready Up Stage")
	gme.playGameStage()
	log.Println("Finished Playing the Game")
	gme.endGameStage()
	log.Println("Game wrapped up; Destroying myself ~OwO~")
	gme.room.Terminate()
}

func NewGameManager(room *ws.Room) *GameManager {
	qr := NewQR()
	qr.LoadRepo()
	return &GameManager{
		repo: qr,
		room: room,
	}
}

func createMessage(action string, data map[string]interface{}) *ws.Message {
	return &ws.Message{
		Action: action,
		Data:   data,
	}
}

func createSimpleMessage(action string, data string) *ws.Message {
	return createMessage(action, map[string]interface{}{
		"message": data,
	})
}
