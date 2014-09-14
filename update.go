package contour

import "fmt"

// Only non-core settings are updateable.
// Flags must use Override* to update settings.
// save it to its environment variable.
func (c *Cfg) updateE(k string, v interface{}) error {
	if !c.canUpdate(k) {
		err := fmt.Errorf("config[%s]: %s is not updateable", c.name, k)
		logger.Warn(err)
		return err
	}

	c.Settings[k].Value = v
	return nil
}

func (c *Cfg) UpdateBoolE(k string, v bool) error {
	return c.updateE(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateIntE(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateStringE(k, v string) error {
	return c.updateE(k, v)
}

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateBool(k string, v bool) {
	c.UpdateBoolE(k, v)
}


// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt(k string, v int) {
	c.UpdateIntE(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateString(k, v string) {
	c.UpdateStringE(k, v)
}

// UpdateBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBoolE(k string, v bool) error {
	return configs[app].updateE(k, v)
}

// UpdateIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateIntE(k string, v int) error {
	return configs[app].updateE(k, v)
}

// UpdateStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateStringE(k, v string) error {
	return configs[app].updateE(k, v)
}

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBool(k string, v bool) {
	configs[app].UpdateBool(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt(k string, v int) {
	configs[app].UpdateInt(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateString(k, v string) {
	configs[app].UpdateString(k, v)
}
