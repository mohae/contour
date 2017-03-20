package contour

import (
	"fmt"
)

// Set contains all of contour's Override functions. Override can set both
// settings whose values were obtained from environment variables and regular
// settings that are also Flags. Override cannot set any configuration setting
// that is not a flag. Also, override cannot set any Immutable or IsCore
// settings.
//
// A common use for overrides is to set values obtained by flags.
func Override(k string, v interface{}) error { return settings.Override(k, v) }
func (s *Settings) Override(k string, v interface{}) error {
	if v == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	// If it can't be overriden,
	st, ok := s.settings[k]
	if !ok {
		return fmt.Errorf("%s not found: cannot override", k)
	}
	if st.IsCore || !st.IsFlag {
		return fmt.Errorf("%s is not a flag: only flags can be overridden", k)
	}
	st.Value = v
	s.settings[k] = st
	return nil
}
