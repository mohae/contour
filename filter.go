package contour

import (
	"fmt"
	// "strconv"
)

// FilterArgs is a convenience function for the global appCfg.
func FilterArgs(args []string) ([]string, error) {
	return appCfg.FilterArgs(args)
}

// FilterArgs takes the passed args and filter's the flags out of them.  The
// populated flags override their settings, according to the override rules.
// Successful overrides result in the relevant Cfg setting being updated along
// with the env variable.
//
// Any args left, after filtering, are returned to the caller.
func (c *Cfg) FilterArgs(args []string) ([]string, error) {
	// Get the flag filters from the config variable information.
	boolFilterNames := c.GetBoolFilterNames()
	var flags int // counter, gets reset with each type
	// Preallocate the worst case scenario.
	boolFilters := make([]*bool, len(boolFilterNames))
	bFilterNames := make([]string, len(boolFilterNames))
	for _, name := range boolFilterNames {
		if c.settings[name].IsFlag {
			boolFilters[flags] = c.flagSet.Bool(name, c.settings[name].Value.(bool), fmt.Sprintf("filter %s", name))
			bFilterNames[flags] = name
			flags++
		}
	}
	// Get the flag filters from the config variable information.
	intFilterNames := c.GetIntFilterNames()
	// Preallocate the worst case scenario.
	intFilters := make([]*int, len(intFilterNames))
	iFilterNames := make([]string, len(intFilterNames))
	flags = 0
	for _, name := range intFilterNames {
		if c.settings[name].IsFlag {
			intFilters[flags] = c.flagSet.Int(name, c.settings[name].Value.(int), fmt.Sprintf("filter %s", name))
			iFilterNames[flags] = name
			flags++
		}
	}
	// Get the flag filters from the config variable information.
	stringFilterNames := c.GetStringFilterNames()
	// Preallocate the worst case scenario.
	stringFilters := make([]*string, len(stringFilterNames))
	sFilterNames := make([]string, len(stringFilterNames))
	flags = 0
	for _, name := range stringFilterNames {
		if c.settings[name].IsFlag {
			stringFilters[flags] = c.flagSet.String(name, c.settings[name].Value.(string), fmt.Sprintf("filter %s", name))
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
	// Process the captured values
	for i, v := range boolFilters {
		if v != c.settings[bFilterNames[i]].Value {
			Override(bFilterNames[i], v)
		}
	}
	for i, v := range intFilters {
		if v != c.settings[iFilterNames[i]].Value {
			Override(iFilterNames[i], v)
		}
	}
	for i, v := range stringFilters {
		if v != c.settings[sFilterNames[i]].Value {
			Override(sFilterNames[i], v)
		}
	}
	c.flagsSet = true
	return cmdArgs, nil
}
