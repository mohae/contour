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
	c.RegisterString(CfgFormat, "")
	format, err := configFormat(v)
	if err != nil {
		return err
	}

	// Now we can update the format, since it wasn't set before, it can be
	// set now before it becomes read only.
	c.RegisterString(CfgFormat, format.String())

	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func (c *Cfg) RegisterSetting(Type string, k string, v interface{}, Code string, IsCore, IsCfg, IsFlag bool) {
	var update bool
	c.Lock.RLock()
	_, ok := configs[app].settings[k]
	if ok {
		c.Lock.RUnlock()
		// Core Settings can't be re-registered.
		if configs[app].settings[k].IsCore {
			return
		}

		if configs[app].settings[k].Value != nil {
			return
		}

		update = true
	}

	c.Lock.RLock()
	c.Lock.Lock()
	defer c.Lock.Lock()

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

	// If the registered setting is allowed to be updated:
	if update {
		c.settings[k].Type = Type
		c.settings[k].Value = v
		c.settings[k].Code = Code
		c.settings[k].IsCore = IsCore
		c.settings[k].IsCfg = IsCfg
		c.settings[k].IsFlag = IsFlag
		return
	}

	// Otherwise register it as a new setting.
	configs[app].settings[k] = &setting{
		Type:      Type,
		Value:     v,
		Code:      Code,
		IsCore:    IsCore,
		IsCfg:     IsCfg,
		IsFlag:    IsFlag,
	}
}

// RegisterCoreBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterCoreBool(k string, v bool) {
	c.RegisterSetting("bool", k, v, "", true, false, false)
	return
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterCoreInt(k string, v int) {
	c.RegisterSetting("int", k, v, "", true, false, false)
	return
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterCoreString(k, v string) {
	c.RegisterSetting("string", k, v, "", true, false, false)
	return
}

// RegisterConfBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterConfBool(k string, v bool) {
	c.RegisterSetting("bool", k, v, "", false, true, false)
	return
}

// RegisterConfInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterConfInt(k string, v bool) {
	c.RegisterSetting("int", k, v, "", false, true, false)
	return
}

// RegisterConfString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterConfString(k string, v bool) {
	c.RegisterSetting("string", k, v, "", false, true, false)
	return
}

// RegisterFlagBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFlagBool(k string, v bool, f string) {
	c.RegisterSetting("bool", k, v, f, false, true, true)
	return
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFlagInt(k string, v int, f string) {
	c.RegisterSetting("int", k, v, f, false, true, true)
	return
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFlagString(k, v, f string) {
	c.RegisterSetting("string", k, v, f, false, true, true)
	return
}


// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBool(k string, v bool) {
	c.RegisterSetting("bool", k, v, "", false, false, false)
	return
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt(k string, v int) {
	c.RegisterSetting("int", k, v, "", false, false, false)
	return
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
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
	configs[app].RegisterString(CfgFormat, "")
	format, err := configFormat(v)
	if err != nil {
		return err
	}

	configs[app].RegisterString(CfgFormat, format.String())

	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func RegisterSetting(Type string, k string, v interface{}, Code string, IsCore, IsCfg, IsFlag bool) {
	configs[app].RegisterSetting(Type, k, v, Code, IsCore, IsCfg, IsFlag)
}

// RegisterCoreBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreBool(k string, v bool) {
	configs[app].RegisterSetting("bool", k, v, "", true, false, false)
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreInt(k string, v int) {
	configs[app].RegisterSetting("int", k, v, "", true, false, false)
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreString(k, v string) {
	configs[app].RegisterSetting("string", k, v, "", true, false, false)
}

// RegisterConfBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterConfBool(k string, v bool) {
	configs[app].RegisterSetting("bool", k, v, "", false, true, false)
}

// RegisterConfInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterConfInt(k string, v bool) {
	configs[app].RegisterSetting("int", k, v, "", false, true, false)
}

// RegisterConfString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterConfString(k string, v bool) {
	configs[app].RegisterSetting("string", k, v, "", false, true, false)
}

// RegisterFlagBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFlagBool(k string, v bool, f string) {
	configs[app].RegisterSetting("bool", k, v, f, false, true, true)
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFlagInt(k string, v int, f string) {
	configs[app].RegisterSetting("int", k, v, f, false, true, true)
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFlagString(k, v, f string) {
	configs[app].RegisterSetting("string", k, v, f, false, true, true)
}


// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterBool(k string, v bool) {
	configs[app].RegisterSetting("bool", k, v, "", false, false, false)
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt(k string, v int) {
	configs[app].RegisterSetting("int", k, v, "", false, false, false)
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterString(k, v string) {
	configs[app].RegisterSetting("string", k, v, "", false, false, false)
}
