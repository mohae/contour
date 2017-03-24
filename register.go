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

import (
	"fmt"
	"strconv"
)

type RegistrationErr struct {
	name string
	slug string
}

func (e RegistrationErr) Error() string {
	if e.slug == "" {
		return fmt.Sprintf("%s: registration failed", e.name)
	}
	if e.name == "" {
		return fmt.Sprintf("registration failed: %s", e.slug)
	}
	return fmt.Sprintf("%s: registration failed: %s", e.name, e.slug)
}

// RegisterSetting registers a setting. For most settings, the data and setting
// type specific registration should be used. If an error occurs, a
// RegistrationErr will be returned. The exception would be when you want to
// customize what can override a setting: e.g. allow updates from env vars and
// flags only. If updating this setting, in some manner, is to be allowed,
// IsCore must be false as that will take precedence over any other type.
//
// The short, dflt, and usage parms only apply to settings with IsFlag set to
// true.
//
// For non string, bool, int, and int64 types, the type must be "interface{}".
func RegisterSetting(typ, name, short string, value interface{}, dflt, usage string, IsCore, IsConfFileVar, IsEnv, IsFlag bool) error {
	return settings.RegisterSetting(typ, name, short, value, dflt, usage, IsCore, IsConfFileVar, IsEnv, IsFlag)
}

// RegisterSetting registers a setting. For most settings, the data and setting
// type specific registration should be used. If an error occurs, a
// RegistrationErr will be returned. The exception would be when you want to
// customize what can override a setting: e.g. allow updates from env vars and
// flags only. If updating this setting, in some manner, is to be allowed,
// IsCore must be false as that will take precedence over any other type.
//
// The short, dflt, and usage parms only apply to settings with IsFlag set to
// true.
//
// For non string, bool, int, and int64 types, the type must be "interface{}".
func (s *Settings) RegisterSetting(typ, name, short string, value interface{}, dflt, usage string, IsCore, IsConfFileVar, IsEnv, IsFlag bool) error {
	dType, err := parseDataType(typ)
	if err != nil {
		return RegistrationErr{name: name, slug: err.Error()}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerSetting(dType, name, short, value, dflt, usage, IsCore, IsConfFileVar, IsEnv, IsFlag)
}

func (s *Settings) registerSetting(typ dataType, name, short string, value interface{}, dflt, usage string, IsCore, IsConfFileVar, IsEnv, IsFlag bool) error {
	if name == "" {
		return RegistrationErr{slug: "setting name was empty"}
	}
	_, ok := s.settings[name]
	if ok {
		// Settings can't be re-registered.
		return RegistrationErr{name: name, slug: "setting exists"}
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
		IsEnv:         IsEnv,
		IsFlag:        IsFlag,
	}
	// if it's a conf file setting, add it to the confFileVars map
	if IsConfFileVar {
		s.confFileVars[name] = struct{}{}
	}
	// mapping shortcodes make lookup easier
	if short != "" && IsFlag {
		_, ok := s.shortFlags[short]
		if ok {
			return RegistrationErr{name: name, slug: fmt.Sprintf("a setting using short flag %s exists; they must be unique", short)}
		}
		s.shortFlags[short] = name
	}
	if IsEnv {
		s.useEnv = IsEnv
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

// Core settings are not overridable via a configuration file, env vars, or
// command-line flags. They cannot be modified in any way once they have been
// registered.

// RegisterBoolCore registers a bool setting with the key k and value v. The
// value of this setting cannot be changed once it is registered. If an error
// occurs, a RegistrationErr will be returned.
func RegisterBoolCore(k string, v bool) error { return settings.RegisterBoolCore(k, v) }

// RegisterBoolCore registers a bool setting with the key k and value v. The
// value of this setting cannot be changed once it is registered. If an error
// occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterBoolCore(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBoolCore(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBoolCore(k string, v bool) error {
	return s.registerSetting(_bool, k, "", v, strconv.FormatBool(v), "", true, false, false, false)
}

// RegisterIntCore registers an int settings with the key k and value v. The
// value of this setting cannot be changed once it is registered. If an error
// occurs, a RegistrationErr will be returned.
func RegisterIntCore(k string, v int) error { return settings.RegisterIntCore(k, v) }

// RegisterIntCore registers an int settings with the key k and value v. The
// value of this setting cannot be changed once it is registered. If an error
// occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterIntCore(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntCore(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerIntCore(k string, v int) error {
	return s.registerSetting(_int, k, "", v, strconv.Itoa(v), "", true, false, false, false)
}

// RegisterInt64Core registers an int64 settings with the key k and value v.
// The value of this setting cannot be changed once it is registered. If an
// error occurs, a RegistrationErr will be returned.
func RegisterInt64Core(k string, v int64) error { return settings.RegisterInt64Core(k, v) }

// RegisterInt64Core registers an int64 settings with the key k and value v.
// The value of this setting cannot be changed once it is registered. If an
// error occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterInt64Core(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64Core(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64Core(k string, v int64) error {
	return s.registerSetting(_int64, k, "", v, strconv.FormatInt(v, 10), "", true, false, false, false)
}

