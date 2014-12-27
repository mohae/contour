package contour

import (
	"fmt"
	"log"
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

	c.lock.Lock()
	defer c.lock.Unlock()

	// If it can't be overriden,
	if c.settings[k].IsCore || !c.settings[k].IsFlag {
		err := fmt.Errorf("%v: setting is not a flag. Only flags can be overridden", k)
		log.Print(err)
		return err
	}

	/*
		// Write to environment variable
		err := c.Setenv(k, v)
		if err != nil {
			logger.Error(err)
			return err
		}
	*/

	c.settings[k].Value = v
	return nil
}

func Override(k string, v interface{}) error {
	return appCfg.Override(k, v)
}
