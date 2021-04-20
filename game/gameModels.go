package game

import (
	"encoding/json"
	"regexp"
)

const (
	MP = iota
	FR
)

type Question interface {
	questionType() int
	validate(userResponse interface{}) bool
}

type MultipleChoice struct {
	QType     int      `firestore:"qType" json:"qType"`
	Answer    int      `firestore:"answer" json:"answer"`
	SQuestion string   `firestore:"question" json:"question"`
	Qid       int      `firestore:"qid" json:"qid"`
	Choices   []string `firestore:"choices" json:"choices"`
}

type FreeResponse struct {
	qType    int
	answer   regexp.Regexp
	question string
}

func (fr *FreeResponse) validate(userResponse interface{}) bool {
	if val, ok := userResponse.(string); ok {
		return fr.answer.Match([]byte(val))
	}
	return false
}

func (mp *MultipleChoice) validate(userResponse interface{}) bool {
	if val, ok := userResponse.(int); ok {
		return val == mp.Answer
	}
	return false
}

func (fr *FreeResponse) questionType() int {
	return fr.qType
}

func (mp *MultipleChoice) questionType() int {
	return mp.QType
}

func structToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}