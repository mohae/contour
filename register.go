package contour

// Register contains all of contour's Register functions.Calling Register
// adds, or registers, the Settings information to the AppConfig variable.
// The setting value, if there is one, is not saved to its environment
// variable at this point.
//
// This allows for
//
// These should be called at app startup to register all configuration
// Settings that the application uses.

import "strconv"

// RegisterSetting registers a setting. For most settings, the data and setting
// type specific registration and add functions should be used. The exception
// would be when you want to customize what can override a setting: e.g. allow
// updates from env vars and flags only. If updating this setting, in some
// manner, is to be allowed, IsCore must be false as that will take precedence
// over any other type.
//
// If a setting with the same name already exists, a SettingExistsErr will be
// returned. If the name is an empty string an ErrNoSettingName will be
// returned.
//
// The short, dflt, and usage parms only apply to settings with IsFlag set to
// true.
//
// When IsCore is true, nothing can modify the setting's value once it is
// registered; usage of AddCore functions should be preferred.
//
// If the setting can be updated by a configuration file, environment variable
// or a flag, the IsConfFileVar, IsEnv, and IsFlag bools should be set to true
// as appropriate. These conditionals are independent; e.g. a setting can have
// both IsConfFileVar and IsFlag set to true if the setting is not to be
// updateable from an environment variable.
//
// If IsCore, IsConfFileVar, IsEnv, and IsFlag are all false, the setting will
// only be modifiable from application code via Update methods; usage of Add
// functions should be preferred. The setting will not be exposed to the
// configuration file, environment variables, or as flags.
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

// ConfFileVar settings are settable via a configuration file.  Only settings
// that are of type ConfFileVar, Env, and Flags can be set via a configuration
// file. ConfFileVar's cannot be set from environment variables or flags.

// RegisterBoolConfFileVar registers a bool setting with the key k and value v.
// The value of this setting can only be changed by a configuration once it is
// registered. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
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

// RegisterIntConfFileVar registers an int setting with the key k and value v.
// The value of this setting can only be changed by a configuration once it is
// registered. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
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

// RegisterInt64ConfFileVar registers an int64 setting with the key k and value
// The value of this setting can only be changed by a configuration once it is
// registered. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
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

// RegisterStringConfFileVar registers a string setting with the key k and
// The value of this setting can only be changed by a configuration once it is
// registered. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
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

// EnvVar settings are settable from the config file and environment variables.

// RegisterBoolEnvVar registers a bool setting with the key k and value v. The
// value of this setting can only be changed by a configuration file or an
// environment variable. If a setting with the same key k exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
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

// RegisterIntEnvVar registers an int setting with the key k and value v. The
// value of this setting can only be changed by a configuration file or an
// environment variable. If a setting with the same key k exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func (s *Settings) RegisterIntEnvVar(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntEnvVar(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) registerIntEnvVar(k string, v int) error {
	return s.registerEnvVar(_int, k, v, strconv.Itoa(v))
}

// RegisterInt64EnvVar registers an int64 setting with the key k and value v.
// The value of this setting can only be changed by a configuration file or an
// environment variable. If a setting with the same key k exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func (s *Settings) RegisterInt64EnvVar(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64EnvVar(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) registerInt64EnvVar(k string, v int64) error {
	return s.registerEnvVar(_int64, k, v, strconv.FormatInt(v, 10))
}

// RegisterStringEnvVar registers a string setting with the key k and value v.
// The value of this setting can only be changed by a configuration file or an
// environment variable. If a setting with the same key k exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
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

// Flag settings are settable from a configuration file, environment variable,
// and flags. If there is a short value, that will be the short flag alias for
// the setting.

// RegisterBoolFlag registers a bool setting with the key k and value v. Flag
// settings can be updated by configuration files, environment variables, and
// flags. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned. If the
// short value already is registered a ShortFlagExistsErr will be returend.
// This error includes the name of the setting that the short value is already
// registered to.
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

// RegisterIntFlag registers an int setting with the key k and value v. Flag
// settings can be updated by configuration files, environment variables, and
// flags. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned. If the
// short value already is registered a ShortFlagExistsErr will be returend.
// This error includes the name of the setting that the short value is already
// registered to.
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

// RegisterInt64Flag registers an int64 setting with the key k and value v. Flag
// settings can be updated by configuration files, environment variables, and
// flags. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned. If the
// short value already is registered a ShortFlagExistsErr will be returend.
// This error includes the name of the setting that the short value is already
// registered to.
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

// RegisterStringFlag registers a string setting with the key k and value v.
// Flag settings can be updated by configuration files, environment variables,
// and flags. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned. If the
// short value already is registered a ShortFlagExistsErr will be returend.
// This error includes the name of the setting that the short value is already
// registered to.
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

