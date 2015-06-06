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

// RegisterConfigFile set's the configuration file's name.  The name is parsed
// for a valid extension--one that is a supported format--and saves that value
// too. If it cannot be determined, the extension info is not set.  These are
// considered core values and cannot be changed from configuration files,
// environment variables, and configuration files.
//
// If the envName is a non-empty value, it is the environment variable name to
// check for a configuration filename.
func (c *Cfg) RegisterCfgFile(k, envName, v string) error {
	if v == "" {
		return fmt.Errorf("RegisterCfgFile expected a cfg filename: none received")
	}
	if k == "" {
		return fmt.Errorf("RegisterCfgFile expected a configuration key: none received")
	}
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
	return nil
}

// RegisterSetting sets appCfg settings.
func RegisterSetting(typ, name, short, envName string, value interface{}, dflt, usage string, IsCore, IsCfg, IsEnv, IsFlag bool) error {
	return appCfg.RegisterSetting(typ, name, short, envName, value, dflt, usage, IsCore, IsCfg, IsEnv, IsFlag)
}

// RegisterSetting checks to see if the entry already exists and adds the new
// setting if it does not.
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
	c.settings[name] = &setting{
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
		c.cfgVars[name] =  struct{}
	}
	// if it's a IsEnv, make a map for it and the short code
	if envName != "" && IsEnv {
		c.envNames[name] = envName
		if short != "" {
			c.envNames[short] = envName
		}
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

// RegisterBoolCoreE adds the information to the global Cfg, appCfg.
func RegisterBoolCoreE(k string, v bool) error {
	return appCfg.RegisterBoolCoreE(k, v)
}

// RegisterBoolCoreE adds the information to the AppsConfig struct, but does not
// save it to its environment variable. E versions return received errors.
func (c *Cfg) RegisterBoolCoreE(k string, v bool) error {
	return c.RegisterSetting("bool", k, "", "", v, strconv.FormatBool(v), "", true, false, false, false)
}

// RegisterBoolCore adds the information to the global Cfg, appCfg.
func RegisterBoolCore(k string, v bool) {
	appCfg.RegisterBoolCore(k, v)
}

// RegisterBoolCore calls RegisterBoolCoreE and ignores any error.
func (c *Cfg) RegisterBoolCore(k string, v bool) {
	c.RegisterBoolCoreE(k, v)
}

// RegisterIntCoreE adds the information to the global Cfg, appCfg.
func RegisterIntCoreE(k string, v int) error {
	return appCfg.RegisterIntCoreE(k, v)
}

// RegisterIntCoreE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterIntCoreE(k string, v int) error {
	return c.RegisterSetting("int", k, "", "", v, strconv.Itoa(v), "", true, false, false, false)
}

// RegisterIntCore adds the information to the global Cfg, appCfg.
func RegisterIntCore(k string, v int) {
	appCfg.RegisterIntCore(k, v)
}

// RegisterIntCore calls RegisterIntCoreE and ignores any error.
func (c *Cfg) RegisterIntCore(k string, v int) {
	c.RegisterIntCoreE(k, v)
}

// RegisterInt64CoreE adds the information to the global Cfg, appCfg.
func RegisterInt64CoreE(k string, v int64) error {
	return appCfg.RegisterInt64CoreE(k, v)
}

// RegisterInt64CoreE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterInt64CoreE(k string, v int64) error {
	return c.RegisterSetting("int64", k, "", "", v, strconv.FormatInt(v, 10), "", true, false, false, false)
}

// RegisterInt64Core adds the information to the global Cfg, appCfg.
func RegisterInt64Core(k string, v int64) {
	appCfg.RegisterInt64Core(k, v)
}

// RegisterInt64Core calls RegisterInt64CoreE and ignores any error.
func (c *Cfg) RegisterInt64Core(k string, v int64) {
	c.RegisterInt64CoreE(k, v)
}

// RegisterStringCoreE adds the information to the global Cfg, appCfg.
func RegisterStringCoreE(k, v string) error {
	return appCfg.RegisterStringCoreE(k, v)
}

// RegisterStringCoreE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterStringCoreE(k, v string) error {
	return c.RegisterSetting("string", k, "", "", v, v, "", true, false, false, false)
}

// RegisterStringCore adds the information to the global Cfg, appCfg.
func RegisterStringCore(k, v string) {
	appCfg.RegisterStringCore(k, v)
}

// RegisterStringCore calls RegisterStringCoreE and ignores any error.
func (c *Cfg) RegisterStringCore(k, v string) {
	c.RegisterStringCoreE(k, v)
}

// RegisterBoolCfgE adds the information to the global Cfg, appCfg.
func RegisterBoolCfgE(k, envName string, v bool) error {
	return appCfg.RegisterBoolCfgE(k, envName, v)
}

// RegisterBoolCfgE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterBoolCfgE(k, envName string, v bool) error {
	return c.RegisterSetting("bool", k, "", envName, v, strconv.FormatBool(v), "", false, true, true, false)
}

