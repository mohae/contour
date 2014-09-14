package contour

import (
	"flag"
	"fmt"
)

// FilterArgs takes the passed args and filter's the flags out of them.
// The populated flags override their settings, according to the override
// rules. Successful overrides result in the relevant AppConfig setting
// being updated along with the environment variable.
//
// Any args left, after filtering, are returned to the caller.
// TODO: refactor this for greater abstraction (as long as it doesn't come
//	at the cost of reflection)
func (c *Cfg) FilterArgs(flagSet *flag.FlagSet, args []string) ([]string, error) {
	// Get the flag filters from the config variable information.
	boolFilterNames := c.GetBoolFilterNames()

	// Preallocate the worst case scenario.
	boolFilters := make([]*bool, len(boolFilterNames))
	bFilterNames := make([]string, len(boolFilterNames))
	var flags int

	for _, name := range boolFilterNames {
		if c.Settings[name].IsFlag {
			boolFilters[flags] = flagSet.Bool(name, c.Settings[name].Value.(bool), fmt.Sprintf("filter %s", name))
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
		if c.Settings[name].IsFlag {
			intFilters[flags] = flagSet.Int(name, c.Settings[name].Value.(int), fmt.Sprintf("filter %s", name))
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
		if c.Settings[name].IsFlag {
			fmt.Println(flags)
			stringFilters[flags] = flagSet.String(name, c.Settings[name].Value.(string), fmt.Sprintf("filter %s", name))
			sFilterNames[flags] = name
			flags++
		}
	}

	// Parse args for flags
	err := flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}

	// Get the remaining args
	cmdArgs := flagSet.Args()

	// Process the captured values
	for i, v := range boolFilters {
		if v != c.Settings[bFilterNames[i]].Value {
			Override(bFilterNames[i], v)
		}
	}

	for i, v := range intFilters {
		if v != c.Settings[iFilterNames[i]].Value {
			Override(iFilterNames[i], v)
		}
	}

	for i, v := range stringFilters {
		if v != c.Settings[sFilterNames[i]].Value {
			Override(sFilterNames[i], v)
		}
	}

	return cmdArgs, nil
}

func FilterArgs(flagSet *flag.FlagSet, args []string) ([]string, error) {
	return configs[app].FilterArgs(flagSet, args)
}
