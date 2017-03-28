package contour

import "fmt"

type UpdateErr struct {
	k    string
	typ  SettingType
	slug string
}

func (e UpdateErr) Error() string {
	return fmt.Sprintf("update of %s failed: %s", e.k, e.slug)
}

type BasicUpdateErr struct {
	typ SettingType
	k   string
}

func (e BasicUpdateErr) Error() string {
	return fmt.Sprintf("%s: %s settings cannot be updated by a basic update", e.k, e.typ)
}

type CoreUpdateErr struct {
	typ SettingType
	k   string
}

func (e CoreUpdateErr) Error() string {
	return fmt.Sprintf("%s: %s settings cannot be updated by a core update", e.k, e.typ)
}

// Only non-core settings are updateable. This assumes that the lock has
// already been obtained by the caller.
func (s *Settings) update(k string, v interface{}) error {
	// if can't update, a false will also return an error explaining why.
	//can := s.canUpdate(typ, k)
	//if !can {
	//	return fmt.Errorf("%se: cannot update using %s", typ)
	//}
	val, _ := s.settings[k]
	val.Value = v
	s.settings[k] = val
	return nil
}

// UpdateBool updates a bool setting. If the setting k doesn't exist, both a
// false and a SettingNotFoundErr will be returned. If the setting k is not
// updateable, both a false and one of the following errors will be returned:
// CoreUpdateErr. BasicUpdateErr, or UpdateErr.
func UpdateBool(k string, v bool) error { return settings.UpdateBool(k, v) }
func (s *Settings) UpdateBool(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateBool(k, v)
}

func (s *Settings) updateBool(k string, v bool) error {
	return s.update(k, v)
}

// UpdateInt updates an int setting. If the setting k doesn't exist, both a
// false and a SettingNotFoundErr will be returned. If the setting k is not
// updateable, both a false and one of the following errors will be returned:
// CoreUpdateErr. BasicUpdateErr, or UpdateErr.
func UpdateInt(k string, v int) error { return settings.UpdateInt(k, v) }
func (s *Settings) UpdateInt(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInt(k, v)
}

func (s *Settings) updateInt(k string, v int) error {
	return s.update(k, v)
}

// UpdateInt64 updates an int64 setting. If the setting k doesn't exist, both a
// false and a SettingNotFoundErr will be returned. If the setting k is not
// updateable, both a false and one of the following errors will be returned:
// CoreUpdateErr. BasicUpdateErr, or UpdateErr.
func UpdateInt64(k string, v int64) error { return settings.UpdateInt64(k, v) }
func (s *Settings) UpdateInt64(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInt64(k, v)
}

func (s *Settings) updateInt64(k string, v int64) error {
	return s.update(k, v)
}

// UpdateString updates a string setting. If the setting k doesn't exist, both
// a false and a SettingNotFoundErr will be returned. If the setting k is not
// updateable, both a false and one of the following errors will be returned:
// CoreUpdateErr. BasicUpdateErr, or UpdateErr.
func UpdateString(k string, v string) error { return settings.UpdateString(k, v) }
func (s *Settings) UpdateString(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateString(k, v)
}

func (s *Settings) updateString(k, v string) error {
	return s.update(k, v)
}

// canUpdate checks to see if the passed setting key is updateable. If the key
// doesn't exist, both a false and a SettingNotFoundErr will be returned. If
// the setting is not updateable, both a false and an Update type specific err
// will be returned, e.g. CoreUpdateErr. This assumes that the lock has already
// been obtained by the caller.
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
		return false, SettingNotFoundErr{name: k}
	}
	// See if there are any settings that prevent it from being overridden.  Core and
	// environment variables are never settable. Core must be set during registration.
	if v.IsCore {
		return false, CoreUpdateErr{k: k}
	}
	// regular, Basic, settings are always updateable as long as this is a Basic
	// update.
	if typ == Basic {
		if !v.IsConfFileVar && !v.IsEnv && !v.IsFlag {
			return true, nil
		}
		var t SettingType
		if v.IsFlag {
			t = Flag
			goto basicErr
		}
		if v.IsEnv {
			t = Env
			goto basicErr
		}
		t = ConfFileVar
	basicErr:
		return false, BasicUpdateErr{typ: t, k: k}
	}

	// check by update type
	switch typ {
	case ConfFileVar:
		if v.IsConfFileVar {
			if !s.confFileVarsSet && !s.envSet && !s.flagsParsed {
				return true, nil
			}
			var set string
			if s.flagsParsed {
				set = "flags"
				goto confErr
			}
			if s.envSet {
				set = "env vars"
				goto confErr
			}
			set = "the configuration file"
		confErr:
			return false, UpdateErr{k: k, slug: fmt.Sprintf("already set from %s", set)}
		}
		return false, UpdateErr{typ: typ, k: k, slug: fmt.Sprintf("is not a %s", ConfFileVar)}
	case Env:
		if v.IsEnv {
			if !s.envSet && !s.flagsParsed {
				return true, nil
			}
			var set string
			if s.flagsParsed {
				set = "flags"
			} else {
				set = "env vars"
			}
			return false, UpdateErr{typ: typ, k: k, slug: fmt.Sprintf("already set from %s", set)}
		}
		return false, UpdateErr{typ: typ, k: k, slug: fmt.Sprintf("is not an %s", Env)}
	case Flag:
		if v.IsFlag {
			if !s.flagsParsed {
				return true, nil
			}
			return false, UpdateErr{typ: typ, k: k, slug: "already set from flags"}
		}
		return false, UpdateErr{typ: typ, k: k, slug: fmt.Sprintf("is not a %s", Flag)}
	}
	// If it was not one of the above, we return a false. It's better to not allow
	// an update if the case isn't handled than be too permissive. Getting here is
	// a sign that something within this func should be updated and/or fixed.
	return false, UpdateErr{typ: typ, k: k, slug: "invalid update type"}
}
