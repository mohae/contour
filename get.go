package contour

import "reflect"

// Get functions and methods.
//
// E versions return an error if one occurs. Non-E versions return the zero
// value if an error occurs.

// GetE returns settings' value for k as an interface{}. A SettingNotFoundError
// is returned if k doesn't exist.
func (s *Settings) GetE(k string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.get(k)
}

// This assumes the lock has already been obtained. Since this is not exported,
// it does not need to end with E to signify it returns an error as it can be
// assumed it does.
func (s *Settings) get(k string) (interface{}, error) {
	_, ok := s.settings[k]
	if !ok {
		return nil, SettingNotFoundError{k: k}
	}
	return s.settings[k].Value, nil
}

// Get returns the settings' value for k as an interface{}. A nil is returned
// if k doesn't exist.
func (s *Settings) Get(k string) interface{} {
	v, _ := s.GetE(k)
	return v
}

// BoolE returns the settings' value for k as a bool. A SettingNotFoundError is
// returned if k doesn't exist. A DataTypeError will be returned if the value
// is not a bool.
func (s *Settings) BoolE(k string) (bool, error) {
	v, err := s.GetE(k)
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
	return false, DataTypeError{k: k, is: reflect.TypeOf(v).String(), not: _bool}
}

// Bool returns the settings' value for k as a bool. A false will be returned
// if k either doesn't exist or if its value is not a bool.
func (s *Settings) Bool(k string) bool {
	v, _ := s.BoolE(k)
	return v
}

// IntE returns the settings' value for k as an int. A SettingNotFoundError is
// returned if k doesn't exist. A DataTypeError will be returned if the value
// is not an int.
func (s *Settings) IntE(k string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.int(k)
}

// This assumes the lock has already been obtained. Unexported methods don't
// need to be suffixed with E to show they return an error.
func (s *Settings) int(k string) (int, error) {
	v, err := s.get(k)
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
	return 0, DataTypeError{k: k, is: reflect.TypeOf(v).String(), not: _int}
}

// Int returns the settings' value for k as an int. A 0 will be returned if k
// either doesn't exist or if its value is nat an int.
func (s *Settings) Int(k string) int {
	v, _ := s.IntE(k)
	return v
}

// Int64E returns the settings value for k as an int64. A SettingNotFoundError
// is returned if k doesn't exist. A DataTypeError will be returned if the
// value is neither an int64 nor an int.
func (s *Settings) Int64E(k string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.int64(k)
}

// It is assumed that the caller has the lock. Unexported methods don't need
// to be suffixed with an E to signify they return an error.
func (s *Settings) int64(k string) (int64, error) {
	v, err := s.get(k)
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
	return 0, DataTypeError{k: k, is: reflect.TypeOf(v).String(), not: _int64}
}

// Int64 returns the settings value for k as an int64. A 0 will be returned if
// k doesn't exist or if its value is neither an int64 nor an int.
func (s *Settings) Int64(k string) int64 {
	v, _ := s.Int64E(k)
	return v
}

// InterfaceE returns the settings' value for k as an interface{}. A
// SettingNotFoundError is returned if k doesn't exist.
func (s *Settings) InterfaceE(k string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.get(k)
}

// Interface returns the settings' value for k as an interface{}. A nil will be
// returned if k doesn't exist.
func (s *Settings) Interface(k string) interface{} {
	v, _ := s.InterfaceE(k)
	return v
}

// StringE returns the settings value for k as a string. A SettingNotFoundError
// is returned if k doesn't exist. A DataTypeError will be returned if the
// value is not a string.
func (s *Settings) StringE(k string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.string(k)
}

// It is assumed that the caller holds the lock. Unexported methods don't need
// to be suffixed with E to signify they return an error.
func (s *Settings) string(k string) (string, error) {
	v, err := s.get(k)
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
	return "", DataTypeError{k: k, is: reflect.TypeOf(v).String(), not: _string}
}

// String returns the settings value for k as a string. An empty string, "",
// will be returned if k either doesn't exist or if its value is not a string.
func (s *Settings) String(k string) string {
	v, _ := s.StringE(k)
	return v
}

// GetE returns the standard settings' value for k as an interface{}. A
// SettingNotFoundError is returned if k doesn't exist.
func GetE(k string) (interface{}, error) { return std.GetE(k) }

// Get returns the standard settings' value for k as an interface{}. A nil is
// returned if k doesn't exist.
func Get(k string) interface{} { return std.Get(k) }

// BoolE returns the standard settings' value for k as a bool. A
//vSettingNotFoundError is returned if k doesn't exist. A DataTypeError will be
// returned if the value is not bool.
func BoolE(k string) (bool, error) { return std.BoolE(k) }

// Bool returns the standard settings' value for k as a bool. A false will be
// returned if k doesn't exist or if its value is not a bool.
func Bool(k string) bool { return std.Bool(k) }

// IntE returns the standard settings' value for k as an int. A
// SettingNotFoundError is returned if k doesn't exist. A DataTypeError will be
// returned if the value is not an int.
func IntE(k string) (int, error) { return std.IntE(k) }

// Int returns the standard settings' value for k as an int. A 0 will be
// returned if k doesn't exist or if its value is not an int.
func Int(k string) int { return std.Int(k) }

// Int64E returns the standard settings' value for k as an int64. A
// SettingNotFoundError is returned if k doesn't exist in. A DataTypeError will
// be returned if the value is neither an int64 nor an int.
func Int64E(k string) (int64, error) { return std.Int64E(k) }

// InterfaceE returns the standard settings' value for k as an interface{}. A
// SettingNotFoundError is returned if k doesn't exist.
func InterfaceE(k string) (interface{}, error) { return std.InterfaceE(k) }

// Interface returns the standard settings' value for k as an interface{}. A nil
// will be returned if k doesn't exist.
func Interface(k string) interface{} { return std.Interface(k) }

// Int64 returns the standard settings' value for k as an int64. A 0 will be
// returned if k doesn't exist or if its value is neither an int64 nor an int
// setting.
func Int64(k string) int64 { return std.Int64(k) }

// StringE returns the standard settings' value for k as a string. A
// SettingNotFoundError is returned if k doesn't exist. A DataTypeError will be
// returned if the value is not a string.
func StringE(k string) (string, error) { return std.StringE(k) }

// String returns the standard settings' value for k as a string. An empty
// string, "", will be returned if k doesn't exist or if its value is not a
// string.
func String(k string) string { return std.String(k) }
