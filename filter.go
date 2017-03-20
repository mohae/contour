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
func FilterArgs(args []string) ([]string, error) { return settings.FilterArgs(args) }
func (s *Settings) FilterArgs(args []string) ([]string, error) {
	// Get the ArgFilter (flag) information and set the flagSet
	err := s.setCfgFlags()
	if err != nil {
		return nil, err
	}
	// the arg slice may be modified, any --flags get normalized to -flag
	args, flags := s.getPassedFlags(args)
	if flags == nil {
		return args, nil
	}
	s.mu.Lock()
	// Parse args for flags
	err = s.flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}
	// Get the remaining args
	cmdArgs := s.flagSet.Args()
	s.mu.Unlock()
	// Process the captured values
	for _, n := range flags {
		s.mu.RLock()
		ptr, ok := s.filterVars[n]
		if !ok {
			continue
		}
		s.mu.RUnlock()
		s.Override(n, ptr)
	}
	s.mu.Lock()
	s.argsFiltered = true
	// nil these out as they are no longer needed
	s.boolFilterNames = nil
	s.intFilterNames = nil
	s.int64FilterNames = nil
	s.stringFilterNames = nil
	s.filterVars = nil
	s.shortFlags = nil
	s.mu.Unlock()
	return cmdArgs, nil
}

// setCfgFlags set's up the argFilter information while setting the Cfg's
// flagset.
func (s *Settings) setCfgFlags() error {
	// Get the flag filters from the config variable information.
	s.mu.Lock()
	for _, v := range s.settings {
		if v.IsFlag {
			switch v.Type {
			case "bool":
				s.filterVars[v.Name] = s.flagSet.Bool(v.Name, v.Value.(bool), v.Usage)
				s.boolFilterNames = append(s.boolFilterNames, v.Name)
			case "int":
				s.filterVars[v.Name] = s.flagSet.Int(v.Name, v.Value.(int), v.Usage)
				s.intFilterNames = append(s.intFilterNames, v.Name)
			case "int64":
				s.filterVars[v.Name] = s.flagSet.Int64(v.Name, v.Value.(int64), v.Usage)
				s.int64FilterNames = append(s.int64FilterNames, v.Name)
			case "string":
				s.filterVars[v.Name] = s.flagSet.String(v.Name, v.Value.(string), v.Usage)
				s.stringFilterNames = append(s.stringFilterNames, v.Name)
			}
		}
	}
	s.mu.Unlock()
	return nil
}

// getPassedFlags returns a slice of flags that exist in the arg list. The
// original args are not affected,
func (s *Settings) getPassedFlags(args []string) ([]string, []string) {
	// keeps track of what flags
	var flags []string
	s.mu.RLock()
	defer s.mu.RUnlock()
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
			n, ok := s.shortFlags[arg]
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
			val, ok := s.settings[arg]
			if !ok {
				continue
			}
			if val.IsFlag {
				flags = append(flags, arg)
			}
		}
	}
	return args, flags
}
