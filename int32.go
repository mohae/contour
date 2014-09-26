package contour

// RegisterInt32Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterInt32Core(k string, v int32) {
	c.RegisterSetting("int32", k, "", v, "", "", true, false, false)
}

// RegisterInt32Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterInt32Core(k string, v int32) {
	configs[0].RegisterInt32Core(k, v)
}

// RegisterInt32Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt32Conf(k string, v int32) {
	c.RegisterSetting("int32", k, "", v, "", "", false, true, false)
}

// RegisterInt32Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt32Conf(k string, v int32) {
	configs[0].RegisterInt32Conf(k, v)
}

// RegisterInt32Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterInt32Flag(k, s string, v int32) {
	c.RegisterSetting("int32", k, s, v, "", "", false, true, true)
}

// RegisterInt32Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt32Flag(k, s string, v int32) {
	configs[0].RegisterInt32Flag(k, s, v)
}

// RegisterInt32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterInt32(k string, v int32) {
	c.RegisterSetting("int32", k, "", v, "",  "", false, false, false)
}

// RegisterInt32 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt32(k string, v int32) {
	configs[0].RegisterInt32(k, v)
}

// UpdateInt32E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt32E(k string, v int32) error {
	return c.updateE(k, v)
}

// UpdateInt32E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt32E(k string, v int32) error {
	return configs[0].UpdateInt32E(k, v)
}

// UpdateInt32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt32(k string, v int32) {
	c.UpdateInt32E(k, v)
}

// UpdateInt32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt32(k string, v int32) {
	configs[0].UpdateInt32(k, v)
}

// GetIntE returns the setting Value as an int.
func (c *Cfg) GetInt32E(k string) (int32, error) {
	v, err := c.GetE(k)
	if err != nil {
		return 0, err
	}

	switch v.(type) {
	case int32:
		return v.(int32), nil
	case *int32:
		return *v.(*int32), nil
	}

	return 0, nil
}

// GetIntE returns the setting Value as an int.
func GetInt32E(k string) (int32, error) {
	return configs[0].GetInt32E(k)
}

// GetInt returns the setting Value as an int.
func (c *Cfg) GetInt32(k string) int32 {
	s, _ := c.GetInt32E(k)
	return s
}


// GetInt returns the setting Value as an int.
func GetInt32(k string) int32 {
	return configs[0].GetInt32(k)
}

