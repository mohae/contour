package contour

// RegisterFloat64Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterFloat64Core(k string, v float64) {
	c.RegisterSetting("float64", k, "", v, "", "", true, false, false)
}

// RegisterFloat64Core adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterFloat64Core(k string, v  float64) {
	configs[0].RegisterFloat64Core(k, v)
}

// RegisterFloat64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFloat64Conf(k string, v  float64) {
	c.RegisterSetting("float64", k, "", v, "", "", false, true, false)
}

// RegisterFloat64Conf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat64Conf(k string, v  float64) {
	configs[0].RegisterFloat64Conf(k, v)
}

// RegisterFloat64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterFloat64Flag(k, s string, v  float64) {
	c.RegisterSetting("float64", k, s, v, "", "", false, true, true)
}

// RegisterFloat64Flag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat64Flag(k, s string, v  float64) {
	configs[0].RegisterFloat64Flag(k, s, v)
}

// RegisterFloat64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterFloat64(k string, v  float64) {
	c.RegisterSetting("float64", k, "", v, "",  "", false, false, false)
}

// RegisterFloat64 adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterFloat64(k string, v  float64) {
	configs[0].RegisterFloat64(k, v)
}

// UpdateFloat64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateFloat64E(k string, v  float64) error {
	return c.updateE(k, v)
}

// UpdateFloat64E adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateFloat64E(k string, v  float64) error {
	return configs[0].UpdateFloat64E(k, v)
}

// UpdateFloat64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateFloat64(k string, v  float64) {
	c.UpdateFloat64E(k, v)
}

// UpdateFloat64 adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateFloat64(k string, v  float64) {
	configs[0].UpdateFloat64(k, v)
}


// GetFloat64E returns the setting Value as an int.
func (c *Cfg) GetFloat64E(k string) (float64, error) {
	v, err := c.GetE(k)
	if err != nil {
		return 0, err
	}

	switch v.(type) {
	case float64:
		return v.(float64), nil
	case *float64:
		return *v.(*float64), nil
	}

	return 0, nil
}

// GetFloat returns the setting Value as an float64.
func GetFloat64E(k string) (float64, error) {
	return configs[0].GetFloat64E(k)
}

// GetFloat returns the setting Value as an float64.
func (c *Cfg) GetFloat64(k string) float64 {
	s, _ := c.GetFloat64E(k)
	return s
}

// GetFloat returns the setting Value as an float64.
func GetFloat64(k string) float64 {
	return configs[0].GetFloat64(k)
}

