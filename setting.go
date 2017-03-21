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
	// IsCore: whether or not this is considered a core setting. Core settings if for
	// settings whose values you don't want overridden or changed, once registered.
	//
	// When IsCore is true, IsConfig and IsFlag are always false. These cannot be changed either.
	IsCore bool
	// IsConfig: whether or not this if a Cfg setting. When true, and a cfg file exists,
	// it will check for this setting in the config file.
	IsCfg bool
	// IsEnv: whether or not this is an env variable.  When true, and the cfg is set to
	// useEnvs, the setting will be settable via env variables. All Cfg and Flag settings
	// result in IsEnv being true.
	IsEnv bool
	// IsFlag:  whether or not this is a flag. When true, IsCfg and IsEnv will also be true.
	IsFlag bool
	// Alias
	Alias []string
}
