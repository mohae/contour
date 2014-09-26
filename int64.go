package contour

// RegisterInt64Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterInt64Core(k string, v int64) {
	c.RegisterSetting("int64", k, "", v, "", "", true, false, false)
}

// RegisterInt64Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterInt64Core(k string, v int64) {
	configs[0].RegisterInt64Core(k, v)
}

// RegisterInt64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt64Conf(k string, v int64) {
	c.RegisterSetting("int64", k, "", v, "", "", false, true, false)
}

// RegisterInt64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64Conf(k string, v int64) {
	configs[0].RegisterInt64Conf(k, v)
}

// RegisterInt64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt64Flag(k, s string, v int64) {
	c.RegisterSetting("int64", k, s, v, "", "", false, true, true)
}

// RegisterInt64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64Flag(k, s string, v int64) {
	configs[0].RegisterInt64Flag(k, s, v)
}

// RegisterInt64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterInt64(k string, v int64) {
	c.RegisterSetting("int64", k, "", v, "",  "", false, false, false)
}

// RegisterInt64 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt64(k string, v int64) {
	configs[0].RegisterInt64(k, v)
}

// UpdateInt64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt64E(k string, v int64) error {
	return c.updateE(k, v)
}

// UpdateInt64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt64E(k string, v int64) error {
	return configs[0].UpdateInt64E(k, v)
}

// UpdateInt64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt64(k string, v int64) {
	c.UpdateInt64E(k, v)
}

// UpdateInt64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt64(k string, v int64) {
	configs[0].UpdateInt64(k, v)
}

// GetIntE returns the setting Value as an int.
func (c *Cfg) GetInt64E(k string) (int64, error) {
	v, err := c.GetE(k)
	if err != nil {
		return 0, err
	}

	switch v.(type) {
	case int64:
		return v.(int64), nil
	case *int64:
		return *v.(*int64), nil
	}

	return 0, nil
}

// GetIntE returns the setting Value as an int.
func GetInt64E(k string) (int64, error) {
	return configs[0].GetInt64E(k)
}

// GetInt returns the setting Value as an int.
func (c *Cfg) GetInt64(k string) int64 {
	s, _ := c.GetInt64E(k)
	return s
}


// GetInt returns the setting Value as an int.
func GetInt64(k string) int64 {
	return configs[0].GetInt64(k)
}

