package db

import (
	"time"
)

// RecShort - short record for info
type RecShort struct {
	ID           uint64
	Name         string
	Applied      bool
	Note         *string
	GITinfo      *string
	Filename     string
	TStamp       time.Time
	Requirements []string
}

// RecFull - full record
type RecFull struct {
	RecShort
	SQLData *string
}

// RecHistory - history record
type RecHistory struct {
	ID       uint64
	When     time.Time
	Name     string
	Applied  bool
	Note     *string
	GITinfo  *string
	Filename string
}
