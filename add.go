package contour

import (
	"fmt"
	"strconv"
)

// Core settings are not overridable via a configuration file, env vars, or
// command-line flags. They cannot be modified in any way once they have been
// registered.

// AddBoolCore adds a Core bool setting to the settings with the key k and
// value v. The value of this setting cannot be changed once it is added. If a
// setting with the same name, k, exists, a SettingExistsErr will be returned.
// If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) AddBoolCore(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addBoolCore(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addBoolCore(k string, v bool) error {
	return s.addCoreSetting(_bool, k, v, strconv.FormatBool(v))
}

// AddIntCore adds a Core int setting to the settings with the key k and value
// v. The value of this setting cannot be changed once it is added. If a
// setting with the same name, k, exists, a SettingExistsErr will be returned.
// If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) AddIntCore(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addIntCore(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addIntCore(k string, v int) error {
	return s.addCoreSetting(_int, k, v, strconv.Itoa(v))
}

// AddInt64Core adds a Core int64 setting to the settings with the key k and
// value v. The value of this setting cannot be changed once it is added. If a
// setting with the same name, k, exists, a SettingExistsErr will be returned.
// If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) AddInt64Core(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addInt64Core(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addInt64Core(k string, v int64) error {
	return s.addCoreSetting(_int64, k, v, strconv.FormatInt(v, 10))
}

// AddInterfaceCore adds a Core interface{} setting to the settings with the
// key k and value v. The value of this setting cannot be changed once it is
// added. If a setting with the same name, k, exists, a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) AddInterfaceCore(k string, v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addInterfaceCore(k, v)
}

func (s *Settings) addInterfaceCore(k string, v interface{}) error {
	return s.addCoreSetting(_interface, k, v, fmt.Sprintf("%v", v))
}

// AddStringCore adds a Core string setting to the settings with the key k and
// value v. The value of this setting cannot be changed once it is added. If a
// setting with the same name, k, exists, a SettingExistsErr will be returned.
// If k is empty, an ErrNoSettingName will be returned.
func (s *Settings) AddStringCore(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addStringCore(k, v)
}

func (s *Settings) addStringCore(k, v string) error {
	return s.addCoreSetting(_string, k, v, v)
}

func (s *Settings) addCoreSetting(typ dataType, k string, v interface{}, dflt string) error {
	return s.registerSetting(Core, typ, k, "", v, dflt, "", true, false, false, false)
}

// AddBool adds a bool setting to the settings with the key k and value f. This
// can be only be updated using the Update functions. If a setting with the
// same name, k, exists, a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned
func (s *Settings) AddBool(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addBool(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addBool(k string, v bool) error {
	return s.addSetting(_bool, k, v, strconv.FormatBool(v))
}

// AddInt adds an int setting to the settings with the key k and value f. This
// can be only be updated using the Update functions. If a setting with the
// same name, k, exists, a SettingExistsErr will be returned. If k is empty,
// an ErrNoSettingName will be returned
func (s *Settings) AddInt(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addInt(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addInt(k string, v int) error {
	return s.addSetting(_int, k, v, strconv.Itoa(v))
}

// AddInt64 adds an int64 setting to the settings with the key k and value f.
// This can be only updated using the Update functions. If a setting with the
// same name, k, exists, a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned
func (s *Settings) AddInt64(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addInt64(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addInt64(k string, v int64) error {
	return s.addSetting(_int64, k, v, strconv.FormatInt(v, 10))
}

// AddInterface adds an interface{} setting to the settings with the key k and
// value v. This can be updated using the Update functions. If a setting with
// the same name, k, exists, a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned
func (s *Settings) AddInterface(k string, v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addInterface(k, v)
}

// assumes the lock has been obtained.
func (s *Settings) addInterface(k string, v interface{}) error {
	return s.addSetting(_interface, k, v, fmt.Sprintf("%v", v))
}

// AddString adds a string setting to the settings with the key k and value f.
// This can be updated using the Update functions. If a setting with the same
// name, k, exists, a SettingExistsErr will be returned. If k is empty, an
// ErrNoSettingName will be returned
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
	return s.registerSetting(Basic, typ, k, "", v, dflt, "", false, false, false, false)
}

// AddBoolCore adds a Core bool setting to the standard settings with the key k
// and value v. The value of this setting cannot be changed once it is added.
// If a setting with the same name, k, exists, a SettingExistsErr will be
// returned. If k is empty, an ErrNoSettingName will be returned.
func AddBoolCore(k string, v bool) error { return std.AddBoolCore(k, v) }

// AddIntCore adds a Core int setting to the standard settings with the key k
// and value v. The value of this setting cannot be changed once it is added.
// If a setting with the same name, k, exists, a SettingExistsErr will be
// returned. If k is empty, an ErrNoSettingName will be returned.
func AddIntCore(k string, v int) error { return std.AddIntCore(k, v) }

// AddInt64Core adds a Core int64 setting to the standard settings with the key
// k and value v. The value of this setting cannot be changed once it is added.
// If a setting with the same name, k, exists, a SettingExistsErr will be
// returned. If k is empty, an ErrNoSettingName will be returned.
func AddInt64Core(k string, v int64) error { return std.AddInt64Core(k, v) }

// AddInterfaceCore adds a Core interface setting to the standard settings with
// the key k and value v. The value of this setting cannot be changed once it
// is added. If a setting with the same name, k, exists, a SettingExistsErr
// will be returned. If k is empty, an ErrNoSettingName will be returned.
func AddInterfaceCore(k string, v interface{}) error { return std.AddInterfaceCore(k, v) }

// AddStringCore adds a Core string setting to the standard settings with the
// key k and value v. The value of this setting cannot be changed once it is
// added. If a setting with the same name, k, exists, a SettingExistsErr will
// be returned. If k is empty, an ErrNoSettingName will be returned.
func AddStringCore(k, v string) error { return std.AddStringCore(k, v) }

// AddBool adds a bool setting to the standard settings with the key k and
// value f. This can be only be updated using the Update functions. If a
// setting with the same name, k, exists, a SettingExistsErr will be returned.
// If k is empty, an ErrNoSettingName will be returned
func AddBool(k string, v bool) error { return std.AddBool(k, v) }

// AddInt adds an int setting to the standard settings with the key k and
// value f. This can be only be updated using the Update functions. If a
// setting with the same name, k, exists, a SettingExistsErr will be returned.
// If k is empty, an ErrNoSettingName will be returned
func AddInt(k string, v int) error { return std.AddInt(k, v) }

// AddInt64 adds an int64 setting to the standard settings with the key k and
// value f. This can be only updated using the Update functions. If a setting
// with the same name, k, exists, a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned
func AddInt64(k string, v int64) error { return std.AddInt64(k, v) }

// AddInterface adds an interface setting to the standard settings with the key
// k and value f. This can be updated using the Update functions. If a setting
// with the same name, k, exists, a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned
func AddInterface(k string, v interface{}) error { return std.AddInterface(k, v) }

// AddString adds a string setting to the standard settings with the key k and
// value f. This can be updated using the Update functions. If a setting with
// the same name, k, exists, a SettingExistsErr will be returned. If k is
// empty, an ErrNoSettingName will be returned
func AddString(k, v string) error { return std.AddString(k, v) }
