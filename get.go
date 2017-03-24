package contour

import "reflect"

// Get functions and methods.
//
// E versions, these return the error. Non-e versions are just wrapped calls to
// these functions with the error dropped.

// GetE returns the setting Value as an interface{}. If its not a valid
// setting, an error is returned.
func GetE(k string) (interface{}, error) { return settings.GetE(k) }
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

func Get(k string) interface{} { return settings.Get(k) }
func (s *Settings) Get(k string) interface{} {
	v, _ := s.GetE(k)
	return v
}

// BoolE returns the key's value as a bool. A SettingNotFoundErr is returned
// if the key is not valid. If the setting's type is not a bool, a DataTypeErr
// will be returned.
func BoolE(k string) (bool, error) { return settings.BoolE(k) }
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
// either doesn't exist or is not a bool setting.
func Bool(k string) bool { return settings.Bool(k) }
func (s *Settings) Bool(k string) bool {
	v, _ := s.BoolE(k)
	return v
}

// IntE returns the key's value as an int. A SettingNotFoundErr is returned if
// the key is not valid. If the setting's type is not an int, a DataTypeErr
// will be returned.
func IntE(k string) (int, error) { return settings.IntE(k) }
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
func (s *Settings) Int(k string) int {
	v, _ := s.IntE(k)
	return v
}

// Int64E returns the key's value as an int64. A SettingNotFoundErr is returned
// if the key is not valid. If the setting's type is neither an int64 nor an
// int, a DataTypeErr will be returned.
func Int64E(k string) (int64, error) { return settings.Int64E(k) }
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
func (s *Settings) Int64(k string) int64 {
	v, _ := s.Int64E(k)
	return v
}

// StringE returns the key's value as a string. A SettingNotFoundErr is
// returned if the key is not gvalid. If the setting's type is not a string, a
// DataTypeErr will be returned.
func StringE(k string) (string, error) { return settings.StringE(k) }
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
func (s *Settings) String(k string) string {
	v, _ := s.StringE(k)
	return v
}

// ConfFilename returns the configuration filename and its format. If the key
// is not registered, or if the format isn't a supported format, an error is
// returned and the format will be Unsupported.
func ConfFilename() (name string, format Format, err error) { return settings.ConfFilename() }
func (s *Settings) ConfFilename() (name string, format Format, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.settings[s.confFilenameVarName]
	if !ok {
		return "", Unsupported, SettingNotFoundErr{Core, s.confFilenameVarName}
	}
	format, _ = ParseFilenameFormat(v.Name)
	return v.Name, format, nil
}
