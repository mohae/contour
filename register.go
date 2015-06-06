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

// RegisterConfigFile set's the configuration file's name.  The name is parsed
// for a valid extension--one that is a supported format--and saves that value
// too. If it cannot be determined, the extension info is not set.  These are
// considered core values and cannot be changed from configuration files,
// environment variables, and configuration files.
//
// If the envName is a non-empty value, it is the environment variable name to
// check for a configuration filename.
func RegisterCfgFile(k, v string) error { return appCfg.RegisterCfgFile(k, v) }
func (c *Cfg) RegisterCfgFile(k, v string) error {
	if v == "" {
		return fmt.Errorf("RegisterCfgFile expected a cfg filename: none received")
	}
	if k == "" {
		return fmt.Errorf("RegisterCfgFile expected a configuration key: none received")
	}
	// check to see if the env var is set
	c.RWMutex.RLock()
	if c.useEnv {
		fname := os.Getenv(fmt.Sprintf("%s_%s", c.name, k))
		if fname != "" {
			v = fname
		}
	}
	c.RWMutex.RUnlock()
	c.RegisterStringCore(k, v)
	// Register it first. If a valid cfg format isn't found, an error/ will be
	// returned, so registering it afterwords would mean the setting would not
	// exist.
	c.RegisterString(CfgFormat, "")
	format, err := formatFromFilename(v)
	if err != nil {
		return err
	}
	// Now we can update the format, since it wasn't set before, it can be set now
	// before it becomes read only.
	c.UpdateString(CfgFormat, format.String())
	c.RWMutex.Lock()
	c.useCfgFile = true
	c.RWMutex.Unlock()
	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the new
// setting if it does not.
func RegisterSetting(typ, name, short, envName string, value interface{}, dflt, usage string, IsCore, IsCfg, IsEnv, IsFlag bool) error {
	return appCfg.RegisterSetting(typ, name, short, envName, value, dflt, usage, IsCore, IsCfg, IsEnv, IsFlag)
}
func (c *Cfg) RegisterSetting(typ, name, short, envName string, value interface{}, dflt string, usage string, IsCore, IsCfg, IsEnv, IsFlag bool) error {
	c.RWMutex.RLock()
	_, ok := c.settings[name]
	if ok {
		// Settings can't be re-registered.
		c.RWMutex.RUnlock()
		return fmt.Errorf("%s is already registered, cannot re-register settings", name)
	}
	c.RWMutex.RUnlock()
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	// Add the setting
	c.settings[name] = setting{
		Type:    typ,
		Name:    name,
		Short:   short,
		EnvName: envName,
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
		c.cfgVars[name] = struct{}{}
	}
	// mapping shortcodes make lookup easier
	if short != "" && IsFlag {
		_, ok := c.shortFlags[short]
		if ok {
			return fmt.Errorf("short flag %q is already in use; short flags must be unique", short)
		}
		c.shortFlags[short] = name
	}
	// Keep track of whether or not a config is being used. If a setting is
	// registered as a config setting, it is assumed a configuration source
	// is being used.
	c.useEnv = IsEnv
	c.useCfgFile = IsCfg
	c.useFlags = IsFlag
	return nil
}

// Core settings are not overridable via cfg file, env vars, or command-line
// flags.  They can only be set via their respective Update() method or func.

// RegisterBoolCoreE adds the information to the AppsConfig struct, but does not
// save it to its environment variable. E versions return received errors.
func RegisterBoolCoreE(k string, v bool) error { return appCfg.RegisterBoolCoreE(k, v) }
func (c *Cfg) RegisterBoolCoreE(k string, v bool) error {
	return c.RegisterSetting("bool", k, "", "", v, strconv.FormatBool(v), "", true, false, false, false)
}

// RegisterBoolCore calls RegisterBoolCoreE and ignores any error.
func RegisterBoolCore(k string, v bool) { appCfg.RegisterBoolCore(k, v) }
func (c *Cfg) RegisterBoolCore(k string, v bool) {
	c.RegisterBoolCoreE(k, v)
}

// RegisterIntCoreE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterIntCoreE(k string, v int) error { return appCfg.RegisterIntCoreE(k, v) }
func (c *Cfg) RegisterIntCoreE(k string, v int) error {
	return c.RegisterSetting("int", k, "", "", v, strconv.Itoa(v), "", true, false, false, false)
}

// RegisterIntCore calls RegisterIntCoreE and ignores any error.
func RegisterIntCore(k string, v int) { appCfg.RegisterIntCore(k, v) }
func (c *Cfg) RegisterIntCore(k string, v int) {
	c.RegisterIntCoreE(k, v)
}

// RegisterInt64CoreE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterInt64CoreE(k string, v int64) error { return appCfg.RegisterInt64CoreE(k, v) }
func (c *Cfg) RegisterInt64CoreE(k string, v int64) error {
	return c.RegisterSetting("int64", k, "", "", v, strconv.FormatInt(v, 10), "", true, false, false, false)
}

// RegisterInt64Core calls RegisterInt64CoreE and ignores any error.
func RegisterInt64Core(k string, v int64) { appCfg.RegisterInt64Core(k, v) }
func (c *Cfg) RegisterInt64Core(k string, v int64) {
	c.RegisterInt64CoreE(k, v)
}

// RegisterStringCoreE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterStringCoreE(k, v string) error { return appCfg.RegisterStringCoreE(k, v) }
func (c *Cfg) RegisterStringCoreE(k, v string) error {
	return c.RegisterSetting("string", k, "", "", v, v, "", true, false, false, false)
}

// RegisterStringCore calls RegisterStringCoreE and ignores any error.
func RegisterStringCore(k, v string) { appCfg.RegisterStringCore(k, v) }
func (c *Cfg) RegisterStringCore(k, v string) {
	c.RegisterStringCoreE(k, v)
}

// Cfg settings are settable via a configuration file.  Only settings that are
// Cfg and Flags can be set via a cfg file. If the setting can be set from
// an environment variable, that variables name is passed via the "envName'
// parameter. If the envName == "" it will not be settable via an environment
// variable.

// RegisterBoolCfgE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterBoolCfgE(k, envName string, v bool) error { return appCfg.RegisterBoolCfgE(k, envName, v) }
func (c *Cfg) RegisterBoolCfgE(k, envName string, v bool) error {
	return c.RegisterSetting("bool", k, "", envName, v, strconv.FormatBool(v), "", false, true, true, false)
}

// RegisterBoolCfg calls RegisterBoolCfgE and ignores any error.
func RegisterBoolCfg(k, envName string, v bool) { appCfg.RegisterBoolCfg(k, envName, v) }
func (c *Cfg) RegisterBoolCfg(k, envName string, v bool) {
	c.RegisterBoolCfgE(k, envName, v)
}

// RegisterIntCfgE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterIntCfgE(k, envName string, v int) error { return appCfg.RegisterIntCfgE(k, envName, v) }
func (c *Cfg) RegisterIntCfgE(k, envName string, v int) error {
	return c.RegisterSetting("int", k, "", envName, v, strconv.Itoa(v), "", false, true, true, false)
}

// RegisterIntCfg calls RegisterIntCfgE and ignores any error.
func RegisterIntCfg(k, envName string, v int) { appCfg.RegisterIntCfg(k, envName, v) }
func (c *Cfg) RegisterIntCfg(k, envName string, v int) {
	c.RegisterIntCfgE(k, envName, v)
}

// RegisterInt64Cfg adds the informatio to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterInt64CfgE(k, envName string, v int64) error {
	return appCfg.RegisterInt64CfgE(k, envName, v)
}
func (c *Cfg) RegisterInt64CfgE(k, envName string, v int64) error {
	return c.RegisterSetting("int64", k, "", envName, v, strconv.FormatInt(v, 10), "", false, true, true, false)
}

