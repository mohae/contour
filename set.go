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
		c.lock.RLock()
		// Find the key in the settings
		_, ok := c.settings[k]
		c.lock.RUnlock()
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
func (c *Cfg) SetSetting(typ, name, short string, v, dflt interface{}, usage string, IsCore, IsCfg, IsFlag bool) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.settings[name]
	if ok {
		err := fmt.Errorf("%s: key already exists, cannot add another setting with the same key")
		logger.Error(err)
		return err
	}

	c.settings[name] = &setting{
		Type:    typ,
		Name:    name,
		Short:   short,
		Value:   v,
		Default: dflt,
		Usage:   usage,
		IsCore:  IsCore,
		IsCfg:   IsCfg,
		IsFlag:  IsFlag,
	}

	return nil
}

// SetFlagBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagBoolE(k, s string, v, dflt bool, u string) error {
	return c.SetSetting("bool", k, s, v, dflt, u, false, true, true)
}

// SetFlagIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagIntE(k, s string, v, dflt int, u string) error {
	return c.SetSetting("int", k, s, v, dflt, u, false, true, true)
}

// SetFlagStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagStringE(k, s, v, dflt, u string) error {
	return c.SetSetting("string", k, s, v, dflt, u, false, true, true)
}

// SetFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagBool(k, s string, v, dflt bool, u string) {
	c.SetFlagBoolE(k, s, v, dflt, u)
}

// SetFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagInt(k, s string, v, dflt int, u string) {
	c.SetFlagIntE(k, s, v, dflt, u)
}

// SetFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) SetFlagString(k, s, v, dflt, u string) {
	c.SetFlagStringE(k, s, v, dflt, u)

}

// Convenience functions for configs[app]
// SetFlagBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagBoolE(k, s string, v, dflt bool, u string) error {
	return appCfg.SetFlagBoolE(k, s, v, dflt, u)
}

// SetFlagIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagIntE(k, s string, v, dflt int, u string) error {
	return appCfg.SetFlagIntE(k, s, v, dflt, u)
}

// SetFlagStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagStringE(k, s, v, dflt, u string) error {
	return appCfg.SetFlagStringE(k, s, v, dflt, u)
}

// SetFlagBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagBool(k, s string, v, dflt bool, u string) {
	appCfg.SetFlagBoolE(k, s, v, dflt, u)
}

// SetFlagInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagInt(k, s string, v, dflt int, u string) {
	appCfg.SetFlagIntE(k, s, v, dflt, u)
}

// SetFlagString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func SetFlagString(k, s, v, dflt, u string) {
	appCfg.SetFlagStringE(k, s, v, dflt, u)
}
