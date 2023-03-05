package models

import "errors"

var (
	ErrEmpty     = errors.New("is empty")
	ErrIncorrect = errors.New("is incorrect")
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

type Quotes []*Quote
