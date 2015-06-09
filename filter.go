package contour

import (
	"fmt"

	//jww "github.com/spf13/jwalterweatherman"
)

// FilterArgs takes the passed args and filter's the flags out of them.  The
// populated flags override their settings, according to the override rules.
// Successful overrides result in the relevant Cfg setting being updated along
// with the env variable.
//
// The flagset is populated right before filtering so that the default values
// reflect current cfg settings.
//
// Any args left, after filtering, are returned to the caller.
func FilterArgs(args []string) ([]string, error) { return appCfg.FilterArgs(args) }
func (c *Cfg) FilterArgs(args []string) ([]string, error) {
	// Get the ArgFilter (flag) information and set the flagSet
	err := c.setCfgFlags()
	if err != nil {
		return nil, err
	}
	c.RWMutex.RLock()
	// Parse args for flags
	err = c.flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}
	// Get the remaining args
	cmdArgs := c.flagSet.Args()
	c.RWMutex.RUnlock()
	// Process the captured values
	for _, v := range c.boolFilterNames {
		if v == "" {
			continue
		}
		c.RWMutex.RLock()
		ptr, ok := c.filterVars[v]
		c.RWMutex.RUnlock()
		if !ok {
			continue
		}
		// see if the pointer is nil, if it is, the flag wasn't present
		if ptr == nil {
			continue
		}
		s, err := c.GetBoolE(v)
		if err != nil {
			return nil, fmt.Errorf("unable to process arg filter %s: %s", v, err.Error())
		}
		var b bool
		switch ptr.(type) {
		case bool:
			b = ptr.(bool)
		case *bool:
			b = *ptr.(*bool)
		}
		if s != b {
			c.Override(v, b)
		}
	}
	for _, v := range c.intFilterNames {
		if v == "" {
			continue
		}
		c.RWMutex.RLock()
		ptr, ok := c.filterVars[v]
		c.RWMutex.RUnlock()
		if !ok {
			continue
		}
		// see if the pointer is nil, if it is, the flag wasn't present
		if ptr == nil {
			continue
		}
		s, err := c.GetIntE(v)
		if err != nil {
			return nil, fmt.Errorf("unable to process arg filter %s: %s", v, err.Error())
		}
		var i int
		switch ptr.(type) {
		case int:
			i = ptr.(int)
		case *int:
			i = *ptr.(*int)
		}
		if s != i {
			c.Override(v, i)
		}
	}
	for _, v := range c.stringFilterNames {
		if v == "" {
			continue
		}
		c.RWMutex.RLock()
		ptr, ok := c.filterVars[v]
		c.RWMutex.RUnlock()
		if !ok {
			continue
		}
		// see if the pointer is nil, if it is, the flag wasn't present
		if ptr == nil {
			continue
		}
		s, err := c.GetStringE(v)
		if err != nil {
			return nil, fmt.Errorf("unable to process arg filter %s: %s", v, err.Error())
		}
		var st string
		switch ptr.(type) {
		case string:
			st = ptr.(string)
		case *string:
			st = *ptr.(*string)
		}
		if s != st {
			c.Override(v, st)
		}
	}
	c.RWMutex.Lock()
	c.argsFiltered = true
	c.RWMutex.Unlock()
	return cmdArgs, nil
}

// setCfgFlags set's up the argFilter information while setting the Cfg's
// flagset.
func (c *Cfg) setCfgFlags() error {
	// Get the flag filters from the config variable information.
	c.RWMutex.Lock()
	for _, s := range c.settings {
		if s.IsFlag {
			switch s.Type {
			case "bool":
				c.filterVars[s.Name] = c.flagSet.Bool(s.Name, s.Value.(bool), s.Usage)
				c.boolFilterNames = append(c.boolFilterNames, s.Name)
				if s.Short != "" {
					c.filterVars[s.Short] = c.flagSet.Bool(s.Short, s.Value.(bool), s.Usage)
					c.boolFilterNames = append(c.boolFilterNames, s.Short)
				}
			case "int":
				c.filterVars[s.Name] = c.flagSet.Int(s.Name, s.Value.(int), s.Usage)
				c.intFilterNames = append(c.intFilterNames, s.Name)
				if s.Short != "" {
					c.filterVars[s.Short] = c.flagSet.Int(s.Short, s.Value.(int), s.Usage)
					c.intFilterNames = append(c.intFilterNames, s.Short)
				}
			case "int64":
				c.filterVars[s.Name] = c.flagSet.Int64(s.Name, s.Value.(int64), s.Usage)
				c.int64FilterNames = append(c.int64FilterNames, s.Name)
				if s.Short != "" {
					c.filterVars[s.Short] = c.flagSet.Int64(s.Short, s.Value.(int64), s.Usage)
					c.int64FilterNames = append(c.int64FilterNames, s.Short)
				}
			case "string":
				c.filterVars[s.Name] = c.flagSet.String(s.Name, s.Value.(string), s.Usage)
				c.stringFilterNames = append(c.stringFilterNames, s.Name)
				if s.Short != "" {
					c.filterVars[s.Short] = c.flagSet.String(s.Short, s.Value.(string), s.Usage)
					c.stringFilterNames = append(c.stringFilterNames, s.Short)
				}
			}
		}
	}
	c.RWMutex.Unlock()
	return nil
}

// Filter Methods obtain a list of flags of the filter type, e.g. boolFilter
// for bool flags, and returns them.
// GetBoolFilterNames returns a list of filter names (flags).
func GetBoolFilterNames() []string { return appCfg.GetBoolFilterNames() }
func (c *Cfg) GetBoolFilterNames() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	var names []string
	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, k)
		}
	}
	return names
}

// GetIntFilterNames returns a list of filter names (flags).
func GetIntFilterNames() []string { return appCfg.GetBoolFilterNames() }
func (c *Cfg) GetIntFilterNames() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	var names []string
	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, k)
		}
	}
	return names
}

// GetInt64FilterNames returns a list of filter names (flags).
func GetInt64FilterNames() []string { return appCfg.GetBoolFilterNames() }
func (c *Cfg) GetInt64FilterNames() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	var names []string
	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "int64" {
			names = append(names, k)
		}
	}
	return names
}

// GetStringFilterNames returns a list of filter names (flags).
func GetStringFilterNames() []string { return appCfg.GetBoolFilterNames() }
func (c *Cfg) GetStringFilterNames() []string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	var names []string
	for k, setting := range c.settings {
		if setting.IsFlag && setting.Type == "string" {
			names = append(names, k)
		}
	}
	return names
}
