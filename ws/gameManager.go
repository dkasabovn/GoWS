package ws

import (
	"log"
	"time"
)

type GameManager struct {
	repo         *QuestionRepo
	timer        *time.Timer
	room         *Room
	players      int
	questionTime int
	readyUpTime  int
}

func (gme *GameManager) readyUpStage() {
	// TODO Currently skips if one player readies up; Probably should be both
	gme.timer = time.NewTimer(time.Duration(gme.readyUpTime) * time.Second)
	readiedUp := map[string]bool{}
	for {
		select {
		case <-gme.timer.C:
			return
		case command := <-gme.room.server.internal:
			if command.Action == ReadyUp {
				gme.room.BroadcastMessage(createMessage(ReadyUp, map[string]interface{}{
					"user": command.Sender.Name,
				}))
				readiedUp[command.Sender.Name] = true
			}
			if len(readiedUp) == gme.room.Active() {
				gme.players = len(readiedUp)
				gme.room.BroadcastMessage(createMessage(GameStarting, map[string]interface{}{
					"time": 5,
				}))
				return
			}
		}
	}
}

func (gme *GameManager) playGameStage() {
	gme.timer = time.NewTimer(5 * time.Second)
	startedFlag := false
	for {
		select {
		case msg := <-gme.room.server.internal:
			log.Println(msg)
			err := gme.repo.validate(msg)
			if err != nil {
				log.Println(err)
				continue
			}

			if gme.players == gme.repo.questionsSubmitted() {
				gme.timer = time.NewTimer(0 * time.Second)
			}
		case <-gme.timer.C:
			if startedFlag {
				sendAnswerPacket := map[string]interface{}{
					"results":     gme.repo.getResults(),
					"leaderboard": gme.repo.getLeaderboard(),
				}
				gme.room.BroadcastMessage(createMessage(SendAnswer, sendAnswerPacket))
				time.Sleep(5 * time.Second)
			} else {
				startedFlag = true
			}
			q := gme.repo.nextQuestion()
			if q == nil {
				return
			}
			gme.sendQuestion(q)
			gme.timer.Reset(time.Duration(gme.questionTime) * time.Second)
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
	payload["time"] = gme.questionTime
	nextQuestionMsg := &Message{
		Action: NextQuestion,
		Data:   payload,
	}
	gme.room.BroadcastMessage(nextQuestionMsg)
}

func (gme *GameManager) endGameStage() {
	endGameMsg := &Message{
		Action: EndGame,
		Data: map[string]interface{}{
			"message": "game is over! Room blowing up in 10 seconds ~OwO~",
		},
	}
	gme.room.BroadcastMessage(endGameMsg)
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

func NewGameManager(room *Room) *GameManager {
	qr := NewQR()
	qr.LoadRepo()
	return &GameManager{
		repo:         qr,
		room:         room,
		questionTime: 30,
		readyUpTime:  120,
	}
}

func createMessage(action string, data map[string]interface{}) *Message {
	return &Message{
		Action: action,
		Data:   data,
	}
}

func createSimpleMessage(action string, data string) *Message {
	return createMessage(action, map[string]interface{}{
		"message": data,
	})
}
