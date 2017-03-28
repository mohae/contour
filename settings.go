package contour

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/mohae/cjson"
	"gopkg.in/yaml.v2"
)

// Settings is a group of settings and holds all of the application setting
// information. Even though contour automatically uses environment variables,
// unless its told to ignore them, it still needs to maintain state
// information about each setting so it knows how to handle attempst to update.
// TODO:
//	* support ignoring environment variables
//
type Settings struct {
	name string
	mu   sync.RWMutex
	// if an attempt to load configuration from a file should error if the file
	// does not exist.
	errOnMissingConfFile bool
	// the name of the configuration filename variable: defaults to the
	// ConfFilenameSettingVarName constant.
	confFilenameVarName string
	// search the path env var, in addition to wd & executalbe dir, for the conf
	// file. The path env var is assumed to be NAMEPATH where name is
	// Settings.name
	searchPath bool
	// file is the name of the configuration file
	confFilename string
	// Encoding is what encoding scheme is used for this config.
	encoding string
	// Tracks the vars that are exposed to the configuration file. Only vars in
	// this map can be set from a configuration file.
	confFileVars map[string]struct{}
	// If settings should be updated from a configuration file.
	useConfFile bool
	// If the settings have been updated from configuration file.
	confFileVarsSet bool
	// If settings should be loaded from environment variables. Environment
	// variable names will be in the form of NAME_VARNAME where NAME is this
	// Setting's name.
	useEnv bool
	// If the settings have been updated from environment variables.
	envSet bool
	// flagset is the set of flags for arg parsing.
	flagSet *flag.FlagSet
	// If the settings should be updated from passed args.
	useFlags bool
	// If the settings have been updated from the passed args.
	flagsParsed bool
	// The map of variables that capture flag information.
	flagVars map[string]interface{}
	// Maps short flags to the long version
	shortFlags map[string]string
	// parsedFlags are flags that were passed and parsed. Short flags are
	// normalized to the flag name.
	parsedFlags []string
	// Settings contains a map of all the configuration settings for this
	// Setting and each setting's information, including current Value.
	settings map[string]setting
}

// New provides an initialized Settings.
func New(name string) *Settings {
	return &Settings{
		name:                 name,
		errOnMissingConfFile: true,
		searchPath:           true,
		flagSet:              flag.NewFlagSet(name, flag.ContinueOnError),
		confFileVars:         map[string]struct{}{},
		flagVars:             map[string]interface{}{},
		shortFlags:           map[string]string{},
		settings:             map[string]setting{},
	}
}

// SetConfFilename sets the Setting's configuration filename and configures
// settings to use a configuration file. If the filename is empty, an error is
// returned.
func SetConfFilename(v string) error { return settings.SetConfFilename(v) }
func (s *Settings) SetConfFilename(v string) error {
	if v == "" {
		return fmt.Errorf("set configuration filename failed: no name provided")
	}

	// store the key value being used as the configuration setting name by caller
	s.mu.Lock()
	defer s.mu.Unlock()
	s.confFilename = v
	s.useConfFile = true
	return nil
}

// Set updates the registered settings according to Settings' configuration:
// it can be updated using a configuration file and/or environment variables;
// in that order of precedence. This is only run once; subsequent calls will
// result in no changes.
//
// Only settings that are set as Environment, Conf, or Flag types are
// updateable from environment variables.
//
// Only settings that are set as Conf or Flag types are updateable from a
// configuration file.
func Set() error { return settings.Set() }
func (s *Settings) Set() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// if this has already been set from env vars and config, don't do it again.
	// TODO: decide if this should be handled differently to allow for reload.
	if s.confFileVarsSet && s.envSet {
		return nil
	}
	err := s.updateFromEnv()
	if err != nil {
		return fmt.Errorf("setting configuration from env failed: %s", err)
	}

	err = s.setFromConfFile()
	if err != nil {
		return fmt.Errorf("setting configuration from file failed: %s", err)
	}
	return nil
}

