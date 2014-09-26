package contour

// RegisterIntCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterIntCore(k string, v int) {
	c.RegisterSetting("int", k, "", v, "", "", true, false, false)
}

// RegisterIntCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterIntCore(k string, v int) {
	configs[0].RegisterIntCore(k, v)
}

// RegisterIntConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterIntConf(k string, v bool) {
	c.RegisterSetting("int", k, "", v, "", "", false, true, false)
}

// RegisterIntConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterIntConf(k string, v bool) {
	configs[0].RegisterIntConf(k, v)
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterIntFlag(k, s string, v int) {
	c.RegisterSetting("int", k, s, v, "", "", false, true, true)
}

// RegisterIntFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterIntFlag(k, s string, v int) {
	configs[0].RegisterIntFlag(k, s, v)
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) RegisterInt(k string, v int) {
	c.RegisterSetting("int", k, "", v, "",  "", false, false, false)
}

// RegisterInt adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterInt(k string, v int) {
	configs[0].RegisterInt(k, v)
}

// UpdateIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateIntE(k string, v int) error {
	return c.updateE(k, v)
}

// UpdateIntE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateIntE(k string, v int) error {
	return configs[0].UpdateIntE(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateInt(k string, v int) {
	c.UpdateIntE(k, v)
}

// UpdateInt adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateInt(k string, v int) {
	configs[0].UpdateInt(k, v)
}

// GetIntE returns the setting Value as an int.
func (c *Cfg) GetIntE(k string) (int, error) {
	v, err := c.GetE(k)
	if err != nil {
		return 0, err
	}

	switch v.(type) {
	case int:
		return v.(int), nil
	case *int:
		return *v.(*int), nil
	}

	return 0, nil
}

// GetIntE returns the setting Value as an int.
func GetIntE(k string) (int, error) {
	return configs[0].GetIntE(k)
}

// GetInt returns the setting Value as an int.
func (c *Cfg) GetInt(k string) int {
	s, _ := c.GetIntE(k)
	return s
}


// GetInt returns the setting Value as an int.
func GetInt(k string) int {
	return configs[0].GetInt(k)
}

