package s

import "time"

type Settings struct {
	Id        string
	Timestamp time.Time
	Int       int64
	String    string
	Bool      bool
}
