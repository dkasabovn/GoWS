package game

import (
	"log"
	"main/config"

	"google.golang.org/api/iterator"
)

type QuestionRepo struct {
	questions       []Question
	sumbissions     map[string][]bool
	currentQuestion int
}

func NewQR() *QuestionRepo {
	return &QuestionRepo{
		questions:       make([]Question, 0),
		currentQuestion: 0,
	}
}

func (qr *QuestionRepo) nextQuestion() *Question {
	if len(qr.questions) <= qr.currentQuestion {
		return nil
	} else {
		q := &qr.questions[qr.currentQuestion]
		qr.currentQuestion++
		return q
	}
}

func (qr *QuestionRepo) uniqueAnswers() int {
	i := 0
	for k := range qr.sumbissions {
		if len(qr.sumbissions[k]) == qr.currentQuestion {
			i++
		}
	}
	return i
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
