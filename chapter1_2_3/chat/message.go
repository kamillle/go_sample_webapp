package main

import (
	"time"
)

// websocketで送る情報を拡張するためのstruct
type message struct {
	Name    string
	Message string
	When    time.Time
}
