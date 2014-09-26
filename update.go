package contour

import "fmt"

// Only non-core settings are updateable.
// Flags must use Override* to update settings.
// save it to its environment variable.
func (c *Cfg) idxUpdateE(i int, v interface{}) error {
	if !c.canUpdate(i) {
		err := fmt.Errorf("config[%s]: %s is not updateable", c.name, c.setting[i].Name)
		logger.Warn(err)
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.settings[i].Value = v
	return nil
}

func (c *Cfg) updateE(k string, v interface{}) error {
	// setting must exist to update
	idx, err := c.settingIndex(k)
	if err != nil {
		logger.Error(err)
		return err
	}

	return c.idxUpdateE(idx, v)
}


func (c *Cfg) UpdateBoolE(k string, v bool) error {
	return c.updateE(k, v)
}

// UpdateFloat32E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateFloat32E(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateFloat64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateFloat64E(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateIntE(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateInt8E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt8E(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateInt32E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt32E(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateInt64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt64E(k string, v int) error {
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

// UpdateFloat32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateFloat32(k string, v int) {
	c.UpdateFloat32E(k, v)
}

// UpdateFloat64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateFloat64(k string, v int) {
	c.UpdateFloat64E(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt(k string, v int) {
	c.UpdateIntE(k, v)
}

// UpdateInt8 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt8(k string, v int) {
	c.UpdateInt8E(k, v)
}

// UpdateInt32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt32(k string, v int) {
	c.UpdateInt32E(k, v)
}

// UpdateInt64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt64(k string, v int) {
	c.UpdateInt64E(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateString(k, v string) {
	c.UpdateStringE(k, v)
}

// UpdateBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBoolE(k string, v bool) error {
	return configs[0].updateE(k, v)
}

// UpdateFloat32E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateFloat32E(k string, v int) error {
	return configs[0].updateE(k, v)
}

// UpdateFloat64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateFloat64E(k string, v int) error {
	return configs[0].updateE(k, v)
}

// UpdateIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateIntE(k string, v int) error {
	return configs[0].updateE(k, v)
}

// UpdateInt8E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt8E(k string, v int) error {
	return configs[0].updateE(k, v)
}

// UpdateInt32E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt32E(k string, v int) error {
	return configs[0].updateE(k, v)
}

// UpdateIn64tE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt64E(k string, v int) error {
	return configs[0].updateE(k, v)
}

// UpdateStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateStringE(k, v string) error {
	return configs[0].updateE(k, v)
}

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBool(k string, v bool) {
	configs[0].UpdateBool(k, v)
}

// UpdateFloat32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateFloat32(k string, v int) {
	configs[0].UpdateInt(k, v)
}

// UpdateFloat64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateFloat64(k string, v int) {
	configs[0].UpdateInt(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt(k string, v int) {
	configs[0].UpdateInt(k, v)
}

// UpdateInt8 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt8(k string, v int) {
	configs[0].UpdateInt(k, v)
}

// UpdateInt32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt32(k string, v int) {
	configs[0].UpdateInt(k, v)
}

// UpdateInt64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt64(k string, v int) {
	configs[0].UpdateInt(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateString(k, v string) {
	configs[0].UpdateString(k, v)
}
