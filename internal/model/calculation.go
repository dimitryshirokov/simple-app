package model

import "time"

type Calculation struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	A         int       `json:"a"`
	B         int       `json:"b"`
	Result    int       `json:"result"`
	Type      string    `json:"type"`
}
