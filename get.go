package contour

import "reflect"

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
		return nil, SettingNotFoundErr{name: k}
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
	// Isn't a bool.
	return false, DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _bool}
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

	// Isn't an int.
	return 0, DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _int}
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

	// Isn't an int64.
	return 0, DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _int64}
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

	// Isn't a string.
	return "", DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _string}
}

func GetString(k string) string { return appCfg.GetString(k) }
func (c *Cfg) GetString(k string) string {
	s, _ := c.GetStringE(k)
	return s
}
