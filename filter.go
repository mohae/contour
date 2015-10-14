package contour

import (
	"fmt"
	"strings"
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
	// the arg slice may be modified, any --flags get normalized to -flag
	args, flags := c.getPassedFlags(args)
	if flags == nil {
		return args, nil
	}
	c.RWMutex.Lock()
	// Parse args for flags
	err = c.flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}
	// Get the remaining args
	cmdArgs := c.flagSet.Args()
	c.RWMutex.Unlock()
	// Process the captured values
	for _, n := range flags {
		c.RWMutex.RLock()
		ptr, ok := c.filterVars[n]
		if !ok {
			continue
		}
		c.RWMutex.RUnlock()
		c.Override(n, ptr)
	}
	c.RWMutex.Lock()
	c.argsFiltered = true
	// nil these out as they are no longer needed
	c.boolFilterNames = nil
	c.intFilterNames = nil
	c.int64FilterNames = nil
	c.stringFilterNames = nil
	c.filterVars = nil
	c.shortFlags = nil
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
			case "int":
				c.filterVars[s.Name] = c.flagSet.Int(s.Name, s.Value.(int), s.Usage)
				c.intFilterNames = append(c.intFilterNames, s.Name)
			case "int64":
				c.filterVars[s.Name] = c.flagSet.Int64(s.Name, s.Value.(int64), s.Usage)
				c.int64FilterNames = append(c.int64FilterNames, s.Name)
			case "string":
				c.filterVars[s.Name] = c.flagSet.String(s.Name, s.Value.(string), s.Usage)
				c.stringFilterNames = append(c.stringFilterNames, s.Name)
			}
		}
	}
	c.RWMutex.Unlock()
	return nil
}

// getPassedFlags returns a slice of flags that exist in the arg list. The
// original args are not affected,
func (c *Cfg) getPassedFlags(args []string) ([]string, []string) {
	// keeps track of what flags
	var flags []string
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	for i, arg := range args {
		if strings.HasPrefix(arg, "--") {
			arg = arg[1:len(arg)]
			args[i] = arg
		}
		if strings.HasPrefix(arg, "-") {
			arg = strings.TrimPrefix(arg, "-")
			// if the key can be split by =, it's the first element
			split := strings.Split(arg, "=")
			arg = split[0]
			// Check to see if this is a short flag, if so, use the full flag name.
			n, ok := c.shortFlags[arg]
			if ok {
				arg = n
				// replace the short version with the arg.
				if len(split) == 1 {
					args[i] = "-" + arg
				} else {
					split[0] = "-" + arg
					args[i] = strings.Join(split, "=")
				}
			}
			// Check to see if it already exists and is a flag before adding it
			s, ok := c.settings[arg]
			if !ok {
				continue
			}
			if s.IsFlag {
				flags = append(flags, arg)
			}
		}
	}
	return args, flags
}
