package contour

// Set contains all of contour's Set functions.Calling Set
// adds, or registers, the settings information to the AppConfig variable.
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
// SetEnvs goes through AppConfig and saves all of the settings to their
// environment variables.
func Setenvs() error {
	if !appConfig.UseEnv() {
		return nil
	}

	var err error
	// For each setting
	for k, setting := range appConfig.settings {
		err = appConfig.Setenv(k, setting)
		if err != nil {
			return err
		}

	}

	return nil

}
*/

// setEnvFromConfigFile goes through all the settings in the configFile and
// checks to see if the setting is updateable; saving those that are to their
// environment variable.
func (c *Cfg) setCfg(cf map[string]interface{}) error {
	if !c.UseEnv() {
		return nil
	}

	for k, v := range cf {
		c.Lock.RLock()
		// Find the key in the settings
		_, ok := c.settings[k]
		c.Lock.RUnlock()
		if !ok {
			// skip settings that don't already exist
			continue
		}

//		err = appConfig.Setenv(k, v)
//		if err != nil {
//			return err
//		}

		err := c.updateE(k, v)
		if err != nil {
			return err
		}

	}

	return nil
}


// SetSetting
func (c *Cfg) SetSetting(Type, k string, v interface{}, Code string, IsCore, IsCfg, IsFlag bool) error {
	c.Lock.Lock()
	defer c.Lock.Unlock()

	_, ok := c.settings[k]
	if ok {
		err := fmt.Errorf("%s: key already exists, cannot add another setting with the same key")
		logger.Error(err)
		return err
	}
	
	c.settings[k] = &setting{
		Type:	Type,
		Value:	v,
		Code:	Code,
		IsCore:	IsCore,
		IsCfg:	IsCfg,
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


