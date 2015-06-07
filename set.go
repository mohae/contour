package contour

import (
	"fmt"
	// "strconv"
)

// Set contains all of contour's Set functions.  Calling Set adds, or
// registers, the settings information to the AppConfig variable.

// setCfg set's the configuration information from the received map.
func (c *Cfg) setCfg(cf map[string]interface{}) error {
	c.RWMutex.Lock()
	if !c.useCfgFile {
		c.RWMutex.Unlock()
		return nil
	}
	c.RWMutex.Unlock()
	for k, v := range cf {
		c.RWMutex.RLock()
		// Find the key in the settings
		_, ok := c.settings[k]
		c.RWMutex.RUnlock()
		if !ok {
			// skip settings that don't already exist
			continue
		}
		err := c.updateE(k, v)
		if err != nil {
			return err
		}

	}
	return nil
}

// SetSetting
func (c *Cfg) SetSetting(typ, name, short string, v interface{}, dflt, usage string, IsCore, IsCfg, IsFlag bool) error {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	_, ok := c.settings[name]
	if ok {
		err := fmt.Errorf("%s: key already exists, cannot add another setting with the same key")
		return err
	}
	c.settings[name] = setting{
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

// SetFlagBoolE adds the information to the appCfg global.
func SetFlagBoolE(k, s string, v bool, dflt, u string) error {
	return appCfg.SetFlagBoolE(k, s, v, dflt, u)
}
func (c *Cfg) SetFlagBoolE(k, s string, v bool, dflt, u string) error {
	return c.SetSetting("bool", k, s, v, dflt, u, false, true, true)
}

// SetFlagBool adds the information to the appCfg global.
func SetFlagBool(k, s string, v bool, dflt, u string) {
	appCfg.SetFlagBool(k, s, v, dflt, u)
}
func (c *Cfg) SetFlagBool(k, s string, v bool, dflt, u string) {
	c.SetFlagBoolE(k, s, v, dflt, u)
}

// SetFlagIntE adds the information to the appCfg global.
func SetFlagIntE(k, s string, v int, dflt, u string) error {
	return appCfg.SetFlagIntE(k, s, v, dflt, u)
}
func (c *Cfg) SetFlagIntE(k, s string, v int, dflt, u string) error {
	return c.SetSetting("int", k, s, v, dflt, u, false, true, true)
}

// SetFlagInt adds the information to the appCfg global.
func SetFlagInt(k, s string, v int, dflt, u string) {
	appCfg.SetFlagInt(k, s, v, dflt, u)
}
func (c *Cfg) SetFlagInt(k, s string, v int, dflt, u string) {
	c.SetFlagIntE(k, s, v, dflt, u)
}

// SetFlagInt64E adds the information to the appCfg global.
func SetFlagInt64E(k, s string, v int64, dflt, u string) error {
	return appCfg.SetFlagInt64E(k, s, v, dflt, u)
}
func (c *Cfg) SetFlagInt64E(k, s string, v int64, dflt, u string) error {
	return c.SetSetting("int64", k, s, v, dflt, u, false, true, true)
}

// SetFlagInt64 adds the information to the appCfg global.
func SetFlagInt64(k, s string, v int64, dflt, u string) {
	appCfg.SetFlagInt64(k, s, v, dflt, u)
}
func (c *Cfg) SetFlagInt64(k, s string, v int64, dflt, u string) {
	c.SetFlagInt64E(k, s, v, dflt, u)
}

// SetFlagStringE adds the information to the appCfg global.
func SetFlagStringE(k, s, v, dflt, u string) error {
	return appCfg.SetFlagStringE(k, s, v, dflt, u)
}
func (c *Cfg) SetFlagStringE(k, s, v, dflt, u string) error {
	return c.SetSetting("string", k, s, v, dflt, u, false, true, true)
}

// SetFlagString adds the information to the appCfg global.
func SetFlagString(k, s, v, dflt, u string) {
	appCfg.SetFlagString(k, s, v, dflt, u)
}
func (c *Cfg) SetFlagString(k, s, v, dflt, u string) {
	c.SetFlagStringE(k, s, v, dflt, u)

}

// SetName set's the cfg's name.
func SetName(name string) { appCfg.SetName(name) }
func (c *Cfg) SetName(name string) {
	c.name = name
}

// Name returns the cfg's name.
func Name() string { return appCfg.Name() }
func (c *Cfg) Name() string {
	return c.name
}
