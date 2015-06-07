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
	var flags int // counter, gets reset with each type
	// Preallocate the worst case scenario.
	// Get the flag filters from the config variable information.
	boolFilterNames := c.GetBoolFilterNames()
	boolFilters := make([]*bool, len(boolFilterNames))
	bFilterNames := make([]string, len(boolFilterNames))
	c.RWMutex.Lock()
	for _, name := range boolFilterNames {
		if c.settings[name].IsFlag {
			boolFilters[flags] = c.flagSet.Bool(name, c.settings[name].Value.(bool), c.settings[name].Usage)
			bFilterNames[flags] = name
			flags++
		}
	}
	c.RWMutex.Unlock()
	// Get the flag filters from the config variable information.
	intFilterNames := c.GetIntFilterNames()
	// Preallocate the worst case scenario.
	intFilters := make([]*int, len(intFilterNames))
	iFilterNames := make([]string, len(intFilterNames))
	flags = 0
	c.RWMutex.Lock()
	for _, name := range intFilterNames {
		if c.settings[name].IsFlag {
			intFilters[flags] = c.flagSet.Int(name, c.settings[name].Value.(int), c.settings[name].Usage)
			iFilterNames[flags] = name
			flags++
		}
	}
	c.RWMutex.Unlock()
	// Get the flag filters from the config variable information.
	stringFilterNames := c.GetStringFilterNames()
	// Preallocate the worst case scenario.
	stringFilters := make([]*string, len(stringFilterNames))
	sFilterNames := make([]string, len(stringFilterNames))
	flags = 0
	c.RWMutex.Lock()
	for _, name := range stringFilterNames {
		if c.settings[name].IsFlag {
			stringFilters[flags] = c.flagSet.String(name, c.settings[name].Value.(string), c.settings[name].Usage)
			sFilterNames[flags] = name
			flags++
		}
	}
	// Parse args for flags
	err := c.flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}
	// Get the remaining args
	cmdArgs := c.flagSet.Args()
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
			Override(bFilterNames[i], v)
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
			Override(iFilterNames[i], v)
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
			Override(sFilterNames[i], v)
		}
	}
	c.RWMutex.Lock()
	c.flagsSet = true
	c.RWMutex.Unlock()
	return cmdArgs, nil
}
