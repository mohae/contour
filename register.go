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

	c.RegisterStringCore(k, v)

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
func (c *Cfg) RegisterSetting(Type, k, Short string, v interface{}, Usage string, Default string, IsCore, IsCfg, IsFlag bool) error  {

	// find the index for the setting. If its found, return an error
	idx, err := c.settingIndex(k)
	if err == nil {
		err := fmt.Errorf("cannot register %q, it is already registered")
		logger.Error(err)
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()

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

	// Otherwise register it as a new setting.
	c.settings[c.settingsCount] = &setting{
		Type:      Type,
		Value:     v,
		Short:     Short,
		Usage:     Usage,
		Default:   Default,
		IsCore:    IsCore,
		IsCfg:     IsCfg,
		IsFlag:    IsFlag,
	}
	
	c.settingsCount++
	return nil
}

// RegisterBoolCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterBoolCore(k string, v bool) {
	c.RegisterSetting("bool", k, "", v, "", "", true, false, false)
}

// RegisterFloat32Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterFloat32Core(k string, v int) {
	c.RegisterSetting("float32", k, "", v, "", "", true, false, false)
}

// RegisterFloat64Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterFloat64Core(k string, v int) {
	c.RegisterSetting("float64", k, "", v, "", "", true, false, false)
}

// RegisterIntCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterIntCore(k string, v int) {
	c.RegisterSetting("int", k, "", v, "", "", true, false, false)
}

// RegisterInt8Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterInt8Core(k string, v int) {
	c.RegisterSetting("int8", k, "", v, "", "", true, false, false)
}

// RegisterInt32Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterInt32Core(k string, v int) {
	c.RegisterSetting("int32", k, "", v, "", "", true, false, false)
}

// RegisterInt64Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterInt64Core(k string, v int) {
	c.RegisterSetting("int64", k, "", v, "", "", true, false, false)
}

// RegisterStringCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterStringCore(k, v string) {
	c.RegisterSetting("string", k, "", v, "", "", true, false, false)
}

// RegisterBoolConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBoolConf(k string, v bool) {
	c.RegisterSetting("bool", k, "", v, "", "", false, true, false)
}

// RegisterFloat32Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFloat32Conf(k string, v bool) {
	c.RegisterSetting("float32", k, "", v, "", "", false, true, false)
}

// RegisterFloat64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFloat64Conf(k string, v bool) {
	c.RegisterSetting("float64", k, "", v, "", "", false, true, false)
}

// RegisterIntConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterIntConf(k string, v bool) {
	c.RegisterSetting("int", k, "", v, "", "", false, true, false)
}

// RegisterInt8Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt8Conf(k string, v bool) {
	c.RegisterSetting("int8", k, "", v, "", "", false, true, false)
}

// RegisterInt32Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt32Conf(k string, v bool) {
	c.RegisterSetting("int32", k, "", v, "", "", false, true, false)
}

// RegisterInt64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt64Conf(k string, v bool) {
	c.RegisterSetting("int64", k, "", v, "", "", false, true, false)
}

// RegisterStringConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterStringConf(k string, v bool) {
	c.RegisterSetting("string", k, "", v, "", "", false, true, false)
}

// RegisterBoolFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBoolFlag(k, s string, v bool) {
	c.RegisterSetting("bool", k, s, v, "", "", false, true, true)
}

// RegisterFloat32Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFLoat32Flag(k, s string, v int) {
	c.RegisterSetting("float32", k, s, v, "", "", false, true, true)
}

// RegisterFloat64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFloat64Flag(k, s string, v int) {
	c.RegisterSetting("float64", k, s, v, "", "", false, true, true)
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterIntFlag(k, s string, v int) {
	c.RegisterSetting("int", k, s, v, "", "", false, true, true)
}

// RegisterInt8Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt8Flag(k, s string, v int) {
	c.RegisterSetting("int8", k, s, v, "", "", false, true, true)
}

// RegisterInt32Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt32Flag(k, s string, v int) {
	c.RegisterSetting("int32", k, s, v, "", "", false, true, true)
}

// RegisterInt64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt64Flag(k, s string, v int) {
	c.RegisterSetting("int64", k, s, v, "", "", false, true, true)
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterStringFlag(k, s, v string) {
	c.RegisterSetting("string", k, s, v, "", "", false, true, true)
}


// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBool(k string, v bool) {
	c.RegisterSetting("bool", k, "", v, "",  "", false, false, false)
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterInt(k string, v int) {
	c.RegisterSetting("int", k, "", v, "",  "", false, false, false)
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterString(k, v string) {
	c.RegisterSetting("string", k, "", v, "",  "", false, false, false)
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

	configs[0].RegisterStringCore(k, v)

	// TODO redo this given new paradigm
	// Register it first. If a valid config format isn't found, an error
	// will be returned, so registering it afterwords would mean the
	// setting would not exist.
	configs[0].RegisterString(CfgFormat, "")
	format, err := configFormat(v)
	if err != nil {
		return err
	}

	configs[0].RegisterString(CfgFormat, format.String())

	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func RegisterSetting(Type string, k string, short string, v interface{}, Usage string, Default string, IsCore, IsCfg, IsFlag bool) {
	configs[0].RegisterSetting(Type, k, short, v, Usage, Default, IsCore, IsCfg, IsFlag)
}

// RegisterCoreBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreBool(k string, v bool) {
	configs[0].RegisterSetting("bool", k, "", v, "", "", true, false, false)
}

// RegisterCoreFloat32 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreFloat32(k string, v int) {
	configs[0].RegisterSetting("float", k, "", v, "", "", true, false, false)
}

