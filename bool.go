package contour

// RegisterBoolCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterBoolCore(k string, v bool) {
	c.RegisterSetting("bool", k, "", v, "", "", true, false, false)
}

// RegisterBoolCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterBoolCore(k string, v bool) {
	configs[0].RegisterBoolCore(k, v)
}

// RegisterBoolConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBoolConf(k string, v bool) {
	c.RegisterSetting("bool", k, "", v, "", "", false, true, false)
}

// RegisterBoolConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterBoolConf(k string, v bool) {
	configs[0].RegisterBoolConf(k, v)
}

// RegisterBoolFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBoolFlag(k, s string, v bool) {
	c.RegisterSetting("bool", k, s, v, "", "", false, true, true)
}

// RegisterBoolFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterBoolFlag(k, s string, v bool) {
	configs[0].RegisterBoolFlag(k, s, v)
}

// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterBool(k string, v bool) {
	c.RegisterSetting("bool", k, "", v, "",  "", false, false, false)
}

// RegisterBool adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterBool(k string, v bool) {
	configs[0].RegisterBool(k, v)
}

// Updates
func (c *Cfg) UpdateBoolE(k string, v bool) error {
	return c.updateE(k, v)
}

// UpdateBoolE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBoolE(k string, v bool) error {
	return configs[0].UpdateBoolE(k, v)
}

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateBool(k string, v bool) {
	c.UpdateBoolE(k, v)
}

// UpdateBool adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateBool(k string, v bool) {
	configs[0].UpdateBool(k, v)
}

// GetBoolE returns the setting Value as a bool.
func (c *Cfg) GetBoolE(k string) (bool, error) {
	v, err := c.GetE(k)
	if err != nil {
		return false, err
	}

	switch v.(type) {
	case bool:
		return v.(bool), nil
	case *bool:
		return *v.(*bool), nil
	}

	// Should never happen, but since we know the setting is there and we
	// expect it to be bool, given this method was called, we assume any
	// non-bool/*bool type == false.
	return false, nil
}

// GetBoolE returns the setting Value as a bool.
func GetBoolE(k string) (bool, error) {
	return configs[0].GetBoolE(k)
}

// GetBool returns the setting Value as a bool.
func (c *Cfg) GetBool(k string) bool {
	s, _ := c.GetBoolE(k)
	return s
}

// GetBool returns the setting Value as a bool.
func GetBool(k string) bool {
	return configs[0].GetBool(k)
}
