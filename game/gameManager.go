package game

import "time"

type GameManager struct {
	questions       []Question
	timer           *time.Timer
	currentQuestion int
	messages        chan []byte
}
