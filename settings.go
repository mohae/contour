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
	errOnMissingFile bool
	// the key for the conf file setting.
	confFileKey string
	// search the path env var, in addition to pwd & executalbe dir, for cfc file.
	searchPath bool
	flagSet    *flag.FlagSet
	// file is the name of the config file
	file string
	// encoding is what encoding scheme is used for this config.
	encoding string
	// Settings contains a map of the configuration Settings for this
	// config.
	settings map[string]setting
	// Whether configuration settings have been registered and set.
	useCfg bool
	// if the settings have been loaded from config
	cfgSet bool
	// tracks the vars that are exposed to cfg
	cfgVars map[string]struct{}
	// useEnv: whether this config writes to and reads from environment
	// variables. If false, Settings are stored only in Config.
	useEnv bool
	envSet bool
	// Whether flags have been registered and set.
	useFlags    bool
	flagsParsed bool
	// interface contains pointer to a variable
	flagVars map[string]interface{}
	// maps short flags to the long version
	shortFlags map[string]string
}

// AppCfg returns the global cfg.
//
// Contour has a set of functions that implicitly interact with configs[app].
// If the application is only going to use one configuration, this is what
// should be used as one can just interact with contour, instead of directly
// with the app config, which is also supported.
func AppSettings() *Settings {
	return settings
}

// New provides an initialized Settings.
func New(name string) *Settings {
	return &Settings{
		name:             name,
		errOnMissingFile: true,
		searchPath:       true,
		flagSet:          flag.NewFlagSet(name, flag.ContinueOnError),
		settings:         map[string]setting{},
		cfgVars:          map[string]struct{}{},
		useCfg:           true,
		useEnv:           true,
		flagVars:         map[string]interface{}{},
		shortFlags:       map[string]string{},
	}
}

// Set updates the registered settings according to Settings' configuration:
// it can be updated using a configuration file and/or environment variables;
// in that order of precedence. This is only run once; subsequent calls will
// result in no changes.
//
// Only settings that are set as Environment, Cfg, or Flag types are updateable
// from environment variables.
//
// Only settings that are set as Cfg or Flag types are updateable from a
// configuration file.
func Set() error { return settings.Set() }
func (s *Settings) Set() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// if this has already been set from env vars and config, don't do it again.
	// TODO: decide if this should be handled differently to allow for reload.
	if s.cfgSet && s.envSet {
		return nil
	}
	err := s.updateFromEnv()
	if err != nil {
		return fmt.Errorf("setting configuration from env failed: %s", err)
	}

	err = s.setFromFile()
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
// A setting's env name is a concatonation of the cfg's name, an underscore
// (_), and the setting name, e.g. a Settings with the name 'foo' and a setting
// whose name is 'bar' will result in 'FOO_BAR'.
func SetFromEnv() error { return settings.SetFromEnv() }

// SetFromEnv sets the settings that are of type Env from env vars if the
// Settings is set to use env vars. If any settings are registered as env
// settings, the use env vars flag will be set to true. This can be overridden.
//
// Once a Settings has been set from environment variables they will not be
// updated again on subsequent calls.
//
// A setting's env name is a concatonation of the cfg's name, an underscore
// (_), and the setting name, e.g. a Settings with the name 'foo' and a setting
// whose name is 'bar' will result in 'FOO_BAR'.
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
				err = s.updateBoolE(k, b)
			case _int:
				i, err := strconv.Atoi(tmp)
				if err != nil {
					return fmt.Errorf("getenv %s: %s", s.GetEnvName(k), err)
				}
				err = s.updateIntE(k, i)
			case _int64:
				i, err := strconv.ParseInt(tmp, 10, 64)
				if err != nil {
					return fmt.Errorf("getenv %s: %s", s.GetEnvName(k), err)
				}
				err = s.updateInt64E(k, i)
			case _string:
				err = s.updateStringE(k, tmp)
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

// SetFromFile set's the Cfg and Flag settings from the information found in
// the configuration file if there is one. If Settings is not set to use a
// configuration file, if the configuration filename is not set, or if it has
// already been set, nothing is done and no error is returned.
func SetFromFile() error {
	return settings.SetFromFile()
}

// SetFromFile set's the Cfg and Flag settings from the information found in
// the configuration file if there is one. If Settings is not set to use a
// configuration file, if the configuration filename is not set, or if it has
// already been set, nothing is done and no error is returned.
func (s *Settings) SetFromFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.setFromFile()
}