// RegisterBoolCfg adds the information to the global Cfg, appCfg.
func RegisterBoolCfg(k, envName string, v bool) {
	appCfg.RegisterBoolCfg(k, envName, v)
}

// RegisterBoolCfg calls RegisterBoolCfgE and ignores any error.
func (c *Cfg) RegisterBoolCfg(k, envName string, v bool) {
	c.RegisterBoolCfgE(k, envName, v)
}

// RegisterIntCfgE adds the information to the global Cfg, appCfg.
func RegisterIntCfgE(k, envName string, v int) error {
	return appCfg.RegisterIntCfgE(k, envName, v)
}

// RegisterIntCfgE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterIntCfgE(k, envName string, v int) error {
	return c.RegisterSetting("int", k, "", envName, v, strconv.Itoa(v), "", false, true, true, false)
}

// RegisterIntCfg adds the information to the global Cfg, appCfg.
func RegisterIntCfg(k, envName string, v int) {
	appCfg.RegisterIntCfg(k, envName, v)
}

// RegisterIntCfg calls RegisterIntCfgE and ignores any error.
func (c *Cfg) RegisterIntCfg(k, envName string, v int) {
	c.RegisterIntCfgE(k, envName, v)
}

// RegisterInt64CfgE adds the information to the global Cfg, appCfg.
func RegisterInt64CfgE(k, envName string, v int64) error {
	return appCfg.RegisterInt64CfgE(k, envName, v)
}

// RegisterInt64Cfg adds the informatio to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterInt64CfgE(k, envName string, v int64) error {
	return c.RegisterSetting("int64", k, "", envName, v, strconv.FormatInt(v, 10), "", false, true, true, false)
}

// RegisterInt64Cfg adds the information to the global Cfg, appCfg.
func RegisterInt64Cfg(k, envName string, v int64) {
	appCfg.RegisterInt64Cfg(k, envName, v)
}

// RegisterInt64Cfg calls RegisterInt64Cfg and ignores any error.
func (c *Cfg) RegisterInt64Cfg(k, envName string, v int64) {
	c.RegisterInt64CfgE(k, envName, v)
}

// RegisterStringCfgE adds the information to the global Cfg, appCfg.
func RegisterStringCfgE(k, envName, v string) error {
	return appCfg.RegisterStringCfgE(k, envName, v)
}

// RegisterStringCfgE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterStringCfgE(k, envName, v string) error {
	return c.RegisterSetting("string", k, "", envName, v, v, "", false, true, true, false)
}

// RegisterStringCfg adds the information to the global Cfg, appCfg.
func RegisterStringCfg(k, envName, v string) {
	appCfg.RegisterStringCfg(k, envName, v)
}

// RegisterStringCfg calls RegisterStringCfgE and ignores any error.
func (c *Cfg) RegisterStringCfg(k, envName, v string) {
	c.RegisterStringCfgE(k, envName, v)
}

// RegisterBoolFlagE adds the information to the global Cfg, appCfg.
func RegisterBoolFlagE(k, s, envName string, v bool, dflt, usage string) error {
	return appCfg.RegisterBoolFlagE(k, s, envName, v, dflt, usage)
}

// RegisterBoolFlagE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterBoolFlagE(k, s, envName string, v bool, dflt, usage string) error {
	return c.RegisterSetting("bool", k, s, envName, v, dflt, usage, false, true, true, true)
}

// RegisterBoolFlag adds the information to the global Cfg, appCfg.
func RegisterBoolFlag(k, s, envName string, v bool, dflt, usage string) {
	appCfg.RegisterBoolFlag(k, s, envName, v, dflt, usage)
}

// RegisterBoolFlag calls RegisterBoolFlagE and ignores any error.
func (c *Cfg) RegisterBoolFlag(k, s, envName string, v bool, dflt, usage string) {
	c.RegisterBoolFlagE(k, s, envName, v, dflt, usage)
}

// RegisterIntFlagE adds the information to the global Cfg, appCfg.
func RegisterIntFlagE(k, s, envName string, v int, dflt, usage string) error {
	return appCfg.RegisterIntFlagE(k, s, envName, v, dflt, usage)
}

// RegisterIntFlagE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterIntFlagE(k, s, envName string, v int, dflt, usage string) error {
	return c.RegisterSetting("int", k, s, envName, v, dflt, usage, false, true, true, true)
}

// RegisterIntFlag adds the information to the global Cfg, appCfg.
func RegisterIntFlag(k, s, envName string, v int, dflt, usage string) {
	appCfg.RegisterIntFlag(k, s, envName, v, dflt, usage)
}

// RegisterIntFlag calls RegisterIntFlagE and ignores any error.
func (c *Cfg) RegisterIntFlag(k, s, envName string, v int, dflt, usage string) {
	c.RegisterIntFlagE(k, s, envName, v, dflt, usage)
}