// RegisterStringCore registers an string settings with the key k and value v.
// The value of this setting cannot be changed once it is registered. If an
// error occurs, a RegistrationErr will be returned.
func RegisterStringCore(k, v string) error { return settings.RegisterStringCore(k, v) }

// RegisterStringCore registers an string settings with the key k and value v.
// The value of this setting cannot be changed once it is registered. If an
// error occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterStringCore(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerStringCore(k, v)
}

func (s *Settings) registerStringCore(k, v string) error {
	return s.registerSetting(_string, k, "", v, v, "", true, false, false, false)
}

// ConfFileVar settings are settable via a configuration file.  Only settings
// that are of type ConfFileVar, Env, and Flags can be set via a configuration
// file. ConfFileVar's cannot be set from environment variables or flags.

// RegisterBoolConfFileVar registers a bool setting with the key k and value v.
// The value of this setting can only be changed by a configuration once it is
// registered. If an error occurs, a RegistrationErr will be returned.
func RegisterBoolConfFileVar(k string, v bool) error { return settings.RegisterBoolConfFileVar(k, v) }

// RegisterBoolConfFileVar registers a bool setting with the key k and value v.
// The value of this setting can only be changed by a configuration once it is
// registered. If an error occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterBoolConfFileVar(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBoolConfFileVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBoolConfFileVar(k string, v bool) error {
	return s.registerSetting(_bool, k, "", v, strconv.FormatBool(v), "", false, true, true, false)
}

// RegisterIntConfFileVar registers an int setting with the key k and value v.
// The value of this setting can only be changed by a configuration once it is
// registered. If an error occurs, a RegistrationErr will be returned.
func RegisterIntConfFileVar(k string, v int) error { return settings.RegisterIntConfFileVar(k, v) }

// RegisterIntConfFileVar registers an int setting with the key k and value v.
// The value of this setting can only be changed by a configuration once it is
// registered. If an error occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterIntConfFileVar(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntConfFileVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerIntConfFileVar(k string, v int) error {
	return s.registerSetting(_int, k, "", v, strconv.Itoa(v), "", false, true, true, false)
}

// RegisterInt64ConfFileVar registers an int64 setting with the key k and value
// The value of this setting can only be changed by a configuration once it is
// registered. If an error occurs, a RegistrationErr will be returned.
func RegisterInt64ConfFileVar(k string, v int64) error { return settings.RegisterInt64ConfFileVar(k, v) }

// RegisterInt64ConfFileVar registers an int64 setting with the key k and value
// The value of this setting can only be changed by a configuration once it is
// registered. If an error occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterInt64ConfFileVar(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64ConfFileVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64ConfFileVar(k string, v int64) error {
	return s.registerSetting(_int64, k, "", v, strconv.FormatInt(v, 10), "", false, true, true, false)
}

// RegisterStringConfFileVar registers a string setting with the key k and
// value v. The value of this setting can only be changed by a configuration
// once it is registered. If an error occurs, a RegistrationErr will be
// returned.
func RegisterStringConfFileVar(k, v string) error { return settings.RegisterStringConfFileVar(k, v) }

