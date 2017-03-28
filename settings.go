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

// Settings is a named group of settings and information related to that set,
// e.g. configuration file, if applicable, if it should emit an error when the
// configuration file is missing, flags parsed, etc. This is safe for
// concurrent use.
//
// The name of the Settings is used for environment variable naming, if
// applicable.
//
// There will always be a package global Settings whose name is the
// application's name.
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
	useEnvVars bool
	// If the settings have been updated from environment variables.
	envVarsSet bool
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

// New provides an initialized Settings named name.
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

// SetConfFilename sets the Settings' configuration filename and configures
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

// Set sets the registered settings according to Settings' configuration: it
// can be updated using a configuration file and/or environment variables; in
// that order of precedence. This is only run once; subsequent calls will
// result in no changes. Only settings that are of type ConfFileVar or EnvVar
// will be affected. This does not handle flags.
func Set() error { return settings.Set() }
func (s *Settings) Set() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// if this has already been set from env vars and config, don't do it again.
	// TODO: decide if this should be handled differently to allow for reload.
	if s.confFileVarsSet && s.envVarsSet {
		return nil
	}
	err := s.updateFromEnvVars()
	if err != nil {
		return fmt.Errorf("setting configuration from env failed: %s", err)
	}

	err = s.setFromConfFile()
	if err != nil {
		return fmt.Errorf("setting configuration from file failed: %s", err)
	}
	return nil
}

// SetFromEnvVars sets the settings that are of type EnvVar from env vars if
// the Settings is set to use env vars. If any settings are registered as an
// EnvVar settings, the use env vars flag will be set to true. This can be
// overridden.
//
// Once a Settings has been set from environment variables they will not be
// updated again on subsequent calls.
//
// A setting's env var name is a concatonation of the setting's name, an
// underscore, (_), and the Settings' name, e.g. a Settings with the name
// 'foo' and a setting whose name is 'bar' will result in 'FOO_BAR'.
func SetFromEnvVars() error { return settings.SetFromEnvVars() }

// SetFromEnvVars sets the settings that are of type EnvVars from env vars if
// the Settings is set to use env vars. If any settings are registered as an
// EnnVars settings, the use env vars flag will be set to true. This can be
// overridden.
//
// Once a Settings has been set from environment variables they will not be
// updated again on subsequent calls.
//
// A setting's env name is a concatonation of the setting's name, an underscore
// (_), and the Settings' name, e.g. a Settings with the name 'foo' and a
// setting whose name is 'bar' will result in 'FOO_BAR'.
func (s *Settings) SetFromEnvVars() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.updateFromEnvVars()
}

func (s *Settings) updateFromEnvVars() error {
	if !s.useEnvVars || s.envVarsSet {
		return nil
	}
	var err error
	for k, v := range s.settings {
		if !v.IsEnvVar {
			continue
		}
		tmp := os.Getenv(s.EnvVarName(k))
		if tmp != "" {
			switch v.Type {
			case _bool:
				b, _ := strconv.ParseBool(tmp)
				err = s.updateBool(EnvVar, k, b)
			case _int:
				i, err := strconv.Atoi(tmp)
				if err != nil {
					return fmt.Errorf("getenv %s: %s", s.EnvVarName(k), err)
				}
				err = s.updateInt(EnvVar, k, i)
			case _int64:
				i, err := strconv.ParseInt(tmp, 10, 64)
				if err != nil {
					return fmt.Errorf("getenv %s: %s", s.EnvVarName(k), err)
				}
				err = s.updateInt64(EnvVar, k, i)
			case _string:
				err = s.updateString(EnvVar, k, tmp)
			default:
				return fmt.Errorf("%s: unsupported env variable type: %s", s.EnvVarName(k), v.Type)
			}
			if err != nil {
				return fmt.Errorf("get env %s: %s", s.EnvVarName(k), err)
			}
			// lock to check next setting, if there is one.
		}
	}
	// Rlock isn't sufficient for updating to close it and get a Lock() for update.
	s.envVarsSet = true
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
// searched when looking for the configuration file.
func SearchPath() bool { return settings.SearchPath() }

// SearchPath returns whether or not the Path environment variable should be
// searched when looking for the configuration file.
func (s *Settings) SearchPath() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.searchPath
}

// SetSearchPath set's whether or not the user's PATH env variable should be
// searched for the configuratiom file.
func SetSearchPath(b bool) { settings.SetSearchPath(b) }

// SetSearchPath set's whether or not the user's PATH env variable should be
// searched for the configuratiom file.
func (s *Settings) SetSearchPath(b bool) {
	s.mu.Lock()
	s.searchPath = b
	s.mu.Unlock()
}

// UseConfFile returns if ConfFileVar settings are to be updated from a
// configuration file.
func UseConfFile() bool { return settings.UseConfFile() }

// UseConfFile returns if ConfFileVar settings are to be updated from a
// configuration file.
func (s *Settings) UseConfFile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useConfFile
}

// UseEnvVars returns whether or not environment variables are used.
func UseEnvVars() bool { return settings.useEnvVars }

// UseEnvVars returns whether or not environment variables are used.
func (s *Settings) UseEnvVars() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useEnvVars
}