// RegisterInt64FlagE adds the information to the global Cfg, appCfg.
func RegisterInt64FlagE(k, s, envName string, v int64, dflt, usage string) error {
	return appCfg.RegisterInt64FlagE(k, s, envName, v, dflt, usage)
}

// RegisterInt64FlagE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterInt64FlagE(k, s, envName string, v int64, dflt, usage string) error {
	return c.RegisterSetting("int64", k, s, envName, v, dflt, usage, false, true, true, true)
}

// RegisterInt64Flag adds the information to the global Cfg, appCfg.
func RegisterInt64Flag(k, s, envName string, v int64, dflt, usage string) {
	appCfg.RegisterInt64Flag(k, s, envName, v, dflt, usage)
}

// RegisterInt64Flag calls RegisterIntFlagE and ignores any error.
func (c *Cfg) RegisterInt64Flag(k, s, envName string, v int64, dflt, usage string) {
	c.RegisterInt64FlagE(k, s, envName, v, dflt, usage)
}

// RegisterStringFlagE adds the information to the global Cfg, appCfg.
func RegisterStringFlagE(k, s, envName, v, dflt, usage string) error {
	return appCfg.RegisterStringFlagE(k, s, envName, v, dflt, usage)
}

// RegisterStringFlagE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterStringFlagE(k, s, envName, v, dflt, usage string) error {
	return c.RegisterSetting("string", k, s, envName, v, dflt, usage, false, true, true, true)
}

// RegisterStringFlag adds the information to the global Cfg, appCfg.
func RegisterStringFlag(k, s, envName, v, dflt, usage string) {
	appCfg.RegisterStringFlag(k, s, envName, v, dflt, usage)
}

// RegisterStringFlag calls RegisterStringFlagE and ignores any error.
func (c *Cfg) RegisterStringFlag(k, s, envName, v, dflt, usage string) {
	c.RegisterStringFlagE(k, s, envName, v, dflt, usage)
}

// RegisterBoolE adds the information to the global Cfg, appCfg.
func RegisterBoolE(k string, v bool) error {
	return appCfg.RegisterBoolE(k, v)
}

// RegisterBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterBoolE(k string, v bool) error {
	return c.RegisterSetting("bool", k, "", "", v, strconv.FormatBool(v), "", false, false, false, false)
}

// RegisterBool adds the information to the global Cfg, appCfg.
func RegisterBool(k string, v bool) {
	appCfg.RegisterBool(k, v)
}

// RegisterBool calls RegisterBoolE and ignores any error.
func (c *Cfg) RegisterBool(k string, v bool) {
	c.RegisterBoolE(k, v)
}

// RegisterIntE adds the information to the global Cfg, appCfg.
func RegisterIntE(k string, v int) error {
	return appCfg.RegisterIntE(k, v)
}

// RegisterIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterIntE(k string, v int) error {
	return c.RegisterSetting("int", k, "", "", v, strconv.Itoa(v), "", false, false, false, false)
}

// RegisterInt adds the information to the global Cfg, appCfg.
func RegisterInt(k string, v int) {
	appCfg.RegisterInt(k, v)
}

// RegisterInt calls RegisterIntE and ignores any error.
func (c *Cfg) RegisterInt(k string, v int) {
	c.RegisterIntE(k, v)
}

// RegisterInt64E adds the information to the global Cfg, appCfg.
func RegisterInt64E(k string, v int64) error {
	return appCfg.RegisterInt64E(k, v)
}

// RegisterInt64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterInt64E(k string, v int64) error {
	return c.RegisterSetting("int64", k, "", "", v, strconv.FormatInt(v, 10), "", false, false, false, false)
}

// RegisterInt64 adds the information to the global Cfg, appCfg.
func RegisterInt64(k string, v int64) {
	appCfg.RegisterInt64(k, v)
}

// RegisterInt64 calls RegisterInt64E and ignores any error.
func (c *Cfg) RegisterInt64(k string, v int64) {
	c.RegisterInt64E(k, v)
}

// RegisterStringE adds the information to the global Cfg, appCfg.
func RegisterStringE(k, v string) error {
	return appCfg.RegisterStringE(k, v)
}

// RegisterStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable: E versions return received errors.
func (c *Cfg) RegisterStringE(k, v string) error {
	return c.RegisterSetting("string", k, "", "", v, v, "", false, false, false, false)
}

// RegisterString adds the information to the global Cfg, appCfg.
func RegisterString(k, v string) {
	appCfg.RegisterString(k, v)
}

// RegisterString calls RegisterStringE and ignores any error.
func (c *Cfg) RegisterString(k, v string) {
	c.RegisterStringE(k, v)
}

//  is a convenience function for the appCfg global config.
func RegisterCfgFile(k, envName, v string) error {
	return appCfg.RegisterCfgFile(k, envName, v)
}
