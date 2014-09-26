package contour

// RegisterInt8Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterInt8Core(k string, v int8) {
	c.RegisterSetting("int8", k, "", v, "", "", true, false, false)
}

// RegisterInt8Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterInt8Core(k string, v int8) {
	configs[0].RegisterInt8Core(k, v)
}

// RegisterInt8Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt8Conf(k string, v int8) {
	c.RegisterSetting("int8", k, "", v, "", "", false, true, false)
}

// RegisterInt8Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt8Conf(k string, v int8) {
	configs[0].RegisterInt8Conf(k, v)
}

// RegisterInt8Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt8Flag(k, s string, v int8) {
	c.RegisterSetting("int8", k, s, v, "", "", false, true, true)
}

// RegisterInt8Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt8Flag(k, s string, v int8) {
	configs[0].RegisterInt8Flag(k, s, v)
}

// RegisterInt8 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterInt8(k string, v int8) {
	c.RegisterSetting("int8", k, "", v, "",  "", false, false, false)
}

// RegisterInt8 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt8(k string, v int8) {
	configs[0].RegisterInt8(k, v)
}

// UpdateInt8E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt8E(k string, v int8) error {
	return c.updateE(k, v)
}

// UpdateInt8E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt8E(k string, v int8) error {
	return configs[0].UpdateInt8E(k, v)
}

// UpdateInt8 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt8(k string, v int8) {
	c.UpdateInt8E(k, v)
}

// UpdateInt8 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt8(k string, v int8) {
	configs[0].UpdateInt8(k, v)
}

// GetIntE returns the setting Value as an int.
func (c *Cfg) GetInt8E(k string) (int8, error) {
	v, err := c.GetE(k)
	if err != nil {
		return 0, err
	}

	switch v.(type) {
	case int8:
		return v.(int8), nil
	case *int8:
		return *v.(*int8), nil
	}

	return 0, nil
}

// GetIntE returns the setting Value as an int.
func GetInt8E(k string) (int8, error) {
	return configs[0].GetInt8E(k)
}

// GetInt returns the setting Value as an int.
func (c *Cfg) GetInt8(k string) int8 {
	s, _ := c.GetInt8E(k)
	return s
}


// GetInt returns the setting Value as an int.
func GetInt8(k string) int8 {
	return configs[0].GetInt8(k)
}

