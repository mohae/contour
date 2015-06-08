package contour

import (
	"fmt"

	jww "github.com/spf13/jwalterweatherman"
)

// FilterArgs takes the passed args and filter's the flags out of them.  The
// populated flags override their settings, according to the override rules.
// Successful overrides result in the relevant Cfg setting being updated along
// with the env variable.
//
// Any args left, after filtering, are returned to the caller.
func FilterArgs(args []string) ([]string, error) { return appCfg.FilterArgs(args) }
func (c *Cfg) FilterArgs(args []string) ([]string, error) {
	// Preallocate the worst case scenario.
	// Get the flag filters from the config variable information.
	boolFilterNames := c.GetBoolFilterNames()
	boolFilters := make([]*bool, 0, len(boolFilterNames))
	bFilterNames := make([]string, 0, len(boolFilterNames))
	c.RWMutex.Lock()
	for _, name := range boolFilterNames {
		s, _ := c.settings[name]
		if s.IsFlag {
			jww.FEEDBACK.Printf("%s is a bool flag, append it\n", name)
			boolFilters = append(boolFilters, c.flagSet.Bool(name, s.Value.(bool), s.Usage))
			bFilterNames = append(bFilterNames, name)
			if s.Short != "" {
				boolFilters = append(boolFilters, c.flagSet.Bool(s.Short, s.Value.(bool), s.Usage))
				bFilterNames = append(bFilterNames, s.Short)
			}
		}
	}
	c.RWMutex.Unlock()
	// Get the flag filters from the config variable information.
	intFilterNames := c.GetIntFilterNames()
	// Preallocate the worst case scenario.
	intFilters := make([]*int, 0, len(intFilterNames))
	iFilterNames := make([]string, 0, len(intFilterNames))
	c.RWMutex.Lock()
	for _, name := range intFilterNames {
		s, _ := c.settings[name]
		if s.IsFlag {
			jww.FEEDBACK.Printf("%s is a int flag, append it\n", name)
			intFilters = append(intFilters, c.flagSet.Int(name, s.Value.(int), s.Usage))
			iFilterNames = append(iFilterNames, name)
			if s.Short != "" {
				intFilters = append(intFilters, c.flagSet.Int(s.Short, s.Value.(int), s.Usage))
				iFilterNames = append(iFilterNames, s.Short)
			}
		}
	}
	c.RWMutex.Unlock()
	// Get the flag filters from the config variable information.
	stringFilterNames := c.GetStringFilterNames()
	// Preallocate the worst case scenario.
	stringFilters := make([]*string, 0, len(stringFilterNames))
	sFilterNames := make([]string, 0, len(stringFilterNames))
	c.RWMutex.Lock()
	for _, name := range stringFilterNames {
		s, _ := c.settings[name]
		if s.IsFlag {
			jww.FEEDBACK.Printf("%s is a string flag, append it\n", name)
			stringFilters = append(stringFilters, c.flagSet.String(name, s.Value.(string), s.Usage))
			sFilterNames = append(sFilterNames, name)
			if s.Short != "" {
				stringFilters = append(stringFilters, c.flagSet.String(s.Short, s.Value.(string), s.Usage))
				sFilterNames = append(sFilterNames, s.Short)
			}
		}
	}
	// Parse args for flags
	err := c.flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}
	// Get the remaining args
	cmdArgs := c.flagSet.Args()
	jww.FEEDBACK.Printf("cmdArgs: %v\n", cmdArgs)
	c.RWMutex.Unlock()
	// Process the captured values
	for i, v := range boolFilters {
		if v == nil {
			jww.CRITICAL.Println("%v was nil", bFilterNames[i])
			continue
		}
		c.RWMutex.RLock()
		s := c.settings[bFilterNames[i]].Value
		c.RWMutex.RUnlock()
		if s != v {
			c.Override(bFilterNames[i], v)
		}
	}
	for i, v := range intFilters {
		if v == nil {
			jww.CRITICAL.Println("%v was nil", iFilterNames[i])
			continue
		}
		c.RWMutex.RLock()
		s := c.settings[iFilterNames[i]].Value
		c.RWMutex.RUnlock()
		if s != v {
			c.Override(iFilterNames[i], v)
		}
	}
	for i, v := range stringFilters {
		if v == nil {
			jww.CRITICAL.Println("%v was nil", sFilterNames[i])
			continue
		}
		c.RWMutex.RLock()
		s := c.settings[sFilterNames[i]].Value
		c.RWMutex.RUnlock()
		if s != v {
			c.Override(sFilterNames[i], v)
		}
	}
	jww.FEEDBACK.Printf("boolFilters: %+v\n", bFilterNames)
	jww.FEEDBACK.Printf("intFilters: %+v\n", iFilterNames)
	jww.FEEDBACK.Printf("stringFilters: %+v\n", sFilterNames)
	c.RWMutex.Lock()
	c.argsFiltered = true
	c.RWMutex.Unlock()
	return cmdArgs, nil
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
