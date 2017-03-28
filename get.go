package contour

import "reflect"

// Get functions and methods.
//
// E versions return an error if one occurs. Non-E versions return the zero
// value if an error occurs.

// GetE returns the key's Value as an interface{}. An SettingNotFoundErr is
// returned if the key doesn't exist.
func GetE(k string) (interface{}, error) { return settings.GetE(k) }

// GetE returns the key's Value as an interface{}. An SettingNotFoundErr is
// returned if the key doesn't exist.
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
		return nil, SettingNotFoundErr{name: k}
	}
	return s.settings[k].Value, nil
}

// Get returns the key's value as an interface{}. A nil is returned if the key
// doesn't exist.
func Get(k string) interface{} { return settings.Get(k) }

// Get returns the key's value as an interface{}. A nil is returned if the key
// doesn't exist.
func (s *Settings) Get(k string) interface{} {
	v, _ := s.GetE(k)
	return v
}

// BoolE returns the key's value as a bool. A SettingNotFoundErr is returned
// if the key doesn't exist. A DataTypeErr will be returned if the setting's
// type is not bool.
func BoolE(k string) (bool, error) { return settings.BoolE(k) }

// BoolE returns the key's value as a bool. A SettingNotFoundErr is returned
// if the key doesn't exist. A DataTypeErr will be returned if the setting's
// type is not bool.
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
	return false, DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _bool}
}

// Bool returns the key's value as a bool. A false will be returned if the key
// either doesn't exist or the setting's type is not bool.
func Bool(k string) bool { return settings.Bool(k) }

// Bool returns the key's value as a bool. A false will be returned if the key
// either doesn't exist or the setting's type is not bool.
func (s *Settings) Bool(k string) bool {
	v, _ := s.BoolE(k)
	return v
}

// IntE returns the key's value as an int. A SettingNotFoundErr is returned if
// the key doesn't exist. A DataTypeErr will be returned if the setting's type
// is not int.
func IntE(k string) (int, error) { return settings.IntE(k) }

// IntE returns the key's value as an int. A SettingNotFoundErr is returned if
// the key doesn't exist. A DataTypeErr will be returned if the setting's type
// is not int.
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
	return 0, DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _int}
}

// Int returns the key's value as an int. A 0 will be returned if the key
// either doesn't exist or is not an int setting.
func Int(k string) int { return settings.Int(k) }

// Int returns the key's value as an int. A 0 will be returned if the key
// either doesn't exist or is not an int setting.
func (s *Settings) Int(k string) int {
	v, _ := s.IntE(k)
	return v
}

// Int64E returns the key's value as an int64. A SettingNotFoundErr is returned
// if the key doesn't exist. A DataTypeErr will be returned if the setting's
// type is neither an int64 nor an int.
func Int64E(k string) (int64, error) { return settings.Int64E(k) }

// Int64E returns the key's value as an int64. A SettingNotFoundErr is returned
// if the key doesn't exist. A DataTypeErr will be returned if the setting's
// type is neither an int64 nor an int.
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
	return 0, DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _int64}
}

// Int64 returns the key's value as an int64. A 0 will be returned if the key
// either doesn't exist or is neither an int64 nor an int setting.
func Int64(k string) int64 { return settings.Int64(k) }

// Int64 returns the key's value as an int64. A 0 will be returned if the key
// either doesn't exist or is neither an int64 nor an int setting.
func (s *Settings) Int64(k string) int64 {
	v, _ := s.Int64E(k)
	return v
}

// StringE returns the key's value as a string. A SettingNotFoundErr is
// returned if the key doesn't exist. A DataTypeErr will be returned if the
// setting's type is not a string.
func StringE(k string) (string, error) { return settings.StringE(k) }

// StringE returns the key's value as a string. A SettingNotFoundErr is
// returned if the key doesn't exist. A DataTypeErr will be returned if the
// setting's type is not a string.
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
	return "", DataTypeErr{name: k, is: reflect.TypeOf(v).String(), not: _string}
}

// String returns the key's value as a string. An empty string, "", will be
// returned if the key either doesn't exist or is not a string setting.
func String(k string) string { return settings.String(k) }

// String returns the key's value as a string. An empty string, "", will be
// returned if the key either doesn't exist or is not a string setting.
func (s *Settings) String(k string) string {
	v, _ := s.StringE(k)
	return v
}

// ConfFilename returns the configuration filename.
func ConfFilename() string { return settings.ConfFilename() }
func (s *Settings) ConfFilename() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.confFilename
}
