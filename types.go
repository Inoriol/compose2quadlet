package compose2quadlet

// WarningLevel classifies how a field mapping issue is surfaced to the consumer.
type WarningLevel int

const (
	WarningSkipped  WarningLevel = iota // feature unavailable at target podman version, skipped with info
	WarningDegraded                     // P3 PodmanArgs fallback used instead of P1 native directive
	WarningFatal                        // cannot proceed (mapping is impossible at this podman version)
)

// Warning describes a single mapping concern for a specific compose field and service.
type Warning struct {
	Level    WarningLevel
	Service  string // service name, empty if not service-scoped
	Field    string // compose field name (e.g. "build", "entrypoint")
	Message  string // human-readable description
	Since    string // minimum podman version required, if applicable (e.g. "5.2.0")
}

// UnitType represents the kind of quadlet unit file.
type UnitType string

const (
	UnitContainer UnitType = "container"
	UnitNetwork   UnitType = "network"
	UnitVolume    UnitType = "volume"
	UnitImage     UnitType = "image"
	UnitBuild     UnitType = "build"
)

// QuadletUnit is a structured representation of a quadlet unit file
// before serialization to ini-format.
type QuadletUnit struct {
	Type     UnitType
	Name     string
	Sections []Section
}

// Section represents a single ini section in a quadlet unit file
// (e.g. [Container], [Service], [Unit], [Install]).
type Section struct {
	Name       string
	Directives []Directive
}

// Directive is a single key-value entry in a quadlet section.
// Values holds the ordered list for directives that accept multiple values
// on separate lines with the same key.
type Directive struct {
	Key    string
	Values []string
}

// Section names used across quadlet unit types.
const (
	SectionUnit      = "Unit"
	SectionService   = "Service"
	SectionInstall   = "Install"
	SectionContainer = "Container"
	SectionNetwork   = "Network"
	SectionVolume    = "Volume"
	SectionImage     = "Image"
	SectionBuild     = "Build"
)
