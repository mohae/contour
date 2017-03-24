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
	"os"
	"strconv"
)

type RegisterErr struct {
	name string
	err  error
}

func (e RegisterErr) Error() string {
	return fmt.Sprintf("%s: registration failed: %s", e.name, e.err)
}

// RegisterConfFilename set's the configuration file's name.  The name is parsed
// for a valid extension--one that is a supported format--and saves that value
// too. If it cannot be determined, the extension info is not set.  These are
// considered core values and cannot be changed from configuration files,
// environment variables, and configuration files.
//
// If the envName is a non-empty value, it is the environment variable name to
// check for a configuration filename.
func RegisterConfFilename(k, v string) error { return settings.RegisterConfFilename(k, v) }
func (s *Settings) RegisterConfFilename(k, v string) error {
	if v == "" {
		return fmt.Errorf("cannot register configuration file: no name provided")
	}

	// update the confFilenameVarName if the value isn't empty; otherwise the default will be used
	if k != "" {
		s.confFilenameVarName = k
	}
	// store the key value being used as the configuration setting name by caller
	s.mu.Lock()
	defer s.mu.Unlock()

	// cache this while we have the lock; technically racy but useEnv shouldn't
	// be modified while a config file is being registered.
	use := s.useEnv
	// check to see if the env var is set
	if use {
		fname := os.Getenv(s.GetEnvName(k))
		if fname != "" {
			v = fname
		}
	}
	s.registerStringCore(s.confFilenameVarName, v)
	s.useConfFile = true
	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the new
// setting if it does not.
func RegisterSetting(typ, name, short string, value interface{}, dflt, usage string, IsCore, IsConfFileVar, IsEnv, IsFlag bool) error {
	return settings.RegisterSetting(typ, name, short, value, dflt, usage, IsCore, IsConfFileVar, IsEnv, IsFlag)
}
func (s *Settings) RegisterSetting(typ, name, short string, value interface{}, dflt string, usage string, IsCore, IsConfFileVar, IsEnv, IsFlag bool) error {
	dType, err := parseDataType(typ)
	if err != nil {
		return RegisterErr{name: name, err: err}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerSetting(dType, name, short, value, dflt, usage, IsCore, IsConfFileVar, IsEnv, IsFlag)
}

func (s *Settings) registerSetting(typ dataType, name, short string, value interface{}, dflt string, usage string, IsCore, IsConfFileVar, IsEnv, IsFlag bool) error {
	if name == "" {
		return fmt.Errorf("cannot register an unnamed setting")
	}
	_, ok := s.settings[name]
	if ok {
		// Settings can't be re-registered.
		return fmt.Errorf("%s is already registered, cannot re-register settings", name)
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
			return fmt.Errorf("short flag %q is already in use; short flags must be unique", short)
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
// command-line flags.

// RegisterBoolCoreE adds the information to the package global, but does not
// save it to its environment variable. E versions return received errors.
func RegisterBoolCoreE(k string, v bool) error { return settings.RegisterBoolCoreE(k, v) }
func (s *Settings) RegisterBoolCoreE(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBoolCore(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBoolCore(k string, v bool) error {
	return s.registerSetting(_bool, k, "", v, strconv.FormatBool(v), "", true, false, false, false)
}

// RegisterBoolCore calls RegisterBoolCoreE and ignores any error.
func RegisterBoolCore(k string, v bool) { settings.RegisterBoolCore(k, v) }
func (s *Settings) RegisterBoolCore(k string, v bool) {
	s.RegisterBoolCoreE(k, v)
}

// RegisterIntCoreE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterIntCoreE(k string, v int) error { return settings.RegisterIntCoreE(k, v) }
func (s *Settings) RegisterIntCoreE(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntCore(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerIntCore(k string, v int) error {
	return s.registerSetting(_int, k, "", v, strconv.Itoa(v), "", true, false, false, false)
}

// RegisterIntCore calls RegisterIntCoreE and ignores any error.
func RegisterIntCore(k string, v int) { settings.RegisterIntCore(k, v) }
func (s *Settings) RegisterIntCore(k string, v int) {
	s.RegisterIntCoreE(k, v)
}

// RegisterInt64CoreE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterInt64CoreE(k string, v int64) error { return settings.RegisterInt64CoreE(k, v) }
func (s *Settings) RegisterInt64CoreE(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64Core(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64Core(k string, v int64) error {
	return s.registerSetting(_int64, k, "", v, strconv.FormatInt(v, 10), "", true, false, false, false)
}

// RegisterInt64Core calls RegisterInt64CoreE and ignores any error.
func RegisterInt64Core(k string, v int64) { settings.RegisterInt64Core(k, v) }
func (s *Settings) RegisterInt64Core(k string, v int64) {
	s.RegisterInt64CoreE(k, v)
}

// RegisterStringCoreE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterStringCoreE(k, v string) error { return settings.RegisterStringCoreE(k, v) }
func (s *Settings) RegisterStringCoreE(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerStringCore(k, v)
}

func (s *Settings) registerStringCore(k, v string) error {
	return s.registerSetting(_string, k, "", v, v, "", true, false, false, false)
}

// RegisterStringCore calls RegisterStringCoreE and ignores any error.
func RegisterStringCore(k, v string) { settings.RegisterStringCore(k, v) }
func (s *Settings) RegisterStringCore(k, v string) {
	s.RegisterStringCoreE(k, v)
}

// ConfFileVar settings are settable via a configuration file.  Only settings that
// are ConfFileVar, Env, and Flags can be set via a configuration file.

// RegisterBoolConfFileVar registers a bool setting using name and its value
// set to v. If an error occurs, a RegistrationErr will be returned.
func RegisterBoolConfFileVar(k string, v bool) error { return settings.RegisterBoolConfFileVar(k, v) }

// RegisterBoolConfFileVar registers a bool setting using and its value set to
// v. If an error occurs, a RegistrationErr will be returned.
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

// RegisterIntConfFileVar registers an int setting using name and its value set
// to v. If an error occurs, a RegistrationErr will be returned.
func RegisterIntConfFileVar(k string, v int) error { return settings.RegisterIntConfFileVar(k, v) }

// RegisterIntConfFileVar registers an int setting using name and its value set
// to v. If an error occurs, a RegistrationErr will be returned.
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

// RegisterInt64ConfFileVar registers an int64 settings using name and its
// value set to v. If an error occurs, a RegistrationErr will be returned.
func RegisterInt64ConfFileVar(k string, v int64) error { return settings.RegisterInt64ConfFileVar(k, v) }

// RegisterInt64ConfFileVar registers an int64 settings using name and its
// value set to v. If an error occurs, a RegistrationErr will be returned.
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

// RegisterStringConfFileVar registers a string setting using name and its
// value set to v. If an error occurs, a RegistrationErr will be returned.
func RegisterStringConfFileVar(k, v string) error { return settings.RegisterStringConfFileVar(k, v) }

// RegisterStringConfFileVar registers a string setting using name and its
// value set to v. If an error occurs, a RegistrationErr will be returned.
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
// Only settings that are Cfg and Flags can be set via a cfg file.  If the
// setting can be set from an environment variable, that variables name is
// passed via the "envName' parameter.  If there is a value for the "short
// code(s)" parameter, that value will be used as that flag's command-line
// short code.  If the envName == "" it will not be settable via an
// environment variable.

// RegisterBoolFlagE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterBoolFlagE(k, short string, v bool, dflt, usage string) error {
	return settings.RegisterBoolFlagE(k, short, v, dflt, usage)
}
func (s *Settings) RegisterBoolFlagE(k, short string, v bool, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBoolFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBoolFlag(k, short string, v bool, dflt, usage string) error {
	return s.registerSetting(_bool, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterBoolFlag calls RegisterBoolFlagE and ignores any error.
func RegisterBoolFlag(k, short string, v bool, dflt, usage string) {
	settings.RegisterBoolFlag(k, short, v, dflt, usage)
}
func (s *Settings) RegisterBoolFlag(k, short string, v bool, dflt, usage string) {
	s.RegisterBoolFlagE(k, short, v, dflt, usage)
}

// RegisterIntFlagE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterIntFlagE(k, short string, v int, dflt, usage string) error {
	return settings.RegisterIntFlagE(k, short, v, dflt, usage)
}
func (s *Settings) RegisterIntFlagE(k, short string, v int, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerIntFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerIntFlag(k, short string, v int, dflt, usage string) error {
	return s.registerSetting(_int, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterIntFlag calls RegisterIntFlagE and ignores any error.
func RegisterIntFlag(k, short string, v int, dflt, usage string) {
	settings.RegisterIntFlag(k, short, v, dflt, usage)
}
func (s *Settings) RegisterIntFlag(k, short string, v int, dflt, usage string) {
	s.RegisterIntFlagE(k, short, v, dflt, usage)
}

// RegisterInt64FlagE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterInt64FlagE(k, short string, v int64, dflt, usage string) error {
	return settings.RegisterInt64FlagE(k, short, v, dflt, usage)
}
func (s *Settings) RegisterInt64FlagE(k, short string, v int64, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64Flag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64Flag(k, short string, v int64, dflt, usage string) error {
	return s.registerSetting(_int64, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterInt64Flag calls RegisterIntFlagE and ignores any error.
func RegisterInt64Flag(k, short string, v int64, dflt, usage string) {
	settings.RegisterInt64Flag(k, short, v, dflt, usage)
}
func (s *Settings) RegisterInt64Flag(k, short string, v int64, dflt, usage string) {
	s.RegisterInt64FlagE(k, short, v, dflt, usage)
}

// RegisterStringFlagE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterStringFlagE(k, short, v, dflt, usage string) error {
	return settings.RegisterStringFlagE(k, short, v, dflt, usage)
}
func (s *Settings) RegisterStringFlagE(k, short, v, dflt, usage string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerStringFlag(k, short, v, dflt, usage)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerStringFlag(k, short, v, dflt, usage string) error {
	return s.registerSetting(_string, k, short, v, dflt, usage, false, true, true, true)
}

// RegisterStringFlag calls RegisterStringFlagE and ignores any error.
func RegisterStringFlag(k, short, v, dflt, usage string) {
	settings.RegisterStringFlag(k, short, v, dflt, usage)
}
func (s *Settings) RegisterStringFlag(k, short, v, dflt, usage string) {
	s.RegisterStringFlagE(k, short, v, dflt, usage)
}

// RegisterBoolE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterBoolE(k string, v bool) error { return settings.RegisterBoolE(k, v) }
func (s *Settings) RegisterBoolE(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerBool(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerBool(k string, v bool) error {
	return s.registerSetting(_bool, k, "", v, strconv.FormatBool(v), "", false, false, false, false)
}

// RegisterBool calls RegisterBoolE and ignores any error.
func RegisterBool(k string, v bool) { settings.RegisterBool(k, v) }
func (s *Settings) RegisterBool(k string, v bool) {
	s.RegisterBoolE(k, v)
}

// RegisterIntE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterIntE(k string, v int) error { return settings.RegisterIntE(k, v) }
func (s *Settings) RegisterIntE(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt(k string, v int) error {
	return s.registerSetting(_int, k, "", v, strconv.Itoa(v), "", false, false, false, false)
}

// RegisterInt calls RegisterIntE and ignores any error.
func RegisterInt(k string, v int) { settings.RegisterInt(k, v) }
func (s *Settings) RegisterInt(k string, v int) {
	s.RegisterIntE(k, v)
}

// RegisterInt64E adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterInt64E(k string, v int64) error { return settings.RegisterInt64E(k, v) }
func (s *Settings) RegisterInt64E(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerInt64(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerInt64(k string, v int64) error {
	return s.registerSetting(_int64, k, "", v, strconv.FormatInt(v, 10), "", false, false, false, false)
}

// RegisterInt64 calls RegisterInt64E and ignores any error.
func RegisterInt64(k string, v int64) { settings.RegisterInt64(k, v) }
func (s *Settings) RegisterInt64(k string, v int64) {
	s.RegisterInt64E(k, v)
}

// RegisterStringE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterStringE(k, v string) error { return settings.RegisterStringE(k, v) }
func (s *Settings) RegisterStringE(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.registerString(k, v)
}

// assumes the lock has been obtained. Unexported register methods always
// return an error.
func (s *Settings) registerString(k, v string) error {
	return s.registerSetting(_string, k, "", v, v, "", false, false, false, false)
}

// RegisterString calls RegisterStringE and ignores any error.
func RegisterString(k, v string) { settings.RegisterString(k, v) }
func (s *Settings) RegisterString(k, v string) {
	s.RegisterStringE(k, v)
}
