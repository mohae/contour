package contour

// Set contains all of contour's Set functions.Calling Set
// adds, or registers, the Settings information to the AppConfig variable.
// The setting value, if there is one, is not saved to its environment
// variable at this point.
//

// SetEnvs writes the current contents of AppConfig to their respective
// environmnet variables.


import (
	"fmt"
	_ "strconv"
)

/*
// SetEnvs goes through AppConfig and saves all of the Settings to their
// environment variables.
func Setenvs() error {
	if !appConfig.UseEnv() {
		return nil
	}

	var err error
	// For each setting
	for k, setting := range appConfig.Settings {
		err = appConfig.Setenv(k, setting)
		if err != nil {
			return err
		}

	}

	return nil

}

// setEnvFromConfigFile goes through all the Settings in the configFile and
// checks to see if the setting is updateable; saving those that are to their
// environment variable.
func setEnvFromConfigFile() error {
	if !appConfig.UseEnv() {
		return nil
	}

	var err error

	for k, v := range configFile {
		// Find the key in the Settings
		_, ok := appConfig.Settings[k]
		if !ok {
			// skip Settings that don't already exist
			continue
		}

		// Skip if Immutable, IsCore, IsEnv since they aren't
		//overridable by ConfigFile.
		if !canUpdate(k) {
			continue
		}

		err = appConfig.Setenv(k, v)
		if err != nil {
			return err
		}

		// Update the setting with file's
		switch appConfig.Settings[k].Type {
		case "string":
			err = UpdateString(k, v.(string))
		case "bool":
			err = UpdateBool(k, v.(bool))
		case "int":
			err = UpdateInt(k, v.(int))
		default:
			return errors.New(k + "'s datatype, " + appConfig.Settings[k].Type + ", is not supported")
		}

		if err != nil {
			return err
		}

	}

	return nil
}
*/

// SetSetting
func (c *Cfg) SetSetting(Type, k string, v interface{}, Code string, IsCore, IsConfig, IsFlag bool) error {
	_, ok := c.Settings[k]
	if ok {
		err := fmt.Errorf("%s: key already exists, cannot add another setting with the same key")
		logger.Error(err)
		return err
	}

	c.Settings[k] = &setting{
		Type:	Type,
		Value:	v,
		Code:	Code,
		IsCore:	IsCore,
		IsConfig:	IsConfig,
		IsFlag:		IsFlag,
	}

	return nil
}

// SetFlagBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagBoolE(k string, v bool, f string) error {
	return c.SetSetting("bool", k, v, f, false, true, true)
}

// SetFlagIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagIntE(k string, v int, f string) error {
	return c.SetSetting("int", k, v, f, false, true, true)
}

// SetFlagStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagStringE(k, v, f string) error {
	return c.SetSetting("string", k, v, f, false, true, true)
}

// SetFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagBool(k string, v bool, f string) {
	c.SetFlagBoolE(k, v, f)
}

// SetFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagInt(k string, v int, f string) {
	c.SetFlagIntE(k, v, f)
}

// SetFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagString(k, v, f string) {
	c.SetFlagStringE(k, v, f)

}

// Convenience functions for configs[app]
// SetFlagBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagBoolE(k string, v bool, f string) error {
	return configs[app].SetSetting("bool", k, v, f, false, true, true)
}

// SetFlagIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagIntE(k string, v int, f string) error {
	return configs[app].SetSetting("int", k, v, f, false, true, true)
}

// SetFlagStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagStringE(k, v, f string) error {
	return configs[app].SetSetting("string", k, v, f, false, true, true)
}

// SetFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagBool(k string, v bool, f string) {
	configs[app].SetFlagBoolE(k, v, f)
}

// SetFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagInt(k string, v int, f string) {
	configs[app].SetFlagIntE(k, v, f)
}

// SetFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagString(k, v, f string) {
	configs[app].SetFlagStringE(k, v, f)
}


