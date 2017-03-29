package contour

import "strconv"

// Core settings are not overridable via a configuration file, env vars, or
// command-line flags. They cannot be modified in any way once they have been
// registered.

// AddBoolCore adds a Core bool setting with the key k and value v. The value
// of this setting cannot be changed once it is added. If a setting with the
// same name, k,
func AddBoolCore(k string, v bool) error { return settings.AddBoolCore(k, v) }

// AddBoolCore adds a Core bool setting with the key k and value v. The value
// of this setting cannot be changed once it is added. If an error occurs, a
// SettingErr will be returned.
func (s *Settings) AddBoolCore(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addBoolCore(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addBoolCore(k string, v bool) error {
	return s.addCoreSetting(_bool, k, v, strconv.FormatBool(v))
}

// AddIntCore adds a Core int setting with the key k and value v. The value
// of this setting cannot be changed once it is added. If an error occurs, a
// SettingErr will be returned.
func AddIntCore(k string, v int) error { return settings.AddIntCore(k, v) }

// AddIntCore adds a Core int setting with the key k and value v. The value
// of this setting cannot be changed once it is added. If an error occurs, a
// SettingErr will be returned.
func (s *Settings) AddIntCore(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addIntCore(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addIntCore(k string, v int) error {
	return s.addCoreSetting(_int, k, v, strconv.Itoa(v))
}

// AddInt64Core adds a Core int64 setting with the key k and value v. The value
// of this setting cannot be changed once it is added. If an error occurs, a
// SettingErr will be returned.
func AddInt64Core(k string, v int64) error { return settings.AddInt64Core(k, v) }

// AddInt64Core adds a Core int64 setting with the key k and value v. The value
// of this setting cannot be changed once it is added. If an error occurs, a
// SettingErr will be returned.
func (s *Settings) AddInt64Core(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addInt64Core(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addInt64Core(k string, v int64) error {
	return s.addCoreSetting(_int64, k, v, strconv.FormatInt(v, 10))
}

// AddStringCore adds a Core string setting with the key k and value v. The
// value of this setting cannot be changed once it is added. If an error
// occurs, a SettingErr will be returned.
func AddStringCore(k, v string) error { return settings.AddStringCore(k, v) }

// AddStringCore adds a Core string setting with the key k and value v. The
// value of this setting cannot be changed once it is added. If an error
// occurs, a SettingErr will be returned.
func (s *Settings) AddStringCore(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addStringCore(k, v)
}

func (s *Settings) addStringCore(k, v string) error {
	return s.addCoreSetting(_string, k, v, v)
}

func (s *Settings) addCoreSetting(typ dataType, k string, v interface{}, dflt string) error {
	return s.registerSetting(typ, k, "", v, dflt, "", true, false, false, false)
}

// AddBool adds a bool setting with they key k and value f. This can be only be
// updated using the Update functions. If an error occurs, a SettingErr will
// returned.
func AddBool(k string, v bool) error { return settings.AddBool(k, v) }

// AddBool adds a bool setting with they key k and value f. This can be only be
// updated using the Update functions. If an error occurs, a SettingErr will
// returned.
func (s *Settings) AddBool(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addBool(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addBool(k string, v bool) error {
	return s.addSetting(_bool, k, v, strconv.FormatBool(v))
}

// AddInt adds an int setting with they key k and value f. This can be only be
// updated using the Update functions. If an error occurs, a SettingErr will
// returned.
func AddInt(k string, v int) error { return settings.AddInt(k, v) }

// AddInt adds an int setting with they key k and value f. This can be only be
// updated using the Update functions. If an error occurs, a SettingErr will
// returned.
func (s *Settings) AddInt(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addInt(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addInt(k string, v int) error {
	return s.addSetting(_int, k, v, strconv.Itoa(v))
}

// AddInt64 adds an int64 setting with they key k and value f. This can be only
// be updated using the Update functions. If an error occurs, a SettingErr will
// returned.
func AddInt64(k string, v int64) error { return settings.AddInt64(k, v) }

// AddInt64 adds an int64 setting with they key k and value f. This can be only
// be updated using the Update functions. If an error occurs, a SettingErr will
// returned.
func (s *Settings) AddInt64(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addInt64(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addInt64(k string, v int64) error {
	return s.addSetting(_int64, k, v, strconv.FormatInt(v, 10))
}

// AddString adds a string setting with they key k and value f. This can be
// only be updated using the Update functions. If an error occurs, a SettingErr
// will returned.
func AddString(k, v string) error { return settings.AddString(k, v) }

// AddString adds a string setting with they key k and value f. This can be
// only be updated using the Update functions. If an error occurs, a SettingErr
// will returned.
func (s *Settings) AddString(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addString(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addString(k, v string) error {
	return s.addSetting(_string, k, v, v)
}

func (s *Settings) addSetting(typ dataType, k string, v interface{}, dflt string) error {
	return s.registerSetting(typ, k, "", v, dflt, "", false, false, false, false)
}
