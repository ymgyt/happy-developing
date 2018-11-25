package app

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
}

// Mode -
type Mode int