// IsSet returns if the Settings has been set from all of its configured
// inputs, as applicable: env vars, configuration file, and flags.
func IsSet() bool { return settings.IsSet() }

// IsSet returns if the Settings has been set from all of its configured
// inputs, as applicable: env vars, configuration file, and flags.
func (s *Settings) IsSet() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.useConfFile && !s.confFileVarsSet {
		return false
	}
	if s.useEnvVars && !s.envVarsSet {
		return false
	}
	if s.useFlags && !s.flagsParsed {
		return false
	}
	// Everything has been updated (is set) according to Settings' configuration.
	return true
}

// SetUsage sets the Usage func.
func SetUsage(f func()) { settings.SetUsage(f) }

// SetUsage sets the Usage func.
func (s *Settings) SetUsage(f func()) {
	s.mu.Lock()
	s.flagSet.Usage = f
	s.mu.Unlock()
}

// Name returns the Settings' name.
func Name() string { return settings.Name() }

// Name returns the Settings' name.
func (s *Settings) Name() string {
	return s.name
}

// IsCoreE returns if setting k is a Core setting. If setting k doesn't exist,
// a SettingNotFoundErr will be returned.
func IsCoreE(name string) (bool, error) { return settings.IsCoreE(name) }

// IsCoreE returns if setting k is a Core setting. If setting k doesn't exist,
// a SettingNotFoundErr will be returned.
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

// IsCore returns if setting k is a Core setting. If setting k doesn't exist,
// a false will be returned.
func IsCore(name string) bool { return settings.IsCore(name) }

// IsCore returns if setting k is a Core setting. If setting k doesn't exist,
// a false will be returned.
func (s *Settings) IsCore(name string) bool {
	b, _ := s.IsCoreE(name)
	return b
}

// IsConfFileVarE returns if setting k is a ConfFileVar setting. If setting k
// doesn't exist, a SettingNotFoundErr will be returned.
func IsConfFileVarE(name string) (bool, error) { return settings.IsConfFileVarE(name) }

// IsConfFileVarE returns if setting k is a ConfFileVar setting. If setting k
// doesn't exist, a SettingNotFoundErr will be returned.
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

// IsConfFileVar returns if setting k is a ConfFileVar setting. If setting k
// doesn't exist, a false will be returned.
func IsConfFileVar(name string) bool { return settings.IsConfFileVar(name) }

// IsConfFileVar returns if setting k is a ConfFileVar setting. If setting k
// doesn't exist, a false will be returned.
func (s *Settings) IsConfFileVar(name string) bool {
	b, _ := s.IsConfFileVarE(name)
	return b
}

// IsEnvVarE returns if setting k is an EnvVar setting. If setting k doesn't
// exist, a SettingNotFoundErr will be returned.
func IsEnvVarE(name string) (bool, error) { return settings.IsEnvVarE(name) }

// IsEnvVarE returns if setting k is an EnvVar setting. If setting k doesn't
// exist, a SettingNotFoundErr will be returned.
func (s *Settings) IsEnvVarE(name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isEnvVar(name)
}

func (s *Settings) isEnvVar(name string) (bool, error) {
	val, ok := s.settings[name]
	if !ok {
		return false, SettingNotFoundErr{settingType: EnvVar, name: name}
	}
	return val.IsEnvVar, nil
}

// IsEnvVar returns if setting k is an EnvVar setting. If setting k doesn't
// exist, a false will be returned.
func IsEnvVar(name string) bool { return settings.IsEnvVar(name) }

// IsEnvVar returns if setting k is an EnvVar setting. If setting k doesn't
// exist, a false will be returned.
func (s *Settings) IsEnvVar(name string) bool {
	b, _ := s.IsEnvVarE(name)
	return b
}

// IsFlagE returns if setting k is a Flag setting. If setting k doesn't exist,
// a SettingNotFoundErr will be returned.
func IsFlagE(name string) (bool, error) { return settings.IsFlagE(name) }

// IsFlagE returns if setting k is a Flag setting. If setting k doesn't exist,
// a SettingNotFoundErr will be returned.
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

// IsFlag returns if setting k is a Flag setting. If setting k doesn't exist,
// a false will be returned.
func IsFlag(name string) bool { return settings.IsFlag(name) }

// IsFlag returns if setting k is a Flag setting. If setting k doesn't exist,
// a false will be returned.
func (s *Settings) IsFlag(name string) bool {
	b, _ := s.IsFlagE(name)
	return b
}

// EnvVarName returns the environment variable name for k. This will be
// NAME_K, where K is k and NAME is the name of Settings. For the pkg global
// Settings, this will be the executable name.
func EnvVarName(k string) string { return settings.EnvVarName(k) }

// EnvVarName returns the environment variable name for k. This will be
// NAME_K, where K is k and NAME is the name of Settings. For the pkg global
// Settings, this will be the executable name.
func (s *Settings) EnvVarName(k string) string {
	return strings.ToUpper(fmt.Sprintf("%s_%s", s.name, k))
}

// Exists returns if setting k exists.
func Exists(k string) bool { return settings.Exists(k) }

// Exists returns if setting k exists.
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
