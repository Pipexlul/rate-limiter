package models

import "time"

type ClientData struct {
	Tokens      int
	LastRequest time.Time
}
