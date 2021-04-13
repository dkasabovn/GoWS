package game

import (
	"main/ws"
	"time"
)

type GameManager struct {
	questions       []Question
	answers         map[string][]interface{}
	timer           *time.Timer
	currentQuestion int
	room            *ws.Room
}

func (gme *GameManager) readyUpStage() {
	// TODO Currently skips if one player readies up; Probably should be both
	gme.timer = time.NewTimer(30 * time.Second)
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
	for {
		select {
		case <-gme.timer.C:
			gme.currentQuestion++
			return
		case command := <-gme.room.Commands:
			if command.Action == ws.QuestionSubmitted {
				gme.answers[command.Sender.ID.String()][gme.currentQuestion] = command.Data["question"]
			}
		}
	}
}

func (gme *GameManager) endGameStage() {
	
}
