package models

import "fmt"

type Seat struct {
	Row       int    `json:"row"`
	Column    string `json:"column"`
	Available bool   `json:"available"`
}

func (seat Seat) String() string {
	return fmt.Sprintf("%d %s", seat.Row, seat.Column)
}
