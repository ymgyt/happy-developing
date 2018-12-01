package app

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// UndefinedMode -
	UndefinedMode Mode = iota
	// DevelopmentMode -
	DevelopmentMode
	// TestingMode -
	TestingMode
	// ProductionMode -
	ProductionMode
)

// Env -
type Env struct {
	Log *logrus.Logger
	Ctx context.Context
	Now func() time.Time
}

// Mode -
type Mode int

// TimeZone -
var TimeZone = time.FixedZone("Asia/Tokyo", 9*60*60)

// Now -
func Now() time.Time {
	return time.Now().In(TimeZone)
}

// ErrCode -
type ErrCode int

const (
	// OK -
	OK ErrCode = iota
)
