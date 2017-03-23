package contour

// Only non-core settings are updateable. Flags must use Override* to update
// settings. This assumes that the lock has already been obtained by the
// caller.
func (s *Settings) update(k string, v interface{}) error {
	// if can't update, a false will also return an error explaining why.
	_, err := s.canUpdate(k)
	if err != nil {
		return err
	}
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

// UpdateCfgFile updates the set config file information.  This only sets
// the filename, the format is not changed.  This does the update
// directly because the cfg filename is a core setting; it will fail the
// canUpdate check.
//
// It is assumed that RegisterCfgFile has already been called, if it hasn't
// nothing will be done.
func UpdateCfgFile(v string) { settings.UpdateCfgFile(v) }
func (s *Settings) UpdateCfgFile(v string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.settings[s.cfgFilenameKey]
	if !ok {
		return
	}
	val.Value = v
	s.settings[s.cfgFilenameKey] = val
}
