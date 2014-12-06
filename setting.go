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

	// IsConfig:
	IsCfg bool

	// IsFlag:  whether or not this is a flag. When IsFlag is true,
	// IsConfig will also be true since a Flag is a subset of Config.
	IsFlag bool

	// Alias
	Alias []string
}
