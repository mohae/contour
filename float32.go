package contour

// RegisterFloat32Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterFloat32Core(k string, v float32) {
	c.RegisterSetting("float32", k, "", v, "", "", true, false, false)
}

// RegisterFloat32Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterFloat32Core(k string, v float32) {
	configs[0].RegisterFloat32Core(k, v)
}

// RegisterFloat32Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFloat32Conf(k string, v float32) {
	c.RegisterSetting("float32", k, "", v, "", "", false, true, false)
}

// RegisterFloat32Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat32Conf(k string, v float32) {
	configs[0].RegisterFloat32Conf(k, v)
}

// RegisterFloat32Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFloat32Flag(k, s string, v float32) {
	c.RegisterSetting("float32", k, s, v, "", "", false, true, true)
}

// RegisterFloat32Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat32Flag(k, s string, v float32) {
	configs[0].RegisterFloat32Flag(k, s, v)
}

// RegisterFloat32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterFloat32(k string, v float32) {
	c.RegisterSetting("float32", k, "", v, "",  "", false, false, false)
}

// RegisterFloat32 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat32(k string, v float32) {
	configs[0].RegisterFloat32(k, v)
}

// UpdateFloat32E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateFloat32E(k string, v float32) error {
	return c.updateE(k, v)
}

// UpdateFloat32E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateFloat32E(k string, v float32) error {
	return configs[0].UpdateFloat32E(k, v)
}

// UpdateFloat32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateFloat32(k string, v float32) {
	c.UpdateFloat32E(k, v)
}

// UpdateFloat32 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateFloat32(k string, v float32) {
	configs[0].UpdateFloat32(k, v)
}

// GetIntE returns the setting Value as an int.
func (c *Cfg) GetFloat32E(k string) (float32, error) {
	v, err := c.GetE(k)
	if err != nil {
		return 0, err
	}

	switch v.(type) {
	case float32:
		return v.(float32), nil
	case *float32:
		return *v.(*float32), nil
	}

	return 0, nil
}

// GetFloat returns the setting Value as an float32.
func GetFloat32E(k string) (float32, error) {
	return configs[0].GetFloat32E(k)
}

// GetFloat returns the setting Value as an float32.
func (c *Cfg) GetFloat32(k string) float32 {
	s, _ := c.GetFloat32E(k)
	return s
}


// GetFloat returns the setting Value as an float32.
func GetFloat32(k string) float32 {
	return configs[0].GetFloat32(k)
}

