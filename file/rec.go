package file

import (
	"time"
)

// RecShort - short file record
type RecShort struct {
	Filename string
	TStamp   time.Time
}

// RecFull - short file record
type RecFull struct {
	RecShort
	Name         string
	Requirements []string
	SQLData      *string
}
