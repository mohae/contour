package contour

import "strconv"

// RegisterSetting registers a setting. For most settings, the data and setting
// type specific registration and add functions should be used. The exception
// would be when more granular control over what can update a registered
// setting is needed. This method allows Is[ConfFileVar|EnvVar|Flag] bools to
// be set independently.
//
// If a setting with the key k already exists, a SettingExistsErr will be
// returned. If k is an empty string an ErrNoSettingName will be returned.
//
// The short, dflt, and usage parms only apply to settings whose IsFlag bool
// is true.
//
// For non-Core settings, IsCore must be false. If IsCore is true, k's value
// cannot be changed after registration, regardless of the truthiness of
// Is[ConfFileVar|IsEnvVar|IsFlag]. For Core settings, AddCore methods should
// be used.
//
// If the setting can be updated by a configuration file, environment variable
// or a flag, the Is[ConfFileVar|IsEnv|IsFlag] bools should be set to true as
// appropriate. These conditionals are independent; e.g. a setting can have
// both IsConfFileVar and IsFlag set to true if the setting is not to be
// updateable from an environment variable.
//
// If Is[Core|ConfFileVar|Env|Flag] are all false, the setting will only be
// updateable by using the Update methods. For these kind of settings, the
// usage of Add functions should be preferred. These setting will not be
// exposed to the configuration file, as an environment variable, or as a flag.
//
// For non string, bool, int, and int64 types, the type must be "interface{}"
func (s *Settings) RegisterSetting(typ, name, short string, value interface{}, dflt, usage string, IsCore, IsConfFileVar, IsEnvVar, IsFlag bool) error {
	dType := parseDataType(typ)
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerSetting(0, dType, name, short, value, dflt, usage, IsCore, IsConfFileVar, IsEnvVar, IsFlag)
}

func (s *Settings) registerSetting(sTyp SettingType, typ dataType, name, short string, value interface{}, dflt, usage string, IsCore, IsConfFileVar, IsEnvVar, IsFlag bool) error {
	if name == "" {
		return ErrNoSettingName
	}
	_, ok := s.settings[name]
	if ok {
		// Settings can't be re-registered.
		return SettingExistsErr{typ: sTyp, k: name}
	}

	// Add the setting
	s.settings[name] = setting{
		Type:          typ,
		Name:          name,
		Short:         short,
		Value:         value,
		Default:       dflt,
		Usage:         usage,
		IsCore:        IsCore,
		IsConfFileVar: IsConfFileVar,
		IsEnvVar:      IsEnvVar,
		IsFlag:        IsFlag,
	}
	// if it's a conf file setting, add it to the confFileVars map
	if IsConfFileVar {
		s.confFileVars[name] = struct{}{}
	}
	// mapping shortcodes make lookup easier
	if short != "" && IsFlag {
		v, ok := s.shortFlags[short]
		if ok {
			return ShortFlagExistsErr{k: name, short: short, shortName: v}
		}
		s.shortFlags[short] = name
	}
	if IsEnvVar {
		s.useEnvVars = IsEnvVar
	}
	// If a setting is a confFile setting, enable using a conf file.
	if IsConfFileVar {
		s.useConfFile = true
	}
	if IsFlag {
		s.useFlags = IsFlag
	}
	return nil
}

// RegisterBoolConfFileVar registers a bool setting using k for its key and v
// for its value. Once registered, the value of this setting can only be
// updated from a configuration file. If k already exists a SettingExistsErr
// will be returned. If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) RegisterBoolConfFileVar(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBoolConfFileVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBoolConfFileVar(k string, v bool) error {
	return s.registerConfFileVar(_bool, k, v, strconv.FormatBool(v))
}

// RegisterIntConfFileVar registers an int setting using k for its key and v
// for its value. Once registered, the value of this setting can only be
// updated from a configuration file. If k already exists a SettingExistsErr
// will be returned. If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) RegisterIntConfFileVar(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntConfFileVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerIntConfFileVar(k string, v int) error {
	return s.registerConfFileVar(_int, k, v, strconv.Itoa(v))
}

// RegisterInt64ConfFileVar registers an int64 setting using k for its key and
// v for its value. Once registered, the value of this setting can only be
// updated from a configuration file. If k already exists a SettingExistsErr
// will be returned. If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) RegisterInt64ConfFileVar(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64ConfFileVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64ConfFileVar(k string, v int64) error {
	return s.registerConfFileVar(_int64, k, v, strconv.FormatInt(v, 10))
}

