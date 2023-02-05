package domain

import "time"

type WebObj struct {
	Theme string    `json:"theme"`
	Link  string    `json:"link"`
	Site  string    `json:"site"`
	Time  time.Time `json:"time"`
}