// SetFromEnv sets the settings that are of type Env from env vars if the
// Settings is set to use env vars. If any settings are registered as env
// settings, the use env vars flag will be set to true. This can be overridden.
//
// Once a Settings has been set from environment variables they will not be
// updated again on subsequent calls.
//
// A setting's env name is a concatonation of the setting's name, an underscore
// (_), and the Settings' name, e.g. a Settings with the name 'foo' and a
// setting whose name is 'bar' will result in 'FOO_BAR'.
func SetFromEnv() error { return settings.SetFromEnv() }

// SetFromEnv sets the settings that are of type Env from env vars if the
// Settings is set to use env vars. If any settings are registered as env
// settings, the use env vars flag will be set to true. This can be overridden.
//
// Once a Settings has been set from environment variables they will not be
// updated again on subsequent calls.
//
// A setting's env name is a concatonation of the setting's name, an underscore
// (_), and the Settings' name, e.g. a Settings with the name 'foo' and a
// setting whose name is 'bar' will result in 'FOO_BAR'.
func (s *Settings) SetFromEnv() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateFromEnv()
}

func (s *Settings) updateFromEnv() error {
	if !s.useEnv || s.envSet {
		return nil
	}
	var err error
	for k, v := range s.settings {
		if !v.IsEnv {
			continue
		}
		tmp := os.Getenv(s.GetEnvName(k))
		if tmp != "" {
			switch v.Type {
			case _bool:
				b, _ := strconv.ParseBool(tmp)
				err = s.updateBool(Env, k, b)
			case _int:
				i, err := strconv.Atoi(tmp)
				if err != nil {
					return fmt.Errorf("getenv %s: %s", s.GetEnvName(k), err)
				}
				err = s.updateInt(Env, k, i)
			case _int64:
				i, err := strconv.ParseInt(tmp, 10, 64)
				if err != nil {
					return fmt.Errorf("getenv %s: %s", s.GetEnvName(k), err)
				}
				err = s.updateInt64(Env, k, i)
			case _string:
				err = s.updateString(Env, k, tmp)
			default:
				return fmt.Errorf("%s: unsupported env variable type: %s", s.GetEnvName(k), v.Type)
			}
			if err != nil {
				return fmt.Errorf("get env %s: %s", s.GetEnvName(k), err)
			}
			// lock to check next setting, if there is one.
		}
	}
	// Rlock isn't sufficient for updating to close it and get a Lock() for update.
	s.envSet = true
	return nil
}

// SetFromConfFile set's the Conf, Env, and Flag settings from the information
// found in the configuration file if there is one. If Settings is not set to
// use a configuration file, if the configuration filename is not set, or if it
// has already been set, nothing is done and no error is returned.
func SetFromConfFile() error {
	return settings.SetFromConfFile()
}

// SetFromConfFile set's the Conf, Env, and Flag settings from the information
// found in the configuration file if there is one. If Settings is not set to
// use a configuration file, if the configuration filename is not set, or if it
// has already been set, nothing is done and no error is returned.
func (s *Settings) SetFromConfFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.setFromConfFile()
}

// setFromFile set's Conf, Env, and Flag settings from the information found in
// the configuration file, if there is one. This assumes the caller already
// holds the lock.
func (s *Settings) setFromConfFile() error {
	if !s.useConfFile || s.confFileVarsSet {
		return nil
	}
	setting, ok := s.settings[s.confFilenameVarName]
	if !ok {
		// Wasn't configured, nothing to do. Not an error.
		return nil
	}
	n := setting.Value.(string)
	if n == "" {
		// This isn't an error as config file is allowed to not exist
		// TODO:
		//	Possible add a confFileRequired flag
		return nil
	}
	// get the file's format from the extension
	f, err := formatFromFilename(n)
	if err != nil {
		return fmt.Errorf("set from file: %s", err)
	}

	b, err := ioutil.ReadFile(n)
	if err != nil {
		return fmt.Errorf("set from file: %s", err)
	}
	cnf, err := unmarshalConfBytes(f, b)
	if err != nil {
		return fmt.Errorf("set from file: %s: %s", n, err)
	}

	// if nothing was returned and no error, nothing to do
	if cnf == nil {
		return nil
	}
	// Go through settings and update setting values.
	for k, v := range cnf.(map[string]interface{}) {
		// otherwise update the setting
		err = s.update(ConfFileVar, k, v)
		if err != nil {
			return fmt.Errorf("update setting: %s", err)
		}
	}
	return nil
}

