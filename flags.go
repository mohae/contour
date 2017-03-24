package contour

import (
	"flag"
	"fmt"
	"sort"
)

// ParseFlags parses the args for flags. Only Flag settings can be set via
// flags. For flag settings, flags have the highest precedence.
//
// If the flags have already been parsed or Settings is set to not use flags,
// nothing will be done and nothing will be returned.
//
// If this is called, the Settings should already have all of its settings
// registered, the config file loaded (if applicable), and the envorinment
// variables loaded (if applicable). In most cases, Settings.Set() should be
// used instead; after registering all of the settings and configuring
// Settings' behavior.
//
// Any args left, after filtering, are returned to the caller.
func ParseFlags(args []string) ([]string, error) { return settings.ParseFlags(args) }
func (s *Settings) ParseFlags(args []string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.parseFlags(args)
}

// assume that the lock has already been obtained.
func (s *Settings) parseFlags(args []string) ([]string, error) {
	// nothing to do
	if !s.useFlags || s.flagsParsed {
		return nil, nil
	}

	// Get the flag information and set the flagSet
	err := s.setFlags()
	if err != nil {
		return nil, err
	}
	// Parse args for flags
	err = s.flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}

	// get the visited flags
	var visited []*flag.Flag
	visitor := func(a *flag.Flag) {
		visited = append(visited, a)
	}

	s.flagSet.Visit(visitor)
	// Update settings with the updated flag values
	for _, f := range visited {
		v, ok := s.settings[f.Name]
		if !ok {
			// see if it's a short flag
			v, ok = s.settings[s.shortFlags[f.Name]]
			if !ok {
				continue
			}
		}
		v.Value = f.Value
		s.settings[v.Name] = v
		s.parsedFlags = append(s.parsedFlags, v.Name)
	}

	// sort the parsed flagsParsed
	sort.Strings(s.parsedFlags)

	// Get the remaining args
	cmdArgs := s.flagSet.Args()

	s.flagsParsed = true

	return cmdArgs, nil
}

// setFlags set's up the flagFilter information while setting the Setting's
// flagset. This assumes that the lock has been obtained by the caller.
func (s *Settings) setFlags() error {
	// Get the flag filters from the config variable information.
	for _, v := range s.settings {
		if v.IsFlag {
			switch v.Type {
			case _bool:
				s.flagVars[v.Name] = s.flagSet.Bool(v.Name, v.Value.(bool), v.Usage)
				if v.Short != "" {
					s.flagSet.BoolVar(s.flagVars[v.Name].(*bool), v.Short, v.Value.(bool), v.Usage)
				}
			case _int:
				s.flagVars[v.Name] = s.flagSet.Int(v.Name, v.Value.(int), v.Usage)
				if v.Short != "" {
					s.flagSet.IntVar(s.flagVars[v.Name].(*int), v.Short, v.Value.(int), v.Usage)
				}
			case _int64:
				s.flagVars[v.Name] = s.flagSet.Int64(v.Name, v.Value.(int64), v.Usage)
				if v.Short != "" {
					s.flagSet.Int64Var(s.flagVars[v.Name].(*int64), v.Short, v.Value.(int64), v.Usage)
				}
			case _string:
				s.flagVars[v.Name] = s.flagSet.String(v.Name, v.Value.(string), v.Usage)
				if v.Short != "" {
					s.flagSet.StringVar(s.flagVars[v.Name].(*string), v.Short, v.Value.(string), v.Usage)
				}
			}
		}
	}
	return nil
}
