package contour

// setting holds the information for a configuration setting.
type setting struct {
	// Type is the datatype for the setting
	Type string
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
	// IsCore: whether or not this is considered a core setting. Core
	// settings if for things like application name, where you don't want
	// anything else overwriting that value, once set. If a setting IsCore
	// it cannot be over-written and the IsCore value cannot be changed.
	//
	// When IsCore is true, IsConfig and IsFlag are always false. These
	// cannot be changed either.
	IsCore bool
	// IsConfig: whether or not this if a Cfg setting. When true, and a config
	// file exists, it will check for this setting in the config file.
	IsCfg bool
	// IsEnv: whether or not this is a Env setting. When true, and the EnvName
	// is != "", it will check for this settings in the environment. When IsCfg
	// is true, IsCfg will also be true as IsEnv is a subset of Cfg.
	IsEnv bool
	// IsFlag:  whether or not this is a flag. When IsFlag is true, IsCfg will
	// also be true since a Flag is a subset of Cfg.
	IsFlag bool
	// Alias
	Alias []string
}
