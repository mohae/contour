package contour

import (
	"fmt"
	"strconv"
)

// Only non-core settings are updateable.
// Flags must use Override* to update settings.
// save it to its environment variable.
func (c *Cfg) updateE(k string, v interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.canUpdate(k) {
		err := fmt.Errorf("config[%s]: %q is not updateable", c.name, k)
		logger.Warn(err)
		return err
	}

	c.settings[k].Value = v
	return nil
}

func (c *Cfg) UpdateBoolE(k, v string) error {
	if v != "" {
		_, err := strconv.ParseBool(v)
		if err != nil {
			v = ""
		}
	}
	return c.updateE(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateIntE(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateInt64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt64E(k string, v int64) error {
	return c.updateE(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateStringE(k, v string) error {
	return c.updateE(k, v)
}

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateBool(k, v string) {
	c.UpdateBoolE(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt(k string, v int) {
	c.UpdateIntE(k, v)
}

// UpdateInt64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt64(k string, v int64) {
	c.UpdateInt64E(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateString(k, v string) {
	c.UpdateStringE(k, v)
}

// UpdateBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBoolE(k, v string) error {
	return appCfg.updateE(k, v)
}

// UpdateIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateIntE(k string, v int) error {
	return appCfg.updateE(k, v)
}

// UpdateInt64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt64E(k string, v int64) error {
	return appCfg.updateE(k, v)
}

// UpdateStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateStringE(k, v string) error {
	return appCfg.updateE(k, v)
}

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBool(k, v string) {
	appCfg.UpdateBool(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt(k string, v int) {
	appCfg.UpdateInt(k, v)
}

// UpdateInt64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt64(k string, v int64) {
	appCfg.UpdateInt64(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateString(k, v string) {
	appCfg.UpdateString(k, v)
}