// ErrOnMissingConfFile returns whether a missing config file should result in
// an error. This only applies when useConf == true
func ErrOnMissingConfFile() bool { return settings.ErrOnMissingConfFile() }
func (s *Settings) ErrOnMissingConfFile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errOnMissingConfFile
}

// SetErrOnMissingConfFile returns whether a missing config file should result
// in an error. This only applies when useFile == true
func SetErrOnMissingConfFile(b bool) { settings.SetErrOnMissingConfFile(b) }
func (s *Settings) SetErrOnMissingConfFile(b bool) {
	s.mu.Lock()
	s.errOnMissingConfFile = b
	s.mu.Unlock()
}

// SearchPath returns whether or not the Path environment variable should be
// searched when looking for the config file.
func SearchPath() bool { return settings.SearchPath() }
func (s *Settings) SearchPath() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.searchPath
}

// SetSearchPath set's whether or not the user's PATH env variable should be
// searched for the configuratiom file.
func SetSearchPath(b bool) { settings.SetSearchPath(b) }
func (s *Settings) SetSearchPath(b bool) {
	s.mu.Lock()
	s.searchPath = b
	s.mu.Unlock()
}

// UseConfFile returns if Conf settings are to be updated from a configuration
// file.
func UseConfFile() bool { return settings.UseConfFile() }
func (s *Settings) UseConfFile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useConfFile
}

// SetUseConfFile sets if Conf settings should be updated from a configuration
// file.
func SetUseConfFile(b bool) { settings.SetUseConfFile(b) }
func (s *Settings) SetUseConfFile(b bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.useConfFile = b
}

// UseEnv is whether or not environment variables are used.
func UseEnv() bool { return settings.useEnv }
func (s *Settings) UseEnv() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useEnv
}

// SetUseEnv set's whether or not environment variables should be used with
// this cfg.
func SetUseEnv(b bool) { settings.SetUseEnv(b) }
func (s *Settings) SetUseEnv(b bool) {
	s.mu.Lock()
	s.useEnv = b
	s.mu.Unlock()
}

// IsSet returns if the Settings has been set from all of its configured
// inputs, as applicable: env vars, configuration file, and flags.
func IsSet() bool { return settings.IsSet() }
func (s *Settings) IsSet() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.useConfFile && !s.confFileVarsSet {
		return false
	}
	if s.useEnv && !s.envSet {
		return false
	}
	if s.useFlags && !s.flagsParsed {
		return false
	}
	// Everything has been updated (is set) according to Settings' configuration.
	return true
}

// SetUsage sets flagSet.Usage
func SetUsage(f func()) { settings.SetUsage(f) }
func (s *Settings) SetUsage(f func()) {
	s.mu.Lock()
	s.flagSet.Usage = f
	s.mu.Unlock()
}

// Name returns the Settings' name.
func Name() string { return settings.Name() }
func (s *Settings) Name() string {
	return s.name
}

// IsCore returns whether the passed setting  is a core setting.
func IsCoreE(name string) (bool, error) { return settings.IsCoreE(name) }
func (s *Settings) IsCoreE(name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isCore(name)
}

func (s *Settings) isCore(name string) (bool, error) {
	val, ok := s.settings[name]
	if !ok {
		return false, SettingNotFoundErr{settingType: Core, name: name}
	}
	return val.IsCore, nil
}

func IsCore(name string) bool { return settings.IsCore(name) }
func (s *Settings) IsCore(name string) bool {
	b, _ := s.IsCoreE(name)
	return b
}

// IsConfE returns if the setting is a Conf setting. If a setting with the
// requested name does not exist, a SettingNotFoundErr will be returned.
func IsConfFileVarE(name string) (bool, error) { return settings.IsConfFileVarE(name) }
func (s *Settings) IsConfFileVarE(name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isConfFileVar(name)
}

func (s *Settings) isConfFileVar(name string) (bool, error) {
	val, ok := s.settings[name]
	if !ok {
		return false, SettingNotFoundErr{settingType: ConfFileVar, name: name}
	}
	return val.IsConfFileVar, nil
}

