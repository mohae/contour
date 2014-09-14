package contour

// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// Config Get Methods

// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// GetE returns the setting Value as an interface{}. If its not a valid
// setting, an error is returned.
func (c *config) GetE(k string) (interface{}, error) {
	_, ok := c.Settings[k]
	if !ok {
		return nil, notFoundErr(k)
	}

	return c.Settings[k].Value, nil
}

// GetBoolE returns the setting Value as a bool.
func (c *config) GetBoolE(k string) (bool, error) {
	_, ok := c.Settings[k]
	if !ok {
		return false, notFoundErr(k)
	}

	return c.Settings[k].Value.(bool), nil
}

// GetIntE returns the setting Value as an int.
func (c *config) GetIntE(k string) (int, error) {
	_, ok := c.Settings[k]
	if !ok {
		return 0, notFoundErr(k)
	}

	return c.Settings[k].Value.(int), nil
}

// GetStringE returns the setting Value as a string.
func (c *config) GetStringE(k string) (string, error) {
	_, ok := c.Settings[k]
	if !ok {
		return "", notFoundErr(k)
	}

	return c.Settings[k].Value.(string), nil
}

// GetInterfaceE is a convenience wrapper function to Get
func (c *config) GetInterfaceE(k string) (interface{}, error) {
	return c.GetE(k)
}

func (c *config) Get(k string) interface{} {
	s, _ := c.GetE(k)
	return s
}

// GetBool returns the setting Value as a bool.
func (c *config) GetBool(k string) bool {
	s, _ := c.GetBoolE(k)
	return s
}

// GetInt returns the setting Value as an int.
func (c *config) GetInt(k string) int {
	s, _ := c.GetIntE(k)
	return s
}

// GetString returns the setting Value as a string.
func (c *config) GetString(k string) string {
	s, _ := c.GetStringE(k)
	return s
}

// GetInterfac returns the setting Value as an interface
func (c *config) GetInterface(k string) interface{} {
	return c.Get(k)
}


// Filter Methods obtain a list of flags of the filter type, e.g. boolFilter 
// for bool flags, and returns them.
// GetBoolFilterNames returns a list of filter names (flags).
func (c *config) GetBoolFilterNames() []string {
	var names []string

	for k, setting := range c.Settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, k)
		}
	}

	return names
}

// GetIntFilterNames returns a list of filter names (flags).
func (c *config) GetIntFilterNames() []string {
	var names []string

	for k, setting := range c.Settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, k)
		}
	}

	return names
}

// GetStringFilterNames returns a list of filter names (flags).
func (c *config) GetStringFilterNames() []string {
	var names []string

	for k, setting := range c.Settings {
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
	_, ok := configs[app].Settings[k]
	if !ok {
		return nil, notFoundErr(k)
	}

	return configs[app].Settings[k].Value, nil
}

// GetBoolE returns the setting Value as a bool.
func GetBoolE(k string) (bool, error) {
	_, ok := configs[app].Settings[k]
	if !ok {
		return false, notFoundErr(k)
	}

	return *configs[app].Settings[k].Value.(*bool), nil
}

// GetIntE returns the setting Value as an int.
func GetIntE(k string) (int, error) {
	_, ok := configs[app].Settings[k]
	if !ok {
		return 0, notFoundErr(k)
	}

	return *configs[app].Settings[k].Value.(*int), nil
}

// GetStringE returns the setting Value as a string.
func GetStringE(k string) (string, error) {
	_, ok := configs[app].Settings[k]
	if !ok {
		return "", notFoundErr(k)
	}

	return configs[app].Settings[k].Value.(string), nil
}

// GetInterfaceE is a convenience wrapper function to Get
func GetInterfaceE(k string) (interface{}, error) {
	return configs[app].GetE(k)
}

func Get(k string) interface{} {
	s, _ := configs[app].GetE(k)
	return s
}


// GetBool returns the setting Value as a bool.
func GetBool(k string) bool {
	s, _ := configs[app].GetBoolE(k)
	return s
}


// GetInt returns the setting Value as an int.
func GetInt(k string) int {
	s, _ := configs[app].GetIntE(k)
	return s
}

// GetString returns the setting Value as a string.
func GetString(k string) string {
	s, _ := configs[app].GetStringE(k)
	return s
}

// GetInterface is a convenience wrapper function to Get
func GetInterface(k string) interface{} {	
	return configs[app].Get(k)
}

// GetBoolFilterNames returns a list of filter names (flags).
func GetBoolFilterNames() []string {
	var names []string

	for k, setting := range configs[app].Settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, k)
		}
	}

	return names
}

// GetIntFilterNames returns a list of filter names (flags).
func GetIntFilterNames() []string {
	var names []string

	for k, setting := range configs[app].Settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, k)
		}
	}

	return names
}

// GetStringFilterNames returns a list of filter names (flags).
func GetStringFilterNames() []string {
	var names []string

	for k, setting := range configs[app].Settings {
		if setting.IsFlag && setting.Type == "string" {
			names = append(names, k)
		}
	}

	return names
}




