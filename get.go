package contour

// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// Config Get Methods

// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// GetE returns the setting Value as an interface{}. If its not a valid
// setting, an error is returned.
func (c *Cfg) GetE(k string) (interface{}, error) {
	idx, err := c.settingIndex(k)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.settings[idx].Value, nil
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

// GetInterfaceE is a convenience wrapper function to Get
func (c *Cfg) GetInterfaceE(k string) (interface{}, error) {
	return c.GetE(k)
}

func (c *Cfg) Get(k string) interface{} {
	s, _ := c.GetE(k)
	return s
}

// GetBool returns the setting Value as a bool.
func (c *Cfg) GetBool(k string) bool {
	s, _ := c.GetBoolE(k)
	return s
}

// GetInt returns the setting Value as an int.
func (c *Cfg) GetInt(k string) int {
	s, _ := c.GetIntE(k)
	return s
}

// GetString returns the setting Value as a string.
func (c *Cfg) GetString(k string) string {
	s, _ := c.GetStringE(k)
	return s
}

// GetInterfac returns the setting Value as an interface
func (c *Cfg) GetInterface(k string) interface{} {
	return c.Get(k)
}

// Convenience functions for configs[app]
// Get returns the setting Value as an interface{}.
// GetE returns the setting Value as an interface{}.
func GetE(k string) (interface{}, error) {
	idx, err := configs[0].settingIndex(k)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return configs[0].settings[idx].Value, nil
}

// GetBoolE returns the setting Value as a bool.
func GetBoolE(k string) (bool, error) {
	v, err := GetE(k)
	if err != nil {
		return false, err
	}

	return *v.(*bool), nil
}

// GetIntE returns the setting Value as an int.
func GetIntE(k string) (int, error) {
	v, err := GetE(k)
	if err != nil {
		return 0, err
	}

	return *v.(*int), nil
}

// GetStringE returns the setting Value as a string.
func GetStringE(k string) (string, error) {
	v, err := GetE(k)
	if err != nil {
		return "", err
	}

	return v.(string), nil
}

// GetInterfaceE is a convenience wrapper function to Get
func GetInterfaceE(k string) (interface{}, error) {
	return configs[0].GetE(k)
}

func Get(k string) interface{} {
	s, _ := configs[0].GetE(k)
	return s
}

// GetBool returns the setting Value as a bool.
func GetBool(k string) bool {
	s, _ := configs[0].GetBoolE(k)
	return s
}

// GetInt returns the setting Value as an int.
func GetInt(k string) int {
	s, _ := configs[0].GetIntE(k)
	return s
}

// GetString returns the setting Value as a string.
func GetString(k string) string {
	s, _ := configs[0].GetStringE(k)
	return s
}

// GetInterface is a convenience wrapper function to Get
func GetInterface(k string) interface{} {
	return configs[0].Get(k)
}
