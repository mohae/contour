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
)


// Config methods
// RegisterConfigFilename set's the configuration file's name. The name is
// parsed for a valid extension--one that is a supported format--and saves
// that value too. If it cannot be determined, the extension info is not set.
// These are considered core values and cannot be changed from command-line
// and configuration files. (IsCore == true).
func (c *Cfg) RegisterConfigFilename(k, v string) error {
	if v == "" {
		return fmt.Errorf("A config filename was expected, none received")
	}

	if k == "" {
		return fmt.Errorf("A key for the config filename setting was expected, none received")
	}

	c.RegisterCoreString(k, v)

	// Register it first. If a valid config format isn't found, an error
	// will be returned, so registering it afterwords would mean the
	// setting would not exist.
	c.RegisterString(EnvCfgFormat, "")
	format, err := configFormat(v)
	if err != nil {
		return err
	}

	c.Settings[EnvCfgFormat].Value = format

	// Now we can update the format, since it wasn't set before, it can be
	// set now before it becomes read only.
	c.RegisterString(EnvCfgFormat, format)

	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func (c *Cfg) RegisterSetting(Type string, k string, v interface{}, Code string, IsCore, IsCfg, IsFlag bool) {
	var update bool
	_, ok := c.Settings[k]
	if ok {
		// Core Settings can't be re-registered.
		if c.Settings[k].IsCore {
			return
		}

		if c.Settings[k].Value != nil {
			return
		}

		update = true
	}

	if update {
		c.Settings[k].Type = Type
		c.Settings[k].Value = v
		c.Settings[k].Code = Code
		c.Settings[k].IsCore = IsCore
		c.Settings[k].IsCfg = IsCfg
		c.Settings[k].IsFlag = IsFlag
		return
	}

	c.Settings[k] = &setting{
		Type:      Type,
		Value:     v,
		Code:      Code,
		IsCore:    IsCore,
		IsCfg:     IsCfg,
		IsFlag:    IsFlag,
	}
}

// RegisterCoreBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func (c *Cfg) RegisterCoreBool(k string, v bool) {
	c.RegisterSetting("bool", k, v, "", true, false, false)
	return
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func (c *Cfg) RegisterCoreInt(k string, v int) {
	c.RegisterSetting("int", k, v, "", true, false, false)
	return
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func (c *Cfg) RegisterCoreString(k, v string) {
	c.RegisterSetting("string", k, v, "", true, false, false)
	return
}

// RegisterConfBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterConfBool(k string, v bool) {
	c.RegisterSetting("bool", k, v, "", false, true, false)
	return
}

// RegisterConfInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterConfInt(k string, v bool) {
	c.RegisterSetting("int", k, v, "", false, true, false)
	return
}

// RegisterConfString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterConfString(k string, v bool) {
	c.RegisterSetting("string", k, v, "", false, true, false)
	return
}

// RegisterFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterFlagBool(k string, v bool, f string) {
	c.RegisterSetting("bool", k, v, f, false, true, true)
	return
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterFlagInt(k string, v int, f string) {
	c.RegisterSetting("int", k, v, f, false, true, true)
	return
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterFlagString(k, v, f string) {
	c.RegisterSetting("string", k, v, f, false, true, true)
	return
}


// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterBool(k string, v bool) {
	c.RegisterSetting("bool", k, v, "", false, false, false)
	return
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterInt(k string, v int) {
	c.RegisterSetting("int", k, v, "", false, false, false)
	return
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterString(k, v string) {
	c.RegisterSetting("string", k, v, "", false, false, false)
	return
}

// Convenience functions for interacting with the configs[app] configuration.

// RegisterConfigFilename set's the configuration file's name. The name is
// parsed for a valid extension--one that is a supported format--and saves
// that value too. If it cannot be determined, the extension info is not set.
// These are considered core values and cannot be changed from command-line
// and configuration files. (IsCore == true).
func RegisterConfigFilename(k, v string) error {
	if v == "" {
		return fmt.Errorf("A config filename was expected, none received")
	}

	if k == "" {
		return fmt.Errorf("A key for the config filename setting was expected, none received")
	}

	configs[app].RegisterCoreString(k, v)

	// TODO redo this given new paradigm
	// Register it first. If a valid config format isn't found, an error
	// will be returned, so registering it afterwords would mean the
	// setting would not exist.
	configs[app].RegisterString(EnvCfgFormat, "")
	format, err := configFormat(v)
	if err != nil {
		return err
	}

	configs[app].Settings[EnvCfgFormat].Value = format

	configs[app].RegisterString(EnvCfgFormat, format)

	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func RegisterSetting(Type string, k string, v interface{}, Code string, IsCore, IsCfg, IsFlag bool) {
	var update bool
	_, ok := configs[app].Settings[k]
	if ok {
		// Core Settings can't be re-registered.
		if configs[app].Settings[k].IsCore {
			return
		}

		if configs[app].Settings[k].Value != nil {
			return
		}

		update = true
	}

	if update {
		configs[app].Settings[k].Type = Type
		configs[app].Settings[k].Value = v
		configs[app].Settings[k].Code = Code
		configs[app].Settings[k].IsCore = IsCore
		configs[app].Settings[k].IsCfg = IsCfg
		configs[app].Settings[k].IsFlag = IsFlag
		return
	}

	configs[app].Settings[k] = &setting{
		Type:      Type,
		Value:     v,
		Code:      Code,
		IsCore:    IsCore,
		IsCfg:     IsCfg,
		IsFlag:    IsFlag,
	}
}

// RegisterCoreBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreBool(k string, v bool) {
	configs[app].RegisterSetting("bool", k, v, "", true, false, false)
	return
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreInt(k string, v int) {
	configs[app].RegisterSetting("int", k, v, "", true, false, false)
	return
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its environment variable
func RegisterCoreString(k, v string) {
	configs[app].RegisterSetting("string", k, v, "", true, false, false)
	return
}

// RegisterConfBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterConfBool(k string, v bool) {
	configs[app].RegisterSetting("bool", k, v, "", false, true, false)
	return
}

// RegisterConfInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterConfInt(k string, v bool) {
	configs[app].RegisterSetting("int", k, v, "", false, true, false)
	return
}

// RegisterConfString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterConfString(k string, v bool) {
	configs[app].RegisterSetting("string", k, v, "", false, true, false)
	return
}

// RegisterFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterFlagBool(k string, v bool, f string) {
	configs[app].RegisterSetting("bool", k, v, f, false, true, true)
	return
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterFlagInt(k string, v int, f string) {
	configs[app].RegisterSetting("int", k, v, f, false, true, true)
	return
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterFlagString(k, v, f string) {
	configs[app].RegisterSetting("string", k, v, f, false, true, true)
	return
}


// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterBool(k string, v bool) {
	configs[app].RegisterSetting("bool", k, v, "", false, false, false)
	return
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterInt(k string, v int) {
	configs[app].RegisterSetting("int", k, v, "", false, false, false)
	return
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func RegisterString(k, v string) {
	configs[app].RegisterSetting("string", k, v, "", false, false, false)
	return
}
