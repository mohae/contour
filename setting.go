package contour

// setting holds the information for a configuration setting.
type setting struct {
	// Type is the setting's datatype
	Type dataType
	// Name
	Name string
	// The short code for the setting, if applicable
	Short string
	// The current value of the setting
	Value interface{}
	// Usage is the usage information for this setting
	Usage string
	// Default is the string version of the default, for usage
	Default string
	// IsCore: a core settings is for settings whose values that cannot be
	// be changed after they are registered.  When this is true, IsConfFileVar,
	// IsEnv, and IsFlag are always false.
	IsCore bool
	// IsConfFileVar: a configuration file setting can only be updated from a
	// configuration file. When this is true, IsEnv and IsFlag aare always false.
	IsConfFileVar bool
	// IsEnvVar: whether or not this is an environment variable.  When true, and
	// the cfg is set to use EnvsVars, the setting will be settable via
	// environment variables. All EnvVars and Flag settings result in IsEnv being
	// true.
	IsEnvVar bool
	// IsFlag:  whether or not this is a flag. When true, IsCfg and IsEnvVar will
	// also be true.
	IsFlag bool
	// Alias
	Alias []string
}
