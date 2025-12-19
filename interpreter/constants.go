package interpreter

// Random generation constants
const (
	// RandomBoolChoices represents the number of choices for boolean random generation (true/false)
	RandomBoolChoices = 2

	// RangeInclusiveOffset is added to make range bounds inclusive
	RangeInclusiveOffset = 1
)

// Unofloat bounds
const (
	UnofloatMin = 0.0
	UnofloatMax = 1.0
)

// Default probability for ifrand without explicit probability
const (
	DefaultIfrandProbability = 0.5
)
