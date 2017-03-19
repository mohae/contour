package contour

// Only non-core settings are updateable. Flags must use Override* to update
// settings.
func (c *Cfg) updateE(k string, v interface{}) error {
	// if can't update, a false will also return an error explaining why.
	_, err := c.canUpdate(k)
	if err != nil {
		return err
	}
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	s, _ := c.settings[k]
	s.Value = v
	c.settings[k] = s
	return nil
}

// UpdateBoolE updates a bool setting, returning any error that occurs.
func UpdateBoolE(k string, v bool) error { return appCfg.UpdateBoolE(k, v) }
func (c *Cfg) UpdateBoolE(k string, v bool) error {
	return c.updateE(k, v)
}

// UpdateBool calls UpdateBoolE and drops the error.
func UpdateBool(k string, v bool) { appCfg.UpdateBool(k, v) }
func (c *Cfg) UpdateBool(k string, v bool) {
	c.UpdateBoolE(k, v)
}

// UpdateIntE updates a int setting, returning any error that occurs.
func UpdateIntE(k string, v int) error { return appCfg.UpdateIntE(k, v) }
func (c *Cfg) UpdateIntE(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateInt calls UpdateIntE and drops the error.
func UpdateInt(k string, v int) { appCfg.UpdateInt(k, v) }
func (c *Cfg) UpdateInt(k string, v int) {
	c.UpdateIntE(k, v)
}

// UpdateInt64E updates a int64 setting, returning any error that occurs.
func UpdateInt64E(k string, v int64) error { return appCfg.UpdateInt64E(k, v) }
func (c *Cfg) UpdateInt64E(k string, v int64) error {
	return c.updateE(k, v)
}

// UpdateInt64 calls UpdateInt64E and drops the error.
func UpdateInt64(k string, v int64) { appCfg.UpdateInt64(k, v) }
func (c *Cfg) UpdateInt64(k string, v int64) {
	c.UpdateInt64E(k, v)
}

// UpdateStringE updates a string setting, returning any error that occurs.
func UpdateStringE(k string, v string) error { return appCfg.UpdateStringE(k, v) }
func (c *Cfg) UpdateStringE(k, v string) error {
	return c.updateE(k, v)
}

// UpdateBool calls UpdateStringE and drops the error.
func UpdateString(k string, v string) { appCfg.UpdateString(k, v) }
func (c *Cfg) UpdateString(k, v string) {
	c.UpdateStringE(k, v)
}

// UpdateCfgFile updates the set config file information.  This only sets
// the filename, the format is not changed.  This does the update
// directly because the cfg filename is a core setting; it will fail the
// canUpdate check.
//
// It is assumed that RegisterCfgFile has already been called, if it hasn't
// nothing will be done.
func UpdateCfgFile(v string) { appCfg.UpdateCfgFile(v) }
func (c *Cfg) UpdateCfgFile(v string) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	s, ok := c.settings[c.confFileKey]
	if !ok {
		return
	}
	s.Value = v
	c.settings[c.confFileKey] = s
}