// RegisterInt64Cfg calls RegisterInt64Cfg and ignores any error.
func RegisterInt64Cfg(k, envName string, v int64) { appCfg.RegisterInt64Cfg(k, envName, v) }
func (c *Cfg) RegisterInt64Cfg(k, envName string, v int64) {
	c.RegisterInt64CfgE(k, envName, v)
}

// RegisterStringCfgE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterStringCfgE(k, envName, v string) error { return appCfg.RegisterStringCfgE(k, envName, v) }
func (c *Cfg) RegisterStringCfgE(k, envName, v string) error {
	return c.RegisterSetting("string", k, "", envName, v, v, "", false, true, true, false)
}

// RegisterStringCfg calls RegisterStringCfgE and ignores any error.
func RegisterStringCfg(k, envName, v string) { appCfg.RegisterStringCfg(k, envName, v) }
func (c *Cfg) RegisterStringCfg(k, envName, v string) {
	c.RegisterStringCfgE(k, envName, v)
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
func RegisterBoolFlagE(k, s, envName string, v bool, dflt, usage string) error {
	return appCfg.RegisterBoolFlagE(k, s, envName, v, dflt, usage)
}
func (c *Cfg) RegisterBoolFlagE(k, s, envName string, v bool, dflt, usage string) error {
	return c.RegisterSetting("bool", k, s, envName, v, dflt, usage, false, true, true, true)
}

// RegisterBoolFlag calls RegisterBoolFlagE and ignores any error.
func RegisterBoolFlag(k, s, envName string, v bool, dflt, usage string) {
	appCfg.RegisterBoolFlag(k, s, envName, v, dflt, usage)
}
func (c *Cfg) RegisterBoolFlag(k, s, envName string, v bool, dflt, usage string) {
	c.RegisterBoolFlagE(k, s, envName, v, dflt, usage)
}

// RegisterIntFlagE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterIntFlagE(k, s, envName string, v int, dflt, usage string) error {
	return appCfg.RegisterIntFlagE(k, s, envName, v, dflt, usage)
}
func (c *Cfg) RegisterIntFlagE(k, s, envName string, v int, dflt, usage string) error {
	return c.RegisterSetting("int", k, s, envName, v, dflt, usage, false, true, true, true)
}

