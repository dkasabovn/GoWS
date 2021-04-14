package game

import "regexp"

const (
	MP = iota
	FR
)

type Question interface {
	validate(userResponse int) bool
}

type MultipleChoice struct {
	qType    int
	answer   int
	question string
}

type FreeResponse struct {
	qType    int
	answer   regexp.Regexp
	question string
}

func (fr *FreeResponse) validate(userResponse string) bool {
	return fr.answer.Match([]byte(userResponse))
}

func (mp *MultipleChoice) validate(userResponse int) bool {
	return userResponse == mp.answer
}
