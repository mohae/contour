package contour

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

// UpdateBoolE updates a bool setting, returning any error that occurs.
func UpdateBoolE(k string, v bool) error { return settings.UpdateBoolE(k, v) }
func (s *Settings) UpdateBoolE(k string, v bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateBoolE(k, v)
}

func (s *Settings) updateBoolE(k string, v bool) error {
	return s.update(k, v)
}

// UpdateBool calls UpdateBoolE and drops the error.
func UpdateBool(k string, v bool) { settings.UpdateBool(k, v) }
func (s *Settings) UpdateBool(k string, v bool) {
	s.UpdateBoolE(k, v)
}

// UpdateIntE updates a int setting, returning any error that occurs.
func UpdateIntE(k string, v int) error { return settings.UpdateIntE(k, v) }
func (s *Settings) UpdateIntE(k string, v int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateIntE(k, v)
}

func (s *Settings) updateIntE(k string, v int) error {
	return s.update(k, v)
}

// UpdateInt calls UpdateIntE and drops the error.
func UpdateInt(k string, v int) { settings.UpdateInt(k, v) }
func (s *Settings) UpdateInt(k string, v int) {
	s.UpdateIntE(k, v)
}

// UpdateInt64E updates a int64 setting, returning any error that occurs.
func UpdateInt64E(k string, v int64) error { return settings.UpdateInt64E(k, v) }
func (s *Settings) UpdateInt64E(k string, v int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateInt64E(k, v)
}

func (s *Settings) updateInt64E(k string, v int64) error {
	return s.update(k, v)
}

// UpdateInt64 calls UpdateInt64E and drops the error.
func UpdateInt64(k string, v int64) { settings.UpdateInt64(k, v) }
func (s *Settings) UpdateInt64(k string, v int64) {
	s.UpdateInt64E(k, v)
}

// UpdateStringE updates a string setting, returning any error that occurs.
func UpdateStringE(k string, v string) error { return settings.UpdateStringE(k, v) }
func (s *Settings) UpdateStringE(k, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateStringE(k, v)
}

func (s *Settings) updateStringE(k, v string) error {
	return s.update(k, v)
}

// UpdateBool calls UpdateStringE and drops the error.
func UpdateString(k string, v string) { settings.UpdateString(k, v) }
func (s *Settings) UpdateString(k, v string) {
	s.UpdateStringE(k, v)
}

// canUpdate checks to see if the passed setting key is updateable. If the key
// doesn't exist, a false will be returned. This assumes that the lock has
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
// checked.
func (s *Settings) canUpdate(typ SettingType, k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be updated.
	v, ok := s.settings[k]
	if !ok {
		return false
	}
	// See if there are any settings that prevent it from being overridden.  Core and
	// environment variables are never settable. Core must be set during registration.
	if v.IsCore {
		return false
	}
	// regular, Basic, settings are always updateable as long as this is a Basic
	// update.
	if !v.IsCore && !v.IsEnv && !v.IsFlag && typ == Basic {
		return true
	}

	// check by update type
	switch typ {
	case ConfFileVar:
		if v.IsConfFileVar {
			if !s.confFileVarsSet && !s.envSet && !s.flagsParsed {
				return true
			}
		}
		return false
	case Env:
		if v.IsEnv {
			if !s.envSet && !s.flagsParsed {
				return true
			}
		}
		return false
	case Flag:
		if v.IsFlag && !s.flagsParsed {
			return true
		}
		return false
	}
	// If it was not one of the above, we return a false. It's better to not allow
	// an update if the case isn't handled than be too permissive. Getting here is
	// a sign that something within this func should be updated and/or fixed.
	return false
}