// RegisterCoreFloat64 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreFloat64(k string, v int) {
	configs[0].RegisterSetting("float", k, "", v, "", "", true, false, false)
}

// RegisterCoreInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreInt(k string, v int) {
	configs[0].RegisterSetting("int", k, "", v, "", "", true, false, false)
}

// RegisterCoreInt8 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreInt8(k string, v int) {
	configs[0].RegisterSetting("int8", k, "", v, "", "", true, false, false)
}

// RegisterCoreInt32 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreInt32(k string, v int) {
	configs[0].RegisterSetting("int32", k, "", v, "", "", true, false, false)
}

// RegisterCoreInt64 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreInt64(k string, v int) {
	configs[0].RegisterSetting("int64", k, "", v, "", "", true, false, false)
}

// RegisterCoreString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterCoreString(k, v string) {
	configs[0].RegisterSetting("string", k, "", v, "", "", true, false, false)
}

// RegisterConfBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterConfBool(k string, v bool) {
	configs[0].RegisterSetting("bool", k, "", v, "", "", false, true, false)
}

// RegisterFloat32Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat32Conf(k string, v bool) {
	configs[0].RegisterSetting("float", k, "", v, "", "", false, true, false)
}

// RegisterFloat64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat64Conf(k string, v bool) {
	configs[0].RegisterSetting("float", k, "", v, "", "", false, true, false)
}

// RegisterIntConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterIntConf(k string, v bool) {
	configs[0].RegisterSetting("int", k, "", v, "", "", false, true, false)
}

// RegisterInt8Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt8Conf(k string, v bool) {
	configs[0].RegisterSetting("int8", k, "", v, "", "", false, true, false)
}

// RegisterInt32Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt32Conf(k string, v bool) {
	configs[0].RegisterSetting("int32", k, "", v, "", "", false, true, false)
}

// RegisterInt64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64Conf(k string, v bool) {
	configs[0].RegisterSetting("int64", k, "", v, "", "", false, true, false)
}

// RegisterStringConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterStringConf(k string, v bool) {
	configs[0].RegisterSetting("string", k, "", v, "", "", false, true, false)
}

// RegisterBoolFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterBoolFlag(k, s string, v bool) {
	configs[0].RegisterSetting("bool", k, s, v, "", "", false, true, true)
}

// RegisterFloat32Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat32Flag(k, s string, v int) {
	configs[0].RegisterSetting("float32", k, s, v, "", "", false, true, true)
}

// RegisterFloat64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat64Flag(k, s string, v int) {
	configs[0].RegisterSetting("float64", k, s, v, "", "", false, true, true)
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterIntFlag(k, s string, v int) {
	configs[0].RegisterSetting("int", k, s, v, "", "", false, true, true)
}

// RegisterInt8Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt8Flag(k, s string, v int) {
	configs[0].RegisterSetting("int8", k, s, v, "", "", false, true, true)
}

// RegisterInt32Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt32Flag(k, s string, v int) {
	configs[0].RegisterSetting("int32", k, s, v, "", "", false, true, true)
}

// RegisterInt64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64Flag(k, s string, v int) {
	configs[0].RegisterSetting("int64", k, s, v, "", "", false, true, true)
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterStringFlag(k, s, v string) {
	configs[0].RegisterSetting("string", k, s, v, "", "", false, true, true)
}


// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterBool(k string, v bool) {
	configs[0].RegisterSetting("bool", k, "", v, "", "", false, false, false)
}

// RegisterFloat32 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat32(k string, v int) {
	configs[0].RegisterSetting("float32", k, "", v, "", "", false, false, false)
}

// RegisterFloat64 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat64(k string, v int) {
	configs[0].RegisterSetting("float64", k, "", v, "", "", false, false, false)
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt(k string, v int) {
	configs[0].RegisterSetting("int", k, "", v, "", "", false, false, false)
}

// RegisterInt8 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt8(k string, v int) {
	configs[0].RegisterSetting("int8", k, "", v, "", "", false, false, false)
}

// RegisterInt32 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt32(k string, v int) {
	configs[0].RegisterSetting("int32", k, "", v, "", "", false, false, false)
}

// RegisterInt64 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64(k string, v int) {
	configs[0].RegisterSetting("int64", k, "", v, "", "", false, false, false)
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterString(k, v string) {
	configs[0].RegisterSetting("string", k, "", v, "", "", false, false, false)
}