// RegisterStringConfFileVar registers a string setting using k for its key and
// v for its value. Once registered, the value of this setting can only be
// updated from a configuration file. If k already exists a SettingExistsErr
// will be returned. If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) RegisterStringConfFileVar(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerStringConfFileVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerStringConfFileVar(k, v string) error {
	return s.registerConfFileVar(_string, k, v, v)
}

func (s *Settings) registerConfFileVar(typ dataType, k string, v interface{}, dflt string) error {
	return s.registerSetting(ConfFileVar, typ, k, "", v, dflt, "", false, true, false, false)
}

// RegisterBoolEnvVar registers a bool setting using k for its key and v
// for its value. Once registered, the value of this setting can only be
// updated from a configuration file or an environment variable. If k already
// exists a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned.
func (s *Settings) RegisterBoolEnvVar(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBoolEnvVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBoolEnvVar(k string, v bool) error {
	return s.registerEnvVar(_bool, k, v, strconv.FormatBool(v))
}

// RegisterIntEnvVar registers an int setting using k for its key and v
// for its value. Once registered, the value of this setting can only be
// updated from a configuration file or an environment variable. If k already
// exists a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned.
func (s *Settings) RegisterIntEnvVar(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntEnvVar(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) registerIntEnvVar(k string, v int) error {
	return s.registerEnvVar(_int, k, v, strconv.Itoa(v))
}

// RegisterInt64EnvVar registers an int64 setting using k for its key and v
// for its value. Once registered, the value of this setting can only be
// updated from a configuration file or an environment variable. If k already
// exists a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned.
func (s *Settings) RegisterInt64EnvVar(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64EnvVar(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) registerInt64EnvVar(k string, v int64) error {
	return s.registerEnvVar(_int64, k, v, strconv.FormatInt(v, 10))
}

// RegisterStringEnvVar registers a string setting using k for its key and v
// for its value. Once registered, the value of this setting can only be
// updated from a configuration file or an environment variable. If k already
// exists a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned.
func (s *Settings) RegisterStringEnvVar(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerStringEnvVar(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) registerStringEnvVar(k, v string) error {
	return s.registerEnvVar(_string, k, v, v)
}

func (s *Settings) registerEnvVar(typ dataType, k string, v interface{}, dflt string) error {
	return s.registerSetting(EnvVar, typ, k, "", v, dflt, "", false, true, true, false)
}

// RegisterBoolFlag registers a bool setting using k for its key and v
// for its value. Once registered, the value of this setting can be updated
// from a configuration file, an environment variable, or a flag. If k already
// exists a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned.
func (s *Settings) RegisterBoolFlag(k, short string, v bool, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBoolFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBoolFlag(k, short string, v bool, dflt, usage string) error {
	return s.registerFlag(_bool, k, short, v, dflt, usage)
}

// RegisterIntFlag registers an int setting using k for its key and v
// for its value. Once registered, the value of this setting can be updated
// from a configuration file, an environment variable, or a flag. If k already
// exists a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned
func (s *Settings) RegisterIntFlag(k, short string, v int, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerIntFlag(k, short string, v int, dflt, usage string) error {
	return s.registerFlag(_int, k, short, v, dflt, usage)
}

// RegisterInt64Flag registers an int64 setting using k for its key and v
// for its value. Once registered, the value of this setting can be updated
// from a configuration file, an environment variable, or a flag. If k already
// exists a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned
func (s *Settings) RegisterInt64Flag(k, short string, v int64, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64Flag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64Flag(k, short string, v int64, dflt, usage string) error {
	return s.registerFlag(_int64, k, short, v, dflt, usage)
}

// RegisterStringFlag registers a string setting using k for its key and v
// for its value. Once registered, the value of this setting can be updated
// from a configuration file, an environment variable, or a flag. If k already
// exists a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned
func (s *Settings) RegisterStringFlag(k, short, v, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerStringFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerStringFlag(k, short, v, dflt, usage string) error {
	return s.registerFlag(_string, k, short, v, dflt, usage)
}

func (s *Settings) registerFlag(typ dataType, k, short string, v interface{}, dflt, usage string) error {
	return s.registerSetting(Flag, typ, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterSetting registers a setting with the standard settings. For most
// settings, the data and setting type specific registration and add functions
// should be used. The exception would be when more granular control over what
// can update a registered setting is needed. This method allows
// Is[ConfFileVar|EnvVar|Flag] bools to be set independently.
//
// If a setting with the key k already exists, a SettingExistsErr will be
// returned. If k is an empty string an ErrNoSettingName will be returned.
//
// The short, dflt, and usage parms only apply to settings whose IsFlag bool
// is true.
//
// For non-Core settings, IsCore must be false. If IsCore is true, k's value
// cannot be changed after registration, regardless of the truthiness of
// Is[ConfFileVar|IsEnvVar|IsFlag]. For Core settings, AddCore methods should
// be used.
//
// If the setting can be updated by a configuration file, environment variable
// or a flag, the Is[ConfFileVar|IsEnv|IsFlag] bools should be set to true as
// appropriate. These conditionals are independent; e.g. a setting can have
// both IsConfFileVar and IsFlag set to true if the setting is not to be
// updateable from an environment variable.
//
// If Is[Core|ConfFileVar|Env|Flag] are all false, the setting will only be
// updateable by using the Update methods. For these kind of settings, the
// usage of Add functions should be preferred. These setting will not be
// exposed to the configuration file, as an environment variable, or as a flag.
//
// For non string, bool, int, and int64 types, the type must be "interface{}"
func RegisterSetting(typ, name, short string, value interface{}, dflt, usage string, IsCore, IsConfFileVar, IsEnv, IsFlag bool) error {
	return std.RegisterSetting(typ, name, short, value, dflt, usage, IsCore, IsConfFileVar, IsEnv, IsFlag)
}

// RegisterBoolConfFileVar registers a bool setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can only be updated from a configuration file. If k already exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func RegisterBoolConfFileVar(k string, v bool) error { return std.RegisterBoolConfFileVar(k, v) }

// RegisterIntConfFileVar registers an int setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can only be updated from a configuration file. If k already exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func RegisterIntConfFileVar(k string, v int) error { return std.RegisterIntConfFileVar(k, v) }

// RegisterInt64ConfFileVar registers an int64 setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can only be updated from a configuration file. If k already exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func RegisterInt64ConfFileVar(k string, v int64) error { return std.RegisterInt64ConfFileVar(k, v) }

// RegisterStringConfFileVar registers a string setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can only be updated from a configuration file. If k already exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func RegisterStringConfFileVar(k, v string) error { return std.RegisterStringConfFileVar(k, v) }

// RegisterBoolEnvVar registers a bool setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can only be updated from a configuration file or an environment
// variable. If k already exists a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned.
func RegisterBoolEnvVar(k string, v bool) error { return std.RegisterBoolEnvVar(k, v) }

// RegisterIntEnvVar registers an int setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can only be updated from a configuration file or an environment
// variable. If k already exists a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned.
func RegisterIntEnvVar(k string, v int) error { return std.RegisterIntEnvVar(k, v) }

// RegisterInt64EnvVar registers an int64 setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can only be updated from a configuration file or an environment
// variable. If k already exists a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned.
func RegisterInt64EnvVar(k string, v int64) error { return std.RegisterInt64EnvVar(k, v) }

// RegisterStringEnvVar registers a string setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can only be updated from a configuration file or an environment
// variable. If k already exists a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned.
func RegisterStringEnvVar(k, v string) error { return std.RegisterStringEnvVar(k, v) }

// RegisterBoolFlag registers a bool setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can be updated from a configuration file, an environment variable,
// or a flag. If k already exists a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned.
func RegisterBoolFlag(k, short string, v bool, dflt, usage string) error {
	return std.RegisterBoolFlag(k, short, v, dflt, usage)
}

// RegisterIntFlag registers an int setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can be updated from a configuration file, an environment variable,
// or a flag. If k already exists a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned.
func RegisterIntFlag(k, short string, v int, dflt, usage string) error {
	return std.RegisterIntFlag(k, short, v, dflt, usage)
}

// RegisterInt64Flag registers an int64 setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can be updated from a configuration file, an environment variable,
// or a flag. If k already exists a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned.
func RegisterInt64Flag(k, short string, v int64, dflt, usage string) error {
	return std.RegisterInt64Flag(k, short, v, dflt, usage)
}

// RegisterStringFlag registers a string setting with the standard settings
// using k for its key and v for its value. Once registered, the value of this
// setting can be updated from a configuration file, an environment variable,
// or a flag. If k already exists a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned.
func RegisterStringFlag(k, short, v, dflt, usage string) error {
	return std.RegisterStringFlag(k, short, v, dflt, usage)
}