// RegisterStringConfFileVar registers a string setting with the key k and
// value v. The value of this setting can only be changed by a configuration
// once it is registered. If an error occurs, a RegistrationErr will be
// returned.
func (s *Settings) RegisterStringConfFileVar(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerStringConfFileVar(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerStringConfFileVar(k, v string) error {
	return s.registerSetting(_string, k, "", v, v, "", false, true, true, false)
}

// Flag settings are settable from the config file and as command-line flags.
// If there is a short value, that will be the short flag alias for the
// setting. Only settings that are of type ConfFileVar, EnvVar, and Flag can be
// set by a flag. A flag can be set by configuration variable, environment
// variable, and command-line argument.

// RegisterBoolFlag registers a bool setting with the key k and value v. The
// value of this setting can be changed by a configuration file, environment
// variable, or a flag. If an error occurs, a RegistrationErr will be returned.
func RegisterBoolFlag(k, short string, v bool, dflt, usage string) error {
	return settings.RegisterBoolFlag(k, short, v, dflt, usage)
}

// RegisterBoolFlag registers a bool setting with the key k and value v. The
// value of this setting can be changed by a configuration file, environment
// variable, or a flag. If an error occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterBoolFlag(k, short string, v bool, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBoolFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBoolFlag(k, short string, v bool, dflt, usage string) error {
	return s.registerSetting(_bool, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterIntFlag registers an int setting with the key k and value v. The
// value of this setting can be changed by a configuration file, environment
// variable, or a flag. If an error occurs, a RegistrationErr will be returned.
func RegisterIntFlag(k, short string, v int, dflt, usage string) error {
	return settings.RegisterIntFlag(k, short, v, dflt, usage)
}

// RegisterIntFlag registers an int setting with the key k and value v. The
// value of this setting can be changed by a configuration file, environment
// variable, or a flag. If an error occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterIntFlag(k, short string, v int, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerIntFlag(k, short string, v int, dflt, usage string) error {
	return s.registerSetting(_int, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterInt64Flag registers an int64 setting with the key k and value v. The
// value of this setting can be changed by a configuration file, environment
// variable, or a flag. If an error occurs, a RegistrationErr will be returned.
func RegisterInt64Flag(k, short string, v int64, dflt, usage string) error {
	return settings.RegisterInt64Flag(k, short, v, dflt, usage)
}

// RegisterInt64Flag registers an int64 setting with the key k and value v. The
// value of this setting can be changed by a configuration file, environment
// variable, or a flag. If an error occurs, a RegistrationErr will be returned.
func (s *Settings) RegisterInt64Flag(k, short string, v int64, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64Flag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64Flag(k, short string, v int64, dflt, usage string) error {
	return s.registerSetting(_int64, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterStringFlag registers a string setting with the key k and value v.
// The value of this setting can be changed by a configuration file,
// environment variable, or a flag. If an error occurs, a RegistrationErr will
// be returned.
func RegisterStringFlag(k, short, v, dflt, usage string) error {
	return settings.RegisterStringFlag(k, short, v, dflt, usage)
}

// RegisterStringFlag registers a string setting with the key k and value v.
// The value of this setting can be changed by a configuration file,
// environment variable, or a flag. If an error occurs, a RegistrationErr will
// be returned.
func (s *Settings) RegisterStringFlag(k, short, v, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerStringFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerStringFlag(k, short, v, dflt, usage string) error {
	return s.registerSetting(_string, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterBool registers a bool setting with they key k and value f. This
// can be updated within the application but is not updated by configuration
// files, environment variables, or flags. If an error occurs, a
// RegistrationErr will be returned.
func RegisterBool(k string, v bool) error { return settings.RegisterBool(k, v) }

// RegisterBool registers a bool setting with they key k and value f. This
// can be updated within the application but is not updated by configuration
// files, environment variables, or flags. If an error occurs, a
// RegistrationErr will be returned.
func (s *Settings) RegisterBool(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBool(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBool(k string, v bool) error {
	return s.registerSetting(_bool, k, "", v, strconv.FormatBool(v), "", false, false, false, false)
}

// RegisterInt registers an int setting with they key k and value f. This can
// be updated within the application but is not updated by configuration files,
// environment variables, or flags. If an error occurs, a RegistrationErr will
// be returned.
func RegisterInt(k string, v int) error { return settings.RegisterInt(k, v) }

// RegisterInt registers an int setting with they key k and value f. This can
// be updated within the application but is not updated by configuration files,
// environment variables, or flags. If an error occurs, a RegistrationErr will
// be returned.
func (s *Settings) RegisterInt(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt(k string, v int) error {
	return s.registerSetting(_int, k, "", v, strconv.Itoa(v), "", false, false, false, false)
}

// RegisterInt64 registers an int64 setting with they key k and value f. This
// can be updated within the application but is not updated by configuration
// files, environment variables, or flags. If an error occurs, a
// RegistrationErr will be returned.
func RegisterInt64(k string, v int64) error { return settings.RegisterInt64(k, v) }

// RegisterInt64 registers an int64 setting with they key k and value f. This
// can be updated within the application but is not updated by configuration
// files, environment variables, or flags. If an error occurs, a
// RegistrationErr will be returned.
func (s *Settings) RegisterInt64(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64(k string, v int64) error {
	return s.registerSetting(_int64, k, "", v, strconv.FormatInt(v, 10), "", false, false, false, false)
}

// RegisterString registers a string setting with they key k and value f. This
// can be updated within the application but is not updated by configuration
// files, environment variables, or flags. If an error occurs, a
// RegistrationErr will be returned.
func RegisterString(k, v string) error { return settings.RegisterString(k, v) }
func (s *Settings) RegisterString(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerString(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerString(k, v string) error {
	return s.registerSetting(_string, k, "", v, v, "", false, false, false, false)
}