// setFromFile set's Cfg and Flag settings from the information found in the
// configuration file, if there is one. This assumes the caller already holds
// the lock.
func (s *Settings) setFromFile() error {
	if !s.useCfg || s.cfgSet {
		return nil
	}
	setting, ok := s.settings[s.confFileKey]
	if !ok {
		// Wasn't configured, nothing to do. Not an error.
		return nil
	}
	n := setting.Value.(string)
	if n == "" {
		// This isn't an error as config file is allowed to not exist
		// TODO:
		//	Possible add a CfgFileRequired flag
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
	cfg, err := unmarshalCfgBytes(f, b)
	if err != nil {
		return fmt.Errorf("set from file: %s: %s", n, err)
	}

	// if nothing was returned and no error, nothing to do
	if cfg == nil {
		return nil
	}
	// Go through settings and update setting values.
	for k, v := range cfg.(map[string]interface{}) {
		// otherwise update the setting
		err = s.update(k, v)
		if err != nil {
			return fmt.Errorf("update setting: %s", err)
		}
	}
	return nil
}

// ErrOnMissingFile returns whether a missing config file should result in an
// error. This only applies when useCfg == true
func ErrOnMissingFile() bool { return settings.ErrOnMissingFile() }
func (s *Settings) ErrOnMissingFile() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errOnMissingFile
}

// SetErrOnMissingFile returns whether a missing config file should result in an
// error. This only applies when useFile == true
func SetErrOnMissingFile(b bool) { settings.SetErrOnMissingFile(b) }
func (s *Settings) SetErrOnMissingFile(b bool) {
	s.mu.Lock()
	s.errOnMissingFile = b
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
// searched for the cfg file.
func SetSearchPath(b bool) { settings.SetSearchPath(b) }
func (s *Settings) SetSearchPath(b bool) {
	s.mu.Lock()
	s.searchPath = b
	s.mu.Unlock()
}

// UseCfg returns whether this cfg uses an external, non env, cfg.
func UseCfg() bool { return settings.UseCfg() }
func (s *Settings) UseCfg() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.useCfg
}

// SetUseCfg set's whether an external, non-env, cfg should be used with this Cfg.
func SetUseCfg(b bool) { settings.SetUseCfg(b) }
func (s *Settings) SetUseCfg(b bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.useCfg = b
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

// ConfFileKey returns the value of confFileKey.
func ConfFileKey() string { return settings.ConfFileKey() }
func (s *Settings) ConfFileKey() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.confFileKey
}

// IsSet returns if the Settings has been set from all of its configured
// inputs, as applicable: env vars, configuration file, and flags.
func IsSet() bool { return settings.IsSet() }
func (s *Settings) IsSet() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.useCfg && !s.cfgSet {
		return false
	}
	if s.useEnv && !s.envSet {
		return false
	}
	if s.useFlags && !s.flagsParsed {
		return false
	}
	// Either post registration cfg isn't being used, or everything is set.
	return true
}

// SetUsage sets flagSet.Usage
func SetUsage(f func()) { settings.SetUsage(f) }
func (s *Settings) SetUsage(f func()) {
	s.mu.Lock()
	s.flagSet.Usage = f
	s.mu.Unlock()
}

// Name returns the cfg's name.
func Name() string { return settings.Name() }
func (s *Settings) Name() string {
	return s.name
}

// IsCore returns whether the passed setting is a core setting.
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

// IsCfg returns whether the passed setting is a cfg setting.
func IsCfgE(name string) (bool, error) { return settings.IsCfgE(name) }
func (s *Settings) IsCfgE(name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isCfg(name)
}

func (s *Settings) isCfg(name string) (bool, error) {
	val, ok := s.settings[name]
	if !ok {
		return false, SettingNotFoundErr{settingType: File, name: name}
	}
	return val.IsCfg, nil
}

func IsCfg(name string) bool { return settings.IsCfg(name) }
func (s *Settings) IsCfg(name string) bool {
	b, _ := s.IsCfgE(name)
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

// canUpdate checks to see if the passed setting key is updateable. If the key
// is not updateable, a false is returned along with an error. This assumes
// that the lock has already been obtained by the caller.
func (s *Settings) canUpdate(k string) (bool, error) {
	// See if the key exists, if it doesn't already exist, it can't be updated.
	v, ok := s.settings[k]
	if !ok {
		return false, SettingNotFoundErr{name: k}
	}
	// See if there are any settings that prevent it from being overridden.  Core and
	// environment variables are never settable. Core must be set during registration.
	if v.IsCore {
		return false, fmt.Errorf("%s: core settings cannot be updated", k)
	}
	if v.IsFlag && s.flagsParsed {
		return false, fmt.Errorf("%s: flag settings cannot be updated after parsing", k)
	}
	// Everything else is updateable.
	return true, nil
}

// canOverride() checks to see if the setting can be overridden. Overrides only
// come from flags. If it can't be overridden, it must be set via application,
// environment variable, or cfg file. This assumes the lock has already been
// obtained by the caller.
func (s *Settings) canOverride(k string) bool {
	// an empty key cannot Override
	if k == "" {
		return false
	}
	// See if the key exists, if it doesn't already exist, it can't be overridden
	v, ok := s.settings[k]
	if !ok {
		return false
	}
	// See if there are any settings that prevent it from being overridden.
	// Core can never be overridden-must be a flag to override.
	if v.IsCore {
		return false
	}
	// flags can only be set prior to arg filtering, after which you must use
	// Override().
	if v.IsFlag && s.flagsParsed {
		return false
	}
	return true
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
func unmarshalCfgBytes(f Format, buff []byte) (interface{}, error) {
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
