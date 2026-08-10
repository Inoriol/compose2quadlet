package compose2quadlet

type transpileConfig struct {
	// Opinionated transforms — all enabled by default.
	selinuxContext  bool
	filePrefix      string
	defaultNetwork  bool
	portOffset      int
	projectName     string
	managedLabel    bool
	autoUpdate      bool
	installSection  bool

	// Network alias injection.
	networkAliases bool

	// Podman version to target. Zero value means "latest" (emit best-available mapping everywhere).
	podmanVersion Version

	// Warnings collected during the transpilation pipeline. Read after Transpile() returns.
	Warnings []Warning
}

// Version holds a semantic version triple. Zero value means "latest / unspecified".
type Version struct {
	Major, Minor, Patch int
}

// AtLeast returns true if v is at least the given major.minor.
// Zero Version always returns true (latest implies everything is available).
func (v Version) AtLeast(major, minor int) bool {
	if v.Major == 0 && v.Minor == 0 {
		return true
	}
	if v.Major > major {
		return true
	}
	if v.Major == major && v.Minor >= minor {
		return true
	}
	return false
}

func defaultConfig() *transpileConfig {
	return &transpileConfig{
		selinuxContext:  true,
		filePrefix:      "cq-",
		defaultNetwork:  true,
		portOffset:      0,
		managedLabel:    true,
		autoUpdate:      false,
		installSection:  true,
		networkAliases:  true,
	}
}

// TranspileOption controls opinionated transforms applied during transpilation.
type TranspileOption func(*transpileConfig)

// WithoutSELinux disables adding :z SELinux context to volume mounts.
func WithoutSELinux() TranspileOption {
	return func(c *transpileConfig) { c.selinuxContext = false }
}

// WithoutPrefix removes the cq-<project>- file name prefix.
func WithoutPrefix() TranspileOption {
	return func(c *transpileConfig) { c.filePrefix = "" }
}

// WithPrefix sets a custom file name prefix (default: "cq-<project>-").
func WithPrefix(prefix string) TranspileOption {
	return func(c *transpileConfig) { c.filePrefix = prefix }
}

// WithoutDefaultNetwork disables automatic injection of the default network.
func WithoutDefaultNetwork() TranspileOption {
	return func(c *transpileConfig) { c.defaultNetwork = false }
}

// WithPortOffset applies a rootless port offset (e.g. 8080 → 8080+offset).
func WithPortOffset(offset int) TranspileOption {
	return func(c *transpileConfig) { c.portOffset = offset }
}

// WithProjectName sets the project name used for labels and prefix generation.
func WithProjectName(name string) TranspileOption {
	return func(c *transpileConfig) { c.projectName = name }
}

// WithoutManagedLabel disables the com.comquad.managed=true label.
func WithoutManagedLabel() TranspileOption {
	return func(c *transpileConfig) { c.managedLabel = false }
}

// WithAutoUpdate enables AutoUpdate=registry on container units.
func WithAutoUpdate() TranspileOption {
	return func(c *transpileConfig) { c.autoUpdate = true }
}

// WithoutInstallSection disables automatic [Install] section injection.
func WithoutInstallSection() TranspileOption {
	return func(c *transpileConfig) { c.installSection = false }
}

// WithoutNetworkAliases disables automatic NetworkAlias=<service> injection.
func WithoutNetworkAliases() TranspileOption {
	return func(c *transpileConfig) { c.networkAliases = false }
}

// WithPodmanVersion sets the target podman version. Mappers use this to decide
// between P1 native directives and P3 PodmanArgs fallbacks. Fields that are
// structurally impossible at the given version (e.g. build on < 5.2.0) produce
// a fatal warning and Transpile() returns an error.
//
// Pass a zero Version or omit the option to target the latest podman release.
func WithPodmanVersion(v Version) TranspileOption {
	return func(c *transpileConfig) { c.podmanVersion = v }
}

// warn appends a warning to the config. Used internally by mappers to surface
// skipped, degraded, and fatal mapping outcomes.
func (c *transpileConfig) warn(w Warning) {
	c.Warnings = append(c.Warnings, w)
}
