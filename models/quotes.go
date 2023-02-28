package models

import "errors"

var (
	ErrEmpty = errors.New("is empty")
)

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

type Quotes []*Quote
