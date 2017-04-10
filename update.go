package contour

import "fmt"

// updateError is any error that happens on an update that isn't one of the
// following: SettingNotFoundError, CoreUpdateError, or UpdateError. This only
// occurs on internal contour operations.
type updateError struct {
	k    string
	typ  SettingType
	slug string
}

func (e updateError) Error() string {
	return fmt.Sprintf("update of %s failed: %s", e.k, e.slug)
}

// UpdateError records information about an Update operation that results in an
// error that isn't neither a SettingNotFoundError nor a CoreUpdateError.
type UpdateError struct {
	typ string
	k   string
}

func (e UpdateError) Error() string {
	return fmt.Sprintf("%s: %s settings cannot be updated", e.k, e.typ)
}

// CoreUpdateError happens when there's an attempt to update a Core setting.
type CoreUpdateError struct {
	k string
}

func (e CoreUpdateError) Error() string {
	return fmt.Sprintf("%s: core settings cannot be updated", e.k)
}

// Only non-core settings are updateable. This assumes that the lock has
// already been obtained by the caller.
func (s *Settings) update(typ SettingType, k string, v interface{}) error {
	// the bool is ignored because a false will always return an error.
	_, err := s.canUpdate(typ, k)
	if err != nil {
		return err
	}
	val, _ := s.settings[k]
	val.Value = v
	s.settings[k] = val
	return nil
}

// UpdateBool updates k with a boo, v. If settings does not have a setting k,
// both a false and a SettingNotFoundError will be returned. If the setting k is
// not updateable, both a false and either a CoreUpdateError or an UpdateError
// will be returned.
func (s *Settings) UpdateBool(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateBool(Basic, k, v)
}

// this assumes the lock is held by the caller.
func (s *Settings) updateBool(typ SettingType, k string, v bool) error {
	return s.update(typ, k, v)
}

// UpdateInt updates k with an int, v. If settings does not have a setting k,
// both a false and a SettingNotFoundError will be returned. If the setting k
// is not updateable, both a false and either a CoreUpdateError or an
// UpdateError will be returned.
func (s *Settings) UpdateInt(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInt(Basic, k, v)
}

// this assumes the lock is held by the caller.
func (s *Settings) updateInt(typ SettingType, k string, v int) error {
	return s.update(typ, k, v)
}

// UpdateInt64 updates k with an int64, v. If settings does not have a setting
// k, both a false and a SettingNotFoundError will be returned. If the setting
// k is not updateable, both a false and either a CoreUpdateError or an
// UpdateError will be returned.
func (s *Settings) UpdateInt64(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInt64(Basic, k, v)
}

// this assumes the lock is held by the caller.
func (s *Settings) updateInt64(typ SettingType, k string, v int64) error {
	return s.update(typ, k, v)
}

// UpdateInterface updates k with an interface, v. If settings does not have a
// setting k, both a false and a SettingNotFoundError will be returned. If the
// setting k is not updateable, both a false and either a CoreUpdateError or an
// UpdateError will be returned.
func (s *Settings) UpdateInterface(k string, v interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInterface(Basic, k, v)
}

// this assumes the lock is held by the caller.
func (s *Settings) updateInterface(typ SettingType, k string, v interface{}) error {
	return s.update(typ, k, v)
}

// UpdateString updates k with a string, v. If settings does not have a setting
// k, both a false and a SettingNotFoundError will be returned. If the setting
// k is not updateable, both a false and either a CoreUpdateError or an
// UpdateError will be returned.
func (s *Settings) UpdateString(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateString(Basic, k, v)
}

// this assumes the lock is held by the caller.
func (s *Settings) updateString(typ SettingType, k, v string) error {
	return s.update(typ, k, v)
}