// RegisterSetting registers a setting. For most settings, the data and setting
// type specific registration and add functions should be used. The exception
// would be when you want to customize what can override a setting: e.g. allow
// updates from env vars and flags only. If updating this setting, in some
// manner, is to be allowed, IsCore must be false as that will take precedence
// over any other type.
//
// If a setting with the same name already exists, a SettingExistsErr will be
// returned. If the name is an empty string an ErrNoSettingName will be
// returned.
//
// The short, dflt, and usage parms only apply to settings with IsFlag set to
// true.
//
// When IsCore is true, nothing can modify the setting's value once it is
// registered; usage of AddCore functions should be preferred.
//
// If the setting can be updated by a configuration file, environment variable
// or a flag, the IsConfFileVar, IsEnv, and IsFlag bools should be set to true
// as appropriate. These conditionals are independent; e.g. a setting can have
// both IsConfFileVar and IsFlag set to true if the setting is not to be
// updateable from an environment variable.
//
// If IsCore, IsConfFileVar, IsEnv, and IsFlag are all false, the setting will
// only be modifiable from application code via Update methods; usage of Add
// functions should be preferred. The setting will not be exposed to the
// configuration file, environment variables, or as flags.
//
// For non string, bool, int, and int64 types, the type must be "interface{}"
// TODO: should typ be allowed be a custom string and treat all unknown
// values as an interface{}, e.g. typ=Foo and the Setting's type would be
// Foo even though returning the Setting's value would result in interface{}.
func RegisterSetting(typ, name, short string, value interface{}, dflt, usage string, IsCore, IsConfFileVar, IsEnv, IsFlag bool) error {
	return settings.RegisterSetting(typ, name, short, value, dflt, usage, IsCore, IsConfFileVar, IsEnv, IsFlag)
}

// RegisterBoolConfFileVar registers a bool setting with the key k and value v.
// The value of this setting can only be changed by a configuration once it is
// registered. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
func RegisterBoolConfFileVar(k string, v bool) error { return settings.RegisterBoolConfFileVar(k, v) }

// RegisterIntConfFileVar registers an int setting with the key k and value v.
// The value of this setting can only be changed by a configuration once it is
// registered. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
func RegisterIntConfFileVar(k string, v int) error { return settings.RegisterIntConfFileVar(k, v) }

// RegisterInt64ConfFileVar registers an int64 setting with the key k and value
// The value of this setting can only be changed by a configuration once it is
// registered. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
func RegisterInt64ConfFileVar(k string, v int64) error { return settings.RegisterInt64ConfFileVar(k, v) }

// RegisterStringConfFileVar registers a string setting with the key k and
// The value of this setting can only be changed by a configuration once it is
// registered. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
func RegisterStringConfFileVar(k, v string) error { return settings.RegisterStringConfFileVar(k, v) }

// RegisterBoolEnvVar registers a bool setting with the key k and value v. The
// value of this setting can only be changed by a configuration file or an
// environment variable. If a setting with the same key k exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func RegisterBoolEnvVar(k string, v bool) error { return settings.RegisterBoolEnvVar(k, v) }

// RegisterIntEnvVar registers an int setting with the key k and value v. The
// value of this setting can only be changed by a configuration file or an
// environment variable. If a setting with the same key k exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func RegisterIntEnvVar(k string, v int) error { return settings.RegisterIntEnvVar(k, v) }

// RegisterInt64EnvVar registers an int64 setting with the key k and value v.
// The value of this setting can only be changed by a configuration file or an
// environment variable. If a setting with the same key k exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func RegisterInt64EnvVar(k string, v int64) error { return settings.RegisterInt64EnvVar(k, v) }

// RegisterStringEnvVar registers a string setting with the key k and value v.
// The value of this setting can only be changed by a configuration file or an
// environment variable. If a setting with the same key k exists a
// SettingExistsErr will be returned. If k is empty, an ErrNoSettingName will
// be returned.
func RegisterStringEnvVar(k, v string) error { return settings.RegisterStringEnvVar(k, v) }

// RegisterBoolFlag registers a bool setting with the key k and value v. Flag
// settings can be updated by configuration files, environment variables, and
// flags. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned. If the
// short value already is registered a ShortFlagExistsErr will be returend.
// This error includes the name of the setting that the short value is already
// registered to.
func RegisterBoolFlag(k, short string, v bool, dflt, usage string) error {
	return settings.RegisterBoolFlag(k, short, v, dflt, usage)
}

// RegisterIntFlag registers an int setting with the key k and value v. Flag
// settings can be updated by configuration files, environment variables, and
// flags. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned. If the
// short value already is registered a ShortFlagExistsErr will be returend.
// This error includes the name of the setting that the short value is already
// registered to.
func RegisterIntFlag(k, short string, v int, dflt, usage string) error {
	return settings.RegisterIntFlag(k, short, v, dflt, usage)
}

// RegisterInt64Flag registers an int64 setting with the key k and value v. Flag
// settings can be updated by configuration files, environment variables, and
// flags. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned. If the
// short value already is registered a ShortFlagExistsErr will be returend.
// This error includes the name of the setting that the short value is already
// registered to.
func RegisterInt64Flag(k, short string, v int64, dflt, usage string) error {
	return settings.RegisterInt64Flag(k, short, v, dflt, usage)
}
// RegisterStringFlag registers a string setting with the key k and value v.
// Flag settings can be updated by configuration files, environment variables,
// and flags. If a setting with the same key k exists a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned. If the
// short value already is registered a ShortFlagExistsErr will be returend.
// This error includes the name of the setting that the short value is already
// registered to.
func RegisterStringFlag(k, short, v, dflt, usage string) error {
	return settings.RegisterStringFlag(k, short, v, dflt, usage)
}
