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