// canUpdate checks to see if the passed setting key is updateable. If the key
// doesn't exist, both a false and a SettingNotFoundError will be returned. If
// the setting is not updateable, both a false and an Update type specific err
// will be returned, e.g. CoreUpdateError. This assumes that the lock has
// already been obtained by the caller.
//
// core settings can never be set.
// settings that are not a ConfFileVar, EnvVar, and Flag, i.e. a Basic setting,
//   cannot be set by a configuration file, environment variable, or a flag as
//   it has not explicitly been exposed to them. They can only be set by the
//   application code, i.e. non-contour code.
//
// all other settings, for whatever reason, may be any combination of types,
// e.g. it could be a conf var and a flag. Settings of type conf var, env var
// or flag can be set if a higher precedence type has not already been set.
//
// examples:
//    a setting that IsConfFileVar && IsFlag can be set by a ConfFile if both
//    flagsParsed and confFileVarsSet are False.
//
//    a setting that IsConfFileVar && IsFlag can be set by flags if both
//    flagsParsed is False.
//
// k is the key of the setting and typ is the type of update that is being
// checked, e.g. an update from an env var will have a typ of EnvVar.
func (s *Settings) canUpdate(typ SettingType, k string) (can bool, err error) {
	// See if the key exists, if it doesn't already exist, it can't be updated.
	v, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundError{k: k}
	}
	// See if there are any settings that prevent it from being overridden. Core and
	// environment variables are never settable. Core must be set during registration.
	if v.IsCore {
		return false, CoreUpdateError{k: k}
	}
	// regular, Basic, settings are always updateable as long as this is a Basic
	// update.
	if typ == Basic {
		if !v.IsConfFileVar && !v.IsEnvVar && !v.IsFlag {
			return true, nil
		}
		var t string
		if v.IsFlag {
			t = Flag.String()
			goto basicErr
		}
		if v.IsEnvVar {
			t = EnvVar.String()
			goto basicErr
		}
		t = "configuration file"
	basicErr:
		return false, UpdateError{typ: t, k: k}
	}

	// check by update type
	switch typ {
	case ConfFileVar:
		if v.IsConfFileVar {
			if !s.confFileVarsSet && !s.envVarsSet && !s.flagsParsed {
				return true, nil
			}
			var set string
			if s.flagsParsed {
				set = "flags"
				goto confErr
			}
			if s.envVarsSet {
				set = "env vars"
				goto confErr
			}
			set = "the configuration file"
		confErr:
			return false, updateError{k: k, slug: fmt.Sprintf("already set from %s", set)}
		}
		return false, updateError{typ: typ, k: k, slug: fmt.Sprintf("is not a %s", ConfFileVar)}
	case EnvVar:
		if v.IsEnvVar {
			if !s.envVarsSet && !s.flagsParsed {
				return true, nil
			}
			var set string
			if s.flagsParsed {
				set = "flags"
			} else {
				set = "env vars"
			}
			return false, updateError{typ: typ, k: k, slug: fmt.Sprintf("already set from %s", set)}
		}
		return false, updateError{typ: typ, k: k, slug: fmt.Sprintf("is not an %s", EnvVar)}
	case Flag:
		if v.IsFlag {
			if !s.flagsParsed {
				return true, nil
			}
			return false, updateError{typ: typ, k: k, slug: "already set from flags"}
		}
		return false, updateError{typ: typ, k: k, slug: fmt.Sprintf("is not a %s", Flag)}
	}
	// If it was not one of the above, we return a false. It's better to not allow
	// an update if the case isn't handled than be too permissive. Getting here is
	// a sign that something within this func should be updated and/or fixed.
	return false, updateError{typ: typ, k: k, slug: "invalid update type"}
}

// UpdateBool updates k with a boo, v. If the standard settings does not have a
// setting k, both a false and a SettingNotFoundError will be returned. If the
// setting k is not updateable, both a false and either a CoreUpdateError or an
// UpdateError will be returned.
func UpdateBool(k string, v bool) error { return std.UpdateBool(k, v) }

// UpdateInt updates k with aan int, v. If the standard settings does not have
// a setting k, both a false and a SettingNotFoundError will be returned. If
// the setting k is not updateable, both a false and either a CoreUpdateError
// or an UpdateError will be returned.
func UpdateInt(k string, v int) error { return std.UpdateInt(k, v) }

// UpdateInt64 updates k with an int64, v. If the standard settings does not
// have a setting k, both a false and a SettingNotFoundError will be returned.
// If the setting k is not updateable, both a false and either a
// CoreUpdateError or an UpdateError will be returned.
func UpdateInt64(k string, v int64) error { return std.UpdateInt64(k, v) }

// UpdateInterface updates k with aan int, v. If the standard settings does not
// have a setting k, both a false and a SettingNotFoundError will be returned.
// If the setting k is not updateable, both a false and either a
// CoreUpdateError or an UpdateError will be returned.
func UpdateInterface(k string, v interface{}) error { return std.UpdateInterface(k, v) }

// UpdateString updates a string setting. If the standard setting k doesn't
// exist, both a false and a SettingNotFoundError will be returned. If the
// setting k is not updateable, both a false and either a CoreUpdateError or an
// UpdateError will be returned.
func UpdateString(k string, v string) error { return std.UpdateString(k, v) }
