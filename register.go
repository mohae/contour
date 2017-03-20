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

// RegisterCfgFile set's the configuration file's name.  The name is parsed
// for a valid extension--one that is a supported format--and saves that value
// too. If it cannot be determined, the extension info is not set.  These are
// considered core values and cannot be changed from configuration files,
// environment variables, and configuration files.
//
// If the envName is a non-empty value, it is the environment variable name to
// check for a configuration filename.
func RegisterCfgFile(k, v string) error { return settings.RegisterCfgFile(k, v) }
func (s *Settings) RegisterCfgFile(k, v string) error {
	if v == "" {
		return fmt.Errorf("cannot register configuration file: no name provided")
	}
	if k == "" {
		return fmt.Errorf("cannot register configuration file: no key provided")
	}
	// store the key value being used as the configuration setting name by caller
	s.mu.Lock()
	s.confFileKey = k
	// cache this while we have the lock; technically racy but useEnv shouldn't
	// be modified while a config file is being registered.
	use := s.useEnv
	s.mu.Unlock()
	// check to see if the env var is set
	if use {
		fname := os.Getenv(s.GetEnvName(k))
		if fname != "" {
			v = fname
		}
	}
	s.RegisterStringCore(k, v)
	s.mu.Lock()
	s.useCfg = true
	s.mu.Unlock()
	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the new
// setting if it does not.
func RegisterSetting(typ, name, short string, value interface{}, dflt, usage string, IsCore, IsCfg, IsEnv, IsFlag bool) error {
	return settings.RegisterSetting(typ, name, short, value, dflt, usage, IsCore, IsCfg, IsEnv, IsFlag)
}
func (s *Settings) RegisterSetting(typ, name, short string, value interface{}, dflt string, usage string, IsCore, IsCfg, IsEnv, IsFlag bool) error {
	if name == "" {
		return fmt.Errorf("cannot register an unnamed setting")
	}
	s.mu.RLock()
	_, ok := s.settings[name]
	s.mu.RUnlock()
	if ok {
		// Settings can't be re-registered.
		return fmt.Errorf("%s is already registered, cannot re-register settings", name)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	// Add the setting
	s.settings[name] = setting{
		Type:    typ,
		Name:    name,
		Short:   short,
		Value:   value,
		Default: dflt,
		Usage:   usage,
		IsCore:  IsCore,
		IsCfg:   IsCfg,
		IsEnv:   IsEnv,
		IsFlag:  IsFlag,
	}
	// if it's a cfg file setting, add it to the cfgNames map
	if IsCfg {
		s.cfgVars[name] = struct{}{}
	}
	// mapping shortcodes make lookup easier
	if short != "" && IsFlag {
		_, ok := s.shortFlags[short]
		if ok {
			return fmt.Errorf("short flag %q is already in use; short flags must be unique", short)
		}
		s.shortFlags[short] = name
	}
	// Keep track of whether or not a cfg is being used. If a setting is registered
	// as a cfg setting, it is assumed a cfg source is being used.
	if IsEnv {
		s.useEnv = IsEnv
	}
	if IsCfg {
		s.useCfg = IsCfg
	}
	if IsFlag {
		s.useFlags = IsFlag
	}
	return nil
}

// Core settings are not overridable via cfg file, env vars, or command-line
// flags.  They can only be set via their respective Update() method or func.

// RegisterBoolCoreE adds the information to the appCfg global, but does not
// save it to its environment variable. E versions return received errors.
func RegisterBoolCoreE(k string, v bool) error { return settings.RegisterBoolCoreE(k, v) }
func (s *Settings) RegisterBoolCoreE(k string, v bool) error {
	return s.RegisterSetting("bool", k, "", v, strconv.FormatBool(v), "", true, false, false, false)
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
	return s.RegisterSetting("int", k, "", v, strconv.Itoa(v), "", true, false, false, false)
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
	return s.RegisterSetting("int64", k, "", v, strconv.FormatInt(v, 10), "", true, false, false, false)
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
	return s.RegisterSetting("string", k, "", v, v, "", true, false, false, false)
}

// RegisterStringCore calls RegisterStringCoreE and ignores any error.
func RegisterStringCore(k, v string) { settings.RegisterStringCore(k, v) }
func (s *Settings) RegisterStringCore(k, v string) {
	s.RegisterStringCoreE(k, v)
}

// Cfg settings are settable via a configuration file.  Only settings that are
// Cfg and Flags can be set via a cfg file. If the setting can be set from
// an environment variable, that variables name is passed via the "envName'
// parameter. If the envName == "" it will not be settable via an environment
// variable.

// RegisterBoolCfgE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterBoolCfgE(k string, v bool) error { return settings.RegisterBoolCfgE(k, v) }
func (s *Settings) RegisterBoolCfgE(k string, v bool) error {
	return s.RegisterSetting("bool", k, "", v, strconv.FormatBool(v), "", false, true, true, false)
}

// RegisterBoolCfg calls RegisterBoolCfgE and ignores any error.
func RegisterBoolCfg(k string, v bool) { settings.RegisterBoolCfg(k, v) }
func (s *Settings) RegisterBoolCfg(k string, v bool) {
	s.RegisterBoolCfgE(k, v)
}

// RegisterIntCfgE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterIntCfgE(k string, v int) error { return settings.RegisterIntCfgE(k, v) }
func (s *Settings) RegisterIntCfgE(k string, v int) error {
	return s.RegisterSetting("int", k, "", v, strconv.Itoa(v), "", false, true, true, false)
}

// RegisterIntCfg calls RegisterIntCfgE and ignores any error.
func RegisterIntCfg(k string, v int) { settings.RegisterIntCfg(k, v) }
func (s *Settings) RegisterIntCfg(k string, v int) {
	s.RegisterIntCfgE(k, v)
}

// RegisterInt64Cfg adds the informatio to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterInt64CfgE(k string, v int64) error { return settings.RegisterInt64CfgE(k, v) }
func (s *Settings) RegisterInt64CfgE(k string, v int64) error {
	return s.RegisterSetting("int64", k, "", v, strconv.FormatInt(v, 10), "", false, true, true, false)
}

// RegisterInt64Cfg calls RegisterInt64Cfg and ignores any error.
func RegisterInt64Cfg(k string, v int64) { settings.RegisterInt64Cfg(k, v) }
func (s *Settings) RegisterInt64Cfg(k string, v int64) {
	s.RegisterInt64CfgE(k, v)
}

// RegisterStringCfgE adds the information to the appCfg global, but does not
// save it to its environment variable: E versions return received errors.
func RegisterStringCfgE(k, v string) error { return settings.RegisterStringCfgE(k, v) }
func (s *Settings) RegisterStringCfgE(k, v string) error {
	return s.RegisterSetting("string", k, "", v, v, "", false, true, true, false)
}

// RegisterStringCfg calls RegisterStringCfgE and ignores any error.
func RegisterStringCfg(k, v string) { settings.RegisterStringCfg(k, v) }
func (s *Settings) RegisterStringCfg(k, v string) {
	s.RegisterStringCfgE(k, v)
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
	return s.RegisterSetting("bool", k, short, v, dflt, usage, false, true, true, true)
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
	return s.RegisterSetting("int", k, short, v, dflt, usage, false, true, true, true)
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
	return s.RegisterSetting("int64", k, short, v, dflt, usage, false, true, true, true)
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
	return s.RegisterSetting("string", k, short, v, dflt, usage, false, true, true, true)
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
	return s.RegisterSetting("bool", k, "", v, strconv.FormatBool(v), "", false, false, false, false)
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
	return s.RegisterSetting("int", k, "", v, strconv.Itoa(v), "", false, false, false, false)
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
	return s.RegisterSetting("int64", k, "", v, strconv.FormatInt(v, 10), "", false, false, false, false)
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
	return s.RegisterSetting("string", k, "", v, v, "", false, false, false, false)
}

// RegisterString calls RegisterStringE and ignores any error.
func RegisterString(k, v string) { settings.RegisterString(k, v) }
func (s *Settings) RegisterString(k, v string) {
	s.RegisterStringE(k, v)
}