// IsConfFileVar returns if the setting is a Conf setting. If a setting with
// the requested name does not exist, a false will also be returned.
func IsConfFileVar(name string) bool { return settings.IsConfFileVar(name) }
func (s *Settings) IsConfFileVar(name string) bool {
	b, _ := s.IsConfFileVarE(name)
	return b
}

// IsEnv returns whether the passed setting is a env setting.
func IsEnvE(name string) (bool, error) { return settings.IsEnvE(name) }
func (s *Settings) IsEnvE(name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isEnv(name)
}

func (s *Settings) isEnv(name string) (bool, error) {
	val, ok := s.settings[name]
	if !ok {
		return false, SettingNotFoundErr{settingType: Env, name: name}
	}
	return val.IsEnv, nil
}

func IsEnv(name string) bool { return settings.IsEnv(name) }
func (s *Settings) IsEnv(name string) bool {
	b, _ := s.IsEnvE(name)
	return b
}

// IsFlag returns whether the passed setting is a flag setting.
func IsFlagE(name string) (bool, error) { return settings.IsFlagE(name) }
func (s *Settings) IsFlagE(name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isFlag(name)
}

func (s *Settings) isFlag(name string) (bool, error) {
	val, ok := s.settings[name]
	if !ok {
		return false, SettingNotFoundErr{settingType: Flag, name: name}
	}
	return val.IsFlag, nil
}

func IsFlag(name string) bool { return settings.IsFlag(name) }
func (s *Settings) IsFlag(name string) bool {
	b, _ := s.IsFlagE(name)
	return b
}

// GetEnvName returns the env variable name version of the passed string.
func GetEnvName(s string) string { return settings.GetEnvName(s) }
func (s *Settings) GetEnvName(v string) string {
	return strings.ToUpper(fmt.Sprintf("%s_%s", s.name, v))
}

// Exists returns if a setting with the key exists.
func Exists(k string) bool { return settings.Exists(k) }
func (s *Settings) Exists(k string) bool {
	s.mu.RLock()
	s.mu.RUnlock()
	return s.exists(k)
}

func (s *Settings) exists(k string) bool {
	_, err := s.get(k)
	if err == nil {
		return true
	}
	return false
}

// Visited returns the names of all flags that were set during argument
// parsing in lexical order.
func Visited() []string { return settings.Visited() }

// Visited returns the names of all flags that were set during argument
// parsing in lexical order.
func (s *Settings) Visited() []string { return s.parsedFlags }

// WasVisited returns if a flag was parsed in the processing of args.
func WasVisited(k string) bool { return settings.WasVisited(k) }

// WasVisited returns if a flag was parsed in the processing of args.
func (s *Settings) WasVisited(k string) bool {
	for i := range s.parsedFlags {
		if s.parsedFlags[i] == k {
			return true
		}
	}
	return false
}

// formatFromFilename gets the format from the passed filename.  An error will
// be returned if either the format isn't supported or the extension doesn't
// exist.  If the passed string has multiple dots, the last dot is assumed to
// be the extension.
func formatFromFilename(s string) (Format, error) {
	if s == "" {
		return Unsupported, fmt.Errorf("no configuration filename")
	}
	parts := strings.Split(s, ".")
	format := ""
	// case 0 has already been evaluated
	switch len(parts) {
	case 1:
		return Unsupported, fmt.Errorf("unable to determine %s's format: no extension", strings.TrimSpace(s))
	case 2:
		format = parts[1]
	default:
		// assume its the last part
		format = parts[len(parts)-1]
	}

	return ParseFormat(format)
}

// unmarshalConfBytes accepts bytes and unmarshals them using the correct
// format. Either the unmarshaled data or an error is returned.
//
// Supported formats:
//   json
//   toml
//   yaml
func unmarshalConfBytes(f Format, buff []byte) (interface{}, error) {
	var ret interface{}
	switch f {
	case JSON:
		err := cjson.Unmarshal(buff, &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	case TOML:
		_, err := toml.Decode(string(buff), &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	case YAML:
		err := yaml.Unmarshal(buff, &ret)
		if err != nil {
			return nil, err
		}
		return ret, nil
	default:
		return nil, UnsupportedFormatErr{f.String()}
	}
}
