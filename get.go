package contour

// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// Config Get Methods

// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// GetE returns the setting Value as an interface{}. If its not a valid
// setting, an error is returned.
func (c *Cfg) GetE(k string) (interface{}, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.settings[k]
	if !ok {
		return nil, notFoundErr(k)
	}

	return c.settings[k].Value, nil
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

	// Should never happen, getting here counts as false
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

// GetInt64E returns the setting Value as an int.
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

// GetInt64 returns the setting Value as an int.
func (c *Cfg) GetInt64(k string) int64 {
	s, _ := c.GetInt64E(k)
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

// Filter Methods obtain a list of flags of the filter type, e.g. boolFilter
// for bool flags, and returns them.
// GetBoolFilterNames returns a list of filter names (flags).
func (c *Cfg) GetBoolFilterNames() []string {
	var names []string

	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, k)
		}
	}

	return names
}

// GetIntFilterNames returns a list of filter names (flags).
func (c *Cfg) GetIntFilterNames() []string {
	var names []string

	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, k)
		}
	}

	return names
}

// GetInt64FilterNames returns a list of filter names (flags).
func (c *Cfg) GetInt64FilterNames() []string {
	var names []string

	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "int64" {
			names = append(names, k)
		}
	}

	return names
}

// GetStringFilterNames returns a list of filter names (flags).
func (c *Cfg) GetStringFilterNames() []string {
	var names []string

	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "string" {
			names = append(names, k)
		}
	}

	return names
}

// Convenience functions for configs[app]
// Get returns the setting Value as an interface{}.
// GetE returns the setting Value as an interface{}.
func GetE(k string) (interface{}, error) {
	return appCfg.GetE(k)
}

// GetBoolE returns the setting Value as a bool.
func GetBoolE(k string) (bool, error) {
	return appCfg.GetBoolE(k)
}

// GetIntE returns the setting Value as an int.
func GetIntE(k string) (int, error) {
	return appCfg.GetIntE(k)
}

// GetInt64E returns the setting Value as an int.
func GetInt64E(k string) (int64, error) {
	return appCfg.GetInt64E(k)
}

// GetStringE returns the setting Value as a string.
func GetStringE(k string) (string, error) {
	return appCfg.GetStringE(k)
}

// GetInterfaceE is a convenience wrapper function to Get
func GetInterfaceE(k string) (interface{}, error) {
	return appCfg.GetE(k)
}

func Get(k string) interface{} {
	s, _ := appCfg.GetE(k)
	return s
}

// GetBool returns the setting Value as a bool.
func GetBool(k string) bool {
	b, _ := appCfg.GetBoolE(k)
	return b
}

// GetInt returns the setting Value as an int.
func GetInt(k string) int {
	s, _ := appCfg.GetIntE(k)
	return s
}

// GetInt64 returns the setting Value as an int.
func GetInt64(k string) int64 {
	s, _ := appCfg.GetInt64E(k)
	return s
}

// GetString returns the setting Value as a string.
func GetString(k string) string {
	s, _ := appCfg.GetStringE(k)
	return s
}

// GetInterface is a convenience wrapper function to Get
func GetInterface(k string) interface{} {
	return appCfg.Get(k)
}

// GetBoolFilterNames returns a list of filter names (flags).
func GetBoolFilterNames() []string {
	var names []string

	for k, setting := range appCfg.settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, k)
		}
	}

	return names
}

// GetIntFilterNames returns a list of filter names (flags).
func GetIntFilterNames() []string {
	var names []string

	for k, setting := range appCfg.settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, k)
		}
	}

	return names
}

// GetInt64FilterNames returns a list of filter names (flags).
func GetInt64FilterNames() []string {
	var names []string

	for k, setting := range appCfg.settings {
		if setting.IsFlag && setting.Type == "int64" {
			names = append(names, k)
		}
	}

	return names
}

// GetStringFilterNames returns a list of filter names (flags).
func GetStringFilterNames() []string {
	var names []string

	for k, setting := range appCfg.settings {
		if setting.IsFlag && setting.Type == "string" {
			names = append(names, k)
		}
	}

	return names
}