// RegisterIntFlag calls RegisterIntFlagE and ignores any error.
func RegisterIntFlag(k, s, envName string, v int, dflt, usage string) {
	appCfg.RegisterIntFlag(k, s, envName, v, dflt, usage)
}
func (c *Cfg) RegisterIntFlag(k, s, envName string, v int, dflt, usage string) {
	c.RegisterIntFlagE(k, s, envName, v, dflt, usage)
}

// RegisterInt64FlagE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterInt64FlagE(k, s, envName string, v int64, dflt, usage string) error {
	return appCfg.RegisterInt64FlagE(k, s, envName, v, dflt, usage)
}
func (c *Cfg) RegisterInt64FlagE(k, s, envName string, v int64, dflt, usage string) error {
	return c.RegisterSetting("int64", k, s, envName, v, dflt, usage, false, true, true, true)
}

// RegisterInt64Flag calls RegisterIntFlagE and ignores any error.
func RegisterInt64Flag(k, s, envName string, v int64, dflt, usage string) {
	appCfg.RegisterInt64Flag(k, s, envName, v, dflt, usage)
}
func (c *Cfg) RegisterInt64Flag(k, s, envName string, v int64, dflt, usage string) {
	c.RegisterInt64FlagE(k, s, envName, v, dflt, usage)
}

// RegisterStringFlagE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterStringFlagE(k, s, envName, v, dflt, usage string) error {
	return appCfg.RegisterStringFlagE(k, s, envName, v, dflt, usage)
}
func (c *Cfg) RegisterStringFlagE(k, s, envName, v, dflt, usage string) error {
	return c.RegisterSetting("string", k, s, envName, v, dflt, usage, false, true, true, true)
}

// RegisterStringFlag calls RegisterStringFlagE and ignores any error.
func RegisterStringFlag(k, s, envName, v, dflt, usage string) {
	appCfg.RegisterStringFlag(k, s, envName, v, dflt, usage)
}
func (c *Cfg) RegisterStringFlag(k, s, envName, v, dflt, usage string) {
	c.RegisterStringFlagE(k, s, envName, v, dflt, usage)
}

// RegisterBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterBoolE(k string, v bool) error { return appCfg.RegisterBoolE(k, v) }
func (c *Cfg) RegisterBoolE(k string, v bool) error {
	return c.RegisterSetting("bool", k, "", "", v, strconv.FormatBool(v), "", false, false, false, false)
}

// RegisterBool calls RegisterBoolE and ignores any error.
func RegisterBool(k string, v bool) { appCfg.RegisterBool(k, v) }
func (c *Cfg) RegisterBool(k string, v bool) {
	c.RegisterBoolE(k, v)
}

// RegisterIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterIntE(k string, v int) error { return appCfg.RegisterIntE(k, v) }
func (c *Cfg) RegisterIntE(k string, v int) error {
	return c.RegisterSetting("int", k, "", "", v, strconv.Itoa(v), "", false, false, false, false)
}

// RegisterInt calls RegisterIntE and ignores any error.
func RegisterInt(k string, v int) { appCfg.RegisterInt(k, v) }
func (c *Cfg) RegisterInt(k string, v int) {
	c.RegisterIntE(k, v)
}

// RegisterInt64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterInt64E(k string, v int64) error { return appCfg.RegisterInt64E(k, v) }
func (c *Cfg) RegisterInt64E(k string, v int64) error {
	return c.RegisterSetting("int64", k, "", "", v, strconv.FormatInt(v, 10), "", false, false, false, false)
}

// RegisterInt64 adds the information to the global Cfg, appCfg.

// RegisterInt64 calls RegisterInt64E and ignores any error.
func RegisterInt64(k string, v int64) { appCfg.RegisterInt64(k, v) }
func (c *Cfg) RegisterInt64(k string, v int64) {
	c.RegisterInt64E(k, v)
}

// RegisterStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func RegisterStringE(k, v string) error { return appCfg.RegisterStringE(k, v) }
func (c *Cfg) RegisterStringE(k, v string) error {
	return c.RegisterSetting("string", k, "", "", v, v, "", false, false, false, false)
}

// RegisterString calls RegisterStringE and ignores any error.
func RegisterString(k, v string) { appCfg.RegisterString(k, v) }
func (c *Cfg) RegisterString(k, v string) {
	c.RegisterStringE(k, v)
}
