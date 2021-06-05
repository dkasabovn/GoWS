package ws

import (
	"errors"
	"log"
	"main/config"

	"google.golang.org/api/iterator"
)

type QuestionRepo struct {
	questions       []Question
	submissions     []map[string]interface{}
	leaderboard     map[string]int
	currentQuestion int
}

func NewQR() *QuestionRepo {
	return &QuestionRepo{
		questions:       make([]Question, 0),
		submissions:     make([]map[string]interface{}, 0),
		leaderboard:     make(map[string]int),
		currentQuestion: 0,
	}
}

func (qr *QuestionRepo) nextQuestion() *Question {
	if len(qr.questions) <= qr.currentQuestion {
		return nil
	} else {
		q := &qr.questions[qr.currentQuestion]
		qr.currentQuestion++
		qr.submissions = append(qr.submissions, map[string]interface{}{})
		return q
	}
}

func (qr *QuestionRepo) getQuestion() Question {
	if len(qr.questions) <= qr.currentQuestion {
		return nil
	}
	return qr.questions[qr.currentQuestion]
}

func (qr *QuestionRepo) questionsSubmitted() int {
	return len(qr.submissions[qr.currentQuestion-1])
}

func (qr *QuestionRepo) validate(m *Message) error {
	if m.Action != QuestionSubmitted {
		return errors.New("Wrong action")
	}
	cq := qr.questions[qr.currentQuestion-1]
	if questionID, ok := m.Data["qid"]; ok && int(questionID.(float64)) != cq.questionID() {
		return errors.New("Incorrect Question was Provided")
	}
	if answer, ok := m.Data["answer"]; ok {
		isCorrect := cq.validate(answer)
		qr.submissions[qr.currentQuestion-1][m.Sender.Name] = isCorrect
		if currentValue, ok := qr.leaderboard[m.Sender.Name]; ok {
			qr.leaderboard[m.Sender.Name] = currentValue + 1
		} else {
			if isCorrect {
				qr.leaderboard[m.Sender.Name] = 1
			} else {
				qr.leaderboard[m.Sender.Name] = 0
			}
		}
		log.Println(qr.submissions[qr.currentQuestion-1])
		return nil
	}
	return errors.New("No answer provided")
}

func (qr *QuestionRepo) getResults() map[string]interface{} {
	return qr.submissions[qr.currentQuestion-1]
}

func (qr *QuestionRepo) getLeaderboard() map[string]int {
	return qr.leaderboard
}

// TODO actually get random question IDs
func (qr *QuestionRepo) getRandomQuestionIDs() interface{} {
	return []int{1}
}

func (qr *QuestionRepo) LoadRepo() {
	docs := config.FsClient.Collection("questions").Where("qid", "in", qr.getRandomQuestionIDs()).Documents(config.CTX)
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Couldn't iterate through all questions")
		}
		var mp MultipleChoice
		if err := doc.DataTo(&mp); err != nil {
			log.Println("Question cannot be parsed as Multiple Choice")
		} else {
			log.Println(mp)
			qr.questions = append(qr.questions, &mp)
			continue
		}
		var fr FreeResponse
		if err := doc.DataTo(&fr); err != nil {
			log.Println("Question cannot be parsed as FreeResponse")
		} else {
			qr.questions = append(qr.questions, &fr)
			continue
		}
	}
}
