package contour

// Register contains all of contour's Register functions.Calling Register
// adds, or registers, the Settings information to the AppConfig variable.
// The setting value, if there is one, is not saved to its ironment
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

// Config methods
// RegisterConfigFilename set's the configuration file's name. The name is
// parsed for a valid extension--one that is a supported format--and saves
// that value too. If it cannot be determined, the extension info is not set.
// These are considered core values and cannot be changed from command-line
// and configuration files. (IsCore == true).
func (c *Cfg) RegisterCfgFilename(k, v string) error {
	if v == "" {
		return fmt.Errorf("A config filename was expected, none received")
	}

	if k == "" {
		return fmt.Errorf("A key for the config filename setting was expected, none received")
	}

	c.RegisterStringCore(k, v)

	// Register it first. If a valid config format isn't found, an error
	// will be returned, so registering it afterwords would mean the
	// setting would not exist.
	c.RegisterString(CfgFormat, "")
	format, err := formatFromFilename(v)
	if err != nil {
		return err
	}

	// Now we can update the format, since it wasn't set before, it can be
	// set now before it becomes read only.
	c.UpdateString(CfgFormat, format.String())
	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func (c *Cfg) RegisterSetting(typ, name, short string, value interface{}, dflt string, usage string, IsCore, IsCfg, IsFlag bool) {
	c.lock.RLock()
	_, ok := c.settings[name]
	if ok {
		// Settings can't be re-registered.
		c.lock.RUnlock()
		return
	}

	c.lock.RUnlock()
	c.lock.Lock()
	defer c.lock.Unlock()

	// Add the setting
	c.settings[name] = &setting{
		Type:    typ,
		Name:    name,
		Short:   short,
		Value:   value,
		Default: dflt,
		Usage:   usage,
		IsCore:  IsCore,
		IsCfg:   IsCfg,
		IsFlag:  IsFlag,
	}

	// Keep track of whether or not a config is being used. If a setting is
	// registered as a config setting, it is assumed a configuration source
	// is being used.
	if IsCfg {
		c.useCfg = true
	}

	// Keep track of whether or not flags are being used. If a setting is
	// registered as a flag setting, it is assumed that flags are being
	// used.
	if IsFlag {
		c.useFlags = true
	}
}

// RegisterBoolCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterBoolCore(k, v string) {
	if v != "" {
		_, err := strconv.ParseBool(v)
		if err != nil {
			v = "" // if the parse error'd set to "", or unset
		}
	}
	c.RegisterSetting("bool", k, "", v, v, "", true, false, false)
	return
}

// RegisterIntCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterIntCore(k string, v int) {
	c.RegisterSetting("int", k, "", v, strconv.Itoa(v), "", true, false, false)
	return
}

// RegisterInt64Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterInt64Core(k string, v int64) {
	c.RegisterSetting("int64", k, "", v, strconv.FormatInt(v, 10), "", true, false, false)
	return
}

// RegisterStringCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterStringCore(k, v string) {
	c.RegisterSetting("string", k, "", v, v, "", true, false, false)
	return
}

// RegisterBoolConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBoolConf(k, v string) {
	if v != "" {
		_, err := strconv.ParseBool(v)
		if err != nil {
			v = "" // if parse results in error, don't set
		}
	}
	c.RegisterSetting("bool", k, "", v, v, "", false, true, false)
	return
}

// RegisterIntConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterIntConf(k string, v int) {
	c.RegisterSetting("int", k, "", v, strconv.Itoa(v), "", false, true, false)
	return
}

// RegisterInt64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt64Conf(k string, v int64) {
	c.RegisterSetting("int64", k, "", v, strconv.FormatInt(v, 10), "", false, true, false)
	return
}

// RegisterStringConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterStringConf(k string, v string) {
	c.RegisterSetting("string", k, "", v, v, "", false, true, false)
	return
}

// RegisterBoolFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBoolFlag(k, s, v, dflt, usage string) {
	c.RegisterSetting("bool", k, s, v, dflt, usage, false, true, true)
	return
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterIntFlag(k, s string, v int, dflt, usage string) {
	c.RegisterSetting("int", k, s, v, dflt, usage, false, true, true)
	return
}

// RegisterInt64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt64Flag(k, s string, v int64, dflt, usage string) {
	c.RegisterSetting("int64", k, s, v, dflt, usage, false, true, true)
	return
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterStringFlag(k, s, v, dflt, usage string) {
	c.RegisterSetting("string", k, s, v, dflt, usage, false, true, true)
	return
}

// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBool(k, v string) {
	if v != "" {
		_, err := strconv.ParseBool(v)
		if err != nil {
			v = "" // parse error == unset
		}
	}
	c.RegisterSetting("bool", k, "", v, v, "", false, false, false)
	return
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt(k string, v int) {
	c.RegisterSetting("int", k, "", v, strconv.Itoa(v), "", false, false, false)
	return
}

// RegisterInt64 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt64(k string, v int64) {
	c.RegisterSetting("int64", k, "", v, strconv.FormatInt(v, 10), "", false, false, false)
	return
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterString(k, v string) {
	c.RegisterSetting("string", k, "", v, v, "", false, false, false)
	return
}

// Convenience functions for interacting with the configs[app] configuration.

// RegisterCfgFilename set's the configuration file's name. The name is
// parsed for a valid extension--one that is a supported format--and saves
// that value too. If it cannot be determined, the extension info is not set.
// These are considered core values and cannot be changed from command-line
// and configuration files. (IsCore == true).
func RegisterCfgFilename(k, v string) error {
	return appCfg.RegisterCfgFilename(k, v)
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func RegisterSetting(typ, name, short string, value interface{}, dflt, usage string, IsCore, IsCfg, IsFlag bool) {
	appCfg.RegisterSetting(typ, name, short, value, dflt, usage, IsCore, IsCfg, IsFlag)
}

// RegisterBoolCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterBoolCore(k, v string) {
	appCfg.RegisterBoolCore(k, v)
}

// RegisterIntCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterIntCore(k string, v int) {
	appCfg.RegisterIntCore(k, v)
}

// RegisterInt64Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterInt64Core(k string, v int64) {
	appCfg.RegisterInt64Core(k, v)
}

// RegisterStringCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterStringCore(k, v string) {
	appCfg.RegisterStringCore(k, v)
}

// RegisterConfBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisteeBoolCore(k, v string) {
	appCfg.RegisterBoolCore(k, v)
}

// RegisterIntConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterIntConf(k string, v int) {
	appCfg.RegisterIntConf(k, v)
}

// RegisterInt64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64Conf(k string, v int64) {
	appCfg.RegisterInt64Conf(k, v)
}

// RegisterStringConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterStringConf(k, v string) {
	appCfg.RegisterStringConf(k, v)
}

// RegisterBoolFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterBoolFlag(k, s, v, dflt, u string) {
	appCfg.RegisterBoolFlag(k, s, v, dflt, u)
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterIntFlag(k, s string, v int, dflt, u string) {
	appCfg.RegisterIntFlag(k, s, v, dflt, u)
}

// RegisterInt64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64Flag(k, s string, v int64, dflt, u string) {
	appCfg.RegisterInt64Flag(k, s, v, dflt, u)
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterStringFlag(k, s, v, dflt, u string) {
	appCfg.RegisterStringFlag(k, s, v, dflt, u)
}

// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterBool(k, v string) {
	appCfg.RegisterBool(k, v)
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt(k string, v int) {
	appCfg.RegisterInt(k, v)
}

// RegisterInt64 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64(k string, v int64) {
	appCfg.RegisterInt64(k, v)
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterString(k, v string) {
	appCfg.RegisterString(k, v)
}
