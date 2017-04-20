package contour

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
)

var (
	ErrFlagsParsed   = errors.New("flags have already been parsed")
	ErrUseFlagsFalse = errors.New("not set to use flags")
)

// ParseFlags parses the args for the settings. Only settings of type Flag can
// be set via ParseFlags. Flags have the highest precedence. After parsing, any
// non-flag args are returned to the caller and a list of flags in the args is
// cached.
//
// If settings is not set to use flags, the args will be returned along with an
// ErrUseFlagsFalse. If settings has already parsed flags, the args are
// returned along with an ErrFlagsParsed.
//
// All of settings' flags must be registered prior to calling ParseFlags.
func (s *Settings) ParseFlags(args []string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.parseFlags(args)
}

// assume that the lock has already been obtained.
func (s *Settings) parseFlags(args []string) ([]string, error) {
	// nothing to do
	if !s.useFlags {
		return args, ErrUseFlagsFalse
	}

	if s.flagsParsed {
		return args, ErrFlagsParsed
	}
	// Get the flag information and set the flagSet
	s.setFlags()

	// Parse args for flags
	err := s.flagSet.Parse(args)
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

// setFlags goes through all the settings and sets the flagset vars for any
// that have IsFlag set to true. It a setting IsFlag but its type is
// interface{} it will not be added to the flagset as parsing interface{} is
// not supported.
func (s *Settings) setFlags() {
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
}

// Visited returns the names of all settings' flags that were set during flag
// parsing, in lexical order.
func (s *Settings) Visited() []string { return s.parsedFlags }

// WasVisited returns if settings' flag k was processed during flag parsing.
func (s *Settings) WasVisited(k string) bool {
	for i := range s.parsedFlags {
		if s.parsedFlags[i] == k {
			return true
		}
	}
	return false
}

// ParseFlags parses the command-line args from os.Args[1:]. Only settings of
// type Flag can be set via ParseFlags. Flags have the highest precedence.
// After parsing, any non-flag args are returned to the caller and a list of
// flags in the args is cached.
//
// If the standard settings has already parsed the flags or standard settings
// is set to not use flags, nothing will be done and nothing will be returned.
//
// All of standard settings' flags must be registered prior to calling.
func ParseFlags() ([]string, error) { return std.ParseFlags(os.Args[1:]) }

// Visited returns the names of all standard settings' flags that were set
// during flag parsing, in lexical order.
func Visited() []string { return std.Visited() }

// WasVisited returns if a standard settings flag k was parsed in the
// processing of args.
func WasVisited(k string) bool { return std.WasVisited(k) }
