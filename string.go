package contour

// RegisterStringCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func (c *Cfg) RegisterStringCore(k, v string) {
	c.RegisterSetting("string", k, "", v, "", "", true, false, false)
}

// RegisterStringCore adds the information to the AppsConfig struct, but does not
// save it to its ironment variable
func RegisterStringCore(k, v string) {
	configs[0].RegisterStringCore(k, v)
}

// RegisterStringConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterStringConf(k string, v bool) {
	c.RegisterSetting("string", k, "", v, "", "", false, true, false)
}

// RegisterStringConf adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterStringConf(k string, v bool) {
	configs[0].RegisterStringConf(k, v)
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterStringFlag(k, s, v string) {
	c.RegisterSetting("string", k, s, v, "", "", false, true, true)
}

// RegisterStringFlag adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterStringFlag(k, s, v string) {
	configs[0].RegisterStringFlag(k, s, v)
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func (c *Cfg) RegisterString(k, v string) {
	c.RegisterSetting("string", k, "", v, "",  "", false, false, false)
}

// RegisterString adds the information to the AppsConfig struct, but does not
// save it to its ironment variable.
func RegisterString(k, v string) {
	configs[0].RegisterString(k, v)
}

// UpdateStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateStringE(k, v string) error {
	return c.updateE(k, v)
}

// UpdateStringE adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateStringE(k, v string) error {
	return configs[0].UpdateStringE(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func (c *Cfg) UpdateString(k, v string) {
	c.UpdateStringE(k, v)
}

// UpdateString adds the information to the AppsConfig struct, but does not
// save it to its environment variable.
func UpdateString(k, v string) {
	configs[0].UpdateString(k, v)
}

// GetStringE returns the setting Value as a string.
func (c *Cfg) GetStringE(k string) (string, error) {
	v, err := c.GetE(k)
	if err != nil {
		return "", err
	}

	switch v.(type) {
	case string:
		return v.(string), nil
	case *string:
		return *v.(*string), nil
	}

	return "", nil
}

// GetStringE returns the setting Value as a string.
func GetStringE(k string) (string, error) {
	return configs[0].GetStringE(k)
}

func (c *Cfg) GetString(k string) string {
	s, _ := c.GetStringE(k)
	return s
}

// GetString returns the setting Value as a string.
func GetString(k string) string {
	return configs[0].GetString(k)
}


