package contour

import (
	"fmt"

	flag "github.com/ogier/pflag"
)

// filter is a struct to hold information about a specific arg/option filter
type filter struct {
	// index of the setting for which this filter exists
	index int
	
	// holds the pointer received from flags
	value interface{}
}

// addFilter adds a filter for a flag. addFilter assumes that it is already
// inside a lock.
func (c *Cfg) addFilter(index int) {
	c.filters = append(c.filters, filter{index: index})
	c.filterCount++
}

// FilterArgs takes the passed args and filter's the flags out of them.
// The populated flags override their settings, according to the override
// rules. Successful overrides result in the relevant AppConfig setting
// being updated along with the environment variable.
//
// Any args left, after filtering, are returned to the caller.
// TODO: refactor this for greater abstraction (as long as it doesn't come
//	at the cost of reflection)
func (c *Cfg) FilterArgs(args []string) ([]string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Preallocate the worst case scenario.
//	boolFilters := make([]*bool, len(boolFilterNames))
//	bFilterNames := make([]string, len(boolFilterNames))
//	var flags int

	for _, filter := range c.filters {
		switch c.settings[filter.index].Type {
		case "string":
			
		case "bool":

		case "int":

		case "float64":

		case "int64":

		case "int8":

		case "int32":

		case "float32":

		default:
			return nil, unsupportedDataTypeErr(c.settings[filter.index].Type)
		}
	}

//	filters[i] = flagSet.BoolP(c.settings[filter.index], c.settings[filter.index].Short), c.settings[index].Value.(bool), fmt.Sprintf("filter %s", boolFilterNames[i]))
//	bFilterNames[flags] = boolFilterNames[i]
//	flags++

/*
	// Get the flag filters from the config variable information.
	intFilterNames, intFilterIndex := c.GetIntFilterNames()

	// Preallocate the worst case scenario.
	intFilters := make([]*int, len(intFilterNames))
	iFilterNames := make([]string, len(intFilterNames))
	flags = 0

	for i, idx := range intFilterIndex {
		if c.settings[idx].IsFlag {
			intFilters[flags] = flagSet.IntP(intFilterNames[i], string(c.settings[idx].Short), c.settings[idx].Value.(int), fmt.Sprintf("filter %s", intFilterNames[i]))
			iFilterNames[flags] = intFilterNames[i]
			flags++
		}
	}
	// Get the flag filters from the config variable information.
	stringFilterNames, stringFilterIndex := c.GetStringFilterNames()

	// Preallocate the worst case scenario.
	stringFilters := make([]*string, len(stringFilterNames))
	sFilterNames := make([]string, len(stringFilterNames))
	flags = 0

	for i, idx := range stringFilterIndex {
		if c.settings[idx].IsFlag {
			stringFilters[idx] = flagSet.StringP(stringFilterNames[i], string(c.settings[idx].Short), c.settings[idx].Value.(string), fmt.Sprintf("filter %s", stringFilterNames[i]))
			sFilterNames[flags] = stringFilterNames[i]
			flags++
		}
	}
*/

	// Parse args for flags
	err := c.flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}

	// Get the remaining args
	cmdArgs := c.flagSet.Args()

/*	var name string
	// Process the captured values
	for i, v := range boolFilters {
		if v != c.settings[bFilterNames[i]].Value {
			Override(bFilterNames[i], v)
		}
	}

	for i, v := range intFilters {
		name = iFilterNames[i]
		if v != c.settings[iFilterNames[i]].Value {
			Override(iFilterNames[i], v)
		}
	}

	for i, v := range stringFilters {
		name = sFilterNames[i]
		if v != c.settings[sFilterNames[i]].Value {
			Override(sFilterNames[i], v)
		}
	}
*/
	c.flagsSet = true	
	return cmdArgs, nil
}

func FilterArgs(args []string) ([]string, error) {
	return configs[0].FilterArgs(args)
}

// Filter Methods obtain a list of flags of the filter type, e.g. boolFilter
// for bool flags, and returns them.
// GetBoolFilterNames returns a list of filter names (flags).
func (c *Cfg) GetBoolFilterNames() (names []string, index []int) {
	for i, setting := range c.settings {
		if setting.IsFlag && setting.Type == "bool" {
			names = append(names, setting.Name)
			index = append(index, i)
		}
	}

	return names, index
}

// GetIntFilterNames returns a list of filter names (flags).
func (c *Cfg) GetIntFilterNames() (names []string, index []int) {
	for i, setting := range c.settings {
		if setting.IsFlag && setting.Type == "int" {
			names = append(names, setting.Name)
			index = append(index, i)
		}
	}

	return names, index
}

// GetStringFilterNames returns a list of filter names (flags).
func (c *Cfg) GetStringFilterNames() (names []string, index []int) {
	for i, setting := range c.settings {
		if setting.IsFlag && setting.Type == "string" {
			names = append(names, setting.Name)
			index = append(index, i)
		}
	}

	return names, index
}

// GetBoolFilterNames returns a list of filter names (flags).
func GetBoolFilterNames() ([]string, []int) {
	return configs[0].GetBoolFilterNames()
}

// GetIntFilterNames returns a list of filter names (flags).
func GetIntFilterNames() ([]string, []int) {
	return configs[0].GetIntFilterNames()}

// GetStringFilterNames returns a list of filter names (flags).
func GetStringFilterNames() ([]string, []int) {
	return configs[0].GetStringFilterNames()
}

func (c *Cfg) newFlagSet() {
	if c.flagSet != nil {
		return
	}

	c.flagSet = flag.NewFlagSet(c.name, flag.ContinueOnError)
}
