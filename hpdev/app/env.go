package app

import (
	"context"
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	Mode      Mode
	Log       *zap.Logger
	Ctx       context.Context
	Now       func() time.Time
	Validator Validator
}

// Mode -
type Mode int

func (m Mode) String() (s string) {
	switch m {
	case DevelopmentMode:
		s = "development"
	case TestingMode:
		s = "testing"
	case ProductionMode:
		s = "production"
	default:
		s = "undefined"
	}
	return s
}

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

const (
	encodeConsole = "console"
	encodeJSON    = "json"
)

// LoggingConfig -
type LoggingConfig struct {
	Out  io.Writer
	Mode Mode
}

// MustLogger -
func MustLogger(cfg *LoggingConfig) *zap.Logger {
	z, err := NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	return z
}

// NewLogger -
func NewLogger(cfg *LoggingConfig) (*zap.Logger, error) {

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "l",
		CallerKey:      "c",
		MessageKey:     "m",
		StacktraceKey:  "s",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zap.AtomicLevel
	var encode string
	switch cfg.Mode {
	case DevelopmentMode:
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
		encode = encodeConsole
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
		encode = encodeJSON
	}

	var encoder zapcore.Encoder
	if encode == encodeConsole {
		encoder = zapcore.NewConsoleEncoder(encCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(cfg.Out), level)
	z := zap.New(core, zap.AddCaller())

	return z, nil
}
