package generate

//go:generate stringer -type=Pill
// ^ go:generate [cmd] [args]

// run...
// $ go generate tool/generate/pill.go

// Pill type
type Pill int

// Types of Pill
const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)
