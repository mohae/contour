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
func (c *Cfg) Override(k string, v interface{}) error {
	if v == nil {
		return nil
	}
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	// If it can't be overriden,
	s, ok := c.settings[k]
	if !ok {
		return fmt.Errorf("%s not found: cannot override")
	}
	if s.IsCore || !s.IsFlag {
		return fmt.Errorf("%s is not a flag: only flags can be overridden", k)
	}
	s.Value = v
	c.settings[k] = s
	return nil
}

func Override(k string, v interface{}) error {
	return appCfg.Override(k, v)
}
