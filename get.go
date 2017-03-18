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

// BoolE returns the key's value as a bool. A SettingNotFoundErr is returned
// if the key is not valid. If the setting's type is not a bool, a DataTypeErr
// will be returned.
func BoolE(k string) (bool, error) { return appCfg.BoolE(k) }
func (c *Cfg) BoolE(k string) (bool, error) {
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

// Bool returns the key's value as a bool. A false will be returned if the key
// either doesn't exist or is not a bool setting.
func Bool(k string) bool { return appCfg.Bool(k) }
func (c *Cfg) Bool(k string) bool {
	s, _ := c.BoolE(k)
	return s
}

// IntE returns the key's value as an int. A SettingNotFoundErr is returned if
// the key is not valid. If the setting's type is not an int, a DataTypeErr
// will be returned.
func IntE(k string) (int, error) { return appCfg.IntE(k) }
func (c *Cfg) IntE(k string) (int, error) {
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

// Int returns the key's value as an int. A 0 will be returned if the key
// either doesn't exist or is not an int setting.
func Int(k string) int { return appCfg.Int(k) }
func (c *Cfg) Int(k string) int {
	s, _ := c.IntE(k)
	return s
}

// Int64E returns the key's value as an int64. A SettingNotFoundErr is returned
// if the key is not valid. If the setting's type is neither an int64 nor an
// int, a DataTypeErr will be returned.
func Int64E(k string) (int64, error) { return appCfg.Int64E(k) }
func (c *Cfg) Int64E(k string) (int64, error) {
	v, err := c.GetE(k)
	if err != nil {
		return 0, err
	}
	switch v.(type) {
	case int64:
		return v.(int64), nil
	case *int64:
		return *v.(*int64), nil
	case int:
		return int64(v.(int)), nil
	case *int:
		return int64(*v.(*int)), nil
	}

	// Is neither an int64 nor an int.
	return 0, DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _int64}
}

// Int64 returns the key's value as an int64. A 0 will be returned if the key
// either doesn't exist or is neither an int64 nor an int setting.
func Int64(k string) int64 { return appCfg.Int64(k) }
func (c *Cfg) Int64(k string) int64 {
	s, _ := c.Int64E(k)
	return s
}

// StringE returns the key's value as a string. A SettingNotFoundErr is
// returned if the key is not gvalid. If the setting's type is not a string, a
// DataTypeErr will be returned.
func StringE(k string) (string, error) { return appCfg.StringE(k) }
func (c *Cfg) StringE(k string) (string, error) {
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

// String returns the key's value as a string. An empty string, "", will be
// returned if the key either doesn't exist or is not a string setting.
func String(k string) string { return appCfg.String(k) }
func (c *Cfg) String(k string) string {
	s, _ := c.StringE(k)
	return s
}
