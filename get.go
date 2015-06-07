package contour

// Get functions and methods.
//
// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// GetE returns the setting Value as an interface{}. If its not a valid
// setting, an error is returned.
func GetE(k string) (interface{}, error) { return appCfg.GetE(k) }
func (c *Cfg) GetE(k string) (interface{}, error) {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	_, ok := c.settings[k]
	if !ok {
		return nil, notFoundErr(k)
	}
	return c.settings[k].Value, nil
}

func Get(k string) interface{} { return appCfg.Get(k) }
func (c *Cfg) Get(k string) interface{} {
	s, _ := c.GetE(k)
	return s
}

// GetBoolE returns the setting Value as a bool.
func GetBoolE(k string) (bool, error) { return appCfg.GetBoolE(k) }
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
	// Should never happen, getting here counts as false
	return false, nil
}

func GetBool(k string) bool { return appCfg.GetBool(k) }
func (c *Cfg) GetBool(k string) bool {
	s, _ := c.GetBoolE(k)
	return s
}

// GetIntE returns the setting Value as an int.
func GetIntE(k string) (int, error) { return appCfg.GetIntE(k) }
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

func GetInt(k string) int { return appCfg.GetInt(k) }
func (c *Cfg) GetInt(k string) int {
	s, _ := c.GetIntE(k)
	return s
}

// GetInt64E returns the setting Value as an int64.
func GetInt64E(k string) (int64, error) { return appCfg.GetInt64E(k) }
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

func GetInt64(k string) int64 { return appCfg.GetInt64(k) }
func (c *Cfg) GetInt64(k string) int64 {
	s, _ := c.GetInt64E(k)
	return s
}

// GetStringE returns the setting Value as a string.
func GetStringE(k string) (string, error) { return appCfg.GetStringE(k) }
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

func GetString(k string) string { return appCfg.GetString(k) }
func (c *Cfg) GetString(k string) string {
	s, _ := c.GetStringE(k)
	return s
}

// GetInterfaceE is a convenience wrapper function to Get
func GetInterfaceE(k string) (interface{}, error) { return appCfg.GetInterfaceE(k) }
func (c *Cfg) GetInterfaceE(k string) (interface{}, error) {
	return c.GetE(k)
}

func GetInterface(k string) interface{} { return appCfg.GetInterface(k) }
func (c *Cfg) GetInterface(k string) interface{} {
	i, _ := c.GetInterfaceE(k)
	return i
}

// Filter Methods obtain a list of flags of the filter type, e.g. boolFilter
// for bool flags, and returns them.
// GetBoolFilterNames returns a list of filter names (flags).
func GetBoolFilterNames() []string { return appCfg.GetBoolFilterNames() }
func (c *Cfg) GetBoolFilterNames() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	var names []string
	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, k)
		}
	}
	return names
}

// GetIntFilterNames returns a list of filter names (flags).
func GetIntFilterNames() []string { return appCfg.GetBoolFilterNames() }
func (c *Cfg) GetIntFilterNames() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	var names []string
	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, k)
		}
	}
	return names
}

// GetInt64FilterNames returns a list of filter names (flags).
func GetInt64FilterNames() []string { return appCfg.GetBoolFilterNames() }
func (c *Cfg) GetInt64FilterNames() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	var names []string
	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "int64" {
			names = append(names, k)
		}
	}
	return names
}

// GetStringFilterNames returns a list of filter names (flags).
func GetStringFilterNames() []string { return appCfg.GetBoolFilterNames() }
func (c *Cfg) GetStringFilterNames() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	var names []string
	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "string" {
			names = append(names, k)
		}
	}
	return names
}
