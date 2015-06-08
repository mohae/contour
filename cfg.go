package contour

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

// Cfg is a group of Settings and holds all of the application setting
// information. Even though contour automatically uses environment variables,
// unless its told to ignore them, it still needs to maintain state
// information about each setting so it knows how to handle attempst to update.
// TODO:
//	* support ignoring environment variables
//
type Cfg struct {
	name string
	sync.RWMutex
	errOnMissingCfg bool
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
	cfgSet bool
	// tracks the vars that are exposed to cfg
	cfgVars map[string]struct{}
	// useEnv: whether this config writes to and reads from environment
	// variables. If false, Settings are stored only in Config.
	useEnv bool
	envSet bool
	// Whether flags have been registered and set.
	useFlags     bool
	argsFiltered bool
	// interface contains pointer to a variable
	filterVars map[string]interface{}
	// flag filters by type
	boolFilters   []string
	intFilters    []string
	int64Filters  []string
	stringFilters []string
	// maps short flags to the long version
	shortFlags map[string]string
}

// AppCfg returns the global cfg.
//
// Contour has a set of functions that implicitly interact with configs[app].
// If the application is only going to use one configuration, this is what
// should be used as one can just interact with contour, instead of directly
// with the app config, which is also supported.
func AppCfg() *Cfg {
	return appCfg
}

// NewConfig returns a *Cfg to the caller
func NewCfg(name string) *Cfg {
	return &Cfg{name: name, errOnMissingCfg: true, searchPath: true, flagSet: flag.NewFlagSet(name, flag.ContinueOnError), settings: map[string]setting{}, cfgVars: map[string]struct{}{}, useEnv: true, shortFlags: map[string]string{}}
}

// UpdateFromEnv updates the cfg settings from env vars: only when the Cfg's
// useEnv flag is set to True.  Cfg settings whose IsEnv flag is set to true
// will be processed. By default, any setting that is registered as a Cfg or
// Flag setting has their IsEnv value set to true. This can be changed.
//
// A setting's env name is a concatonation of the cfg's name, an underscore
// (_), and the setting name, e.g. a Cfg with the name 'rancher' and a setting
// whose name is 'log' will result in 'rancher_log'.
func UpdateFromEnv() error { return appCfg.UpdateFromEnv() }
func (c *Cfg) UpdateFromEnv() error {
	c.RWMutex.RLock()
	if !c.useEnv {
		c.RWMutex.RUnlock()
		return nil
	}
	name := c.name // cache it so I don't have to worry about the lock later
	var err error
	for k, v := range c.settings {
		if !v.IsEnv {
			continue
		}
		tmp := os.Getenv(fmt.Sprintf("%s_%s", name, k))
		if tmp != "" {
			c.RWMutex.RUnlock()
			switch v.Type {
			case "bool":
				b, _ := strconv.ParseBool(tmp)
				err = c.UpdateBoolE(k, b)
			case "int":
				i, err := strconv.Atoi(tmp)
				if err != nil {
					return fmt.Errorf("Loadenv error while parsing %s: %s", fmt.Sprintf("%s_%s", name, k), err)
				}
				err = c.UpdateIntE(k, i)
			case "int64":
				i, err := strconv.ParseInt(tmp, 10, 64)
				if err != nil {
					return fmt.Errorf("Loadenv error while parsing %s: %s", fmt.Sprintf("%s_%s", name, k), err)
				}
				err = c.UpdateInt64E(k, i)
			case "string":
				err = c.UpdateStringE(k, tmp)
			default:
				return fmt.Errorf("%s has an unsupported env variable type: %s", k, v.Type)
			}
			if err != nil {
				return fmt.Errorf("Loadenv error while setting %s: %s", fmt.Sprintf("%s_%s", name, k), err)
			}
			// lock to check next setting, if there is one.
			c.RWMutex.RLock()
		}
	}
	// Rlock isn't sufficient for updating to close it and get a Lock() for update.
	c.RWMutex.RUnlock()
	c.RWMutex.Lock()
	c.envSet = true
	c.RWMutex.Unlock()
	return nil
}

// ErrOnMissingCfg returns whether a missing config file should result in an
// error. This only applies when useCfg == true
func ErrOnMissingCfg() bool { return appCfg.ErrOnMissingCfg() }
func (c *Cfg) ErrOnMissingCfg() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.errOnMissingCfg
}

// SetErrOnMissingCfg returns whether a missing config file should result in an
// error. This only applies when useCfg == true
func SetErrOnMissingCfg(b bool) { appCfg.SetErrOnMissingCfg(b) }
func (c *Cfg) SetErrOnMissingCfg(b bool) {
	c.RWMutex.Lock()
	c.errOnMissingCfg = b
	c.RWMutex.Unlock()
}

// SearchPath returns whether or not the Path environment variable should be
// searched when looking for the config file.
func SearchPath() bool { return appCfg.SearchPath() }
func (c *Cfg) SearchPath() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.searchPath
}

// SetSearchPath set's whether or not the user's PATH env variable should be
// searched for the cfg file.
func SetSearchPath(b bool) { appCfg.SetSearchPath(b) }
func (c *Cfg) SetSearchPath(b bool) {
	c.RWMutex.Lock()
	c.searchPath = b
	c.RWMutex.Unlock()
}

// UseCfgFile returns whether this cfg uses an external, non env, cfg.
func UseCfg() bool { return appCfg.UseCfg() }
func (c *Cfg) UseCfg() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.useCfg
}

// SetUseCfg set's whether an external, non-env, cfg should be used with this Cfg.
func SetUseCfg(b bool) { appCfg.SetUseCfg(b) }
func (c *Cfg) SetUseCfg(b bool) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	c.useCfg = b
}

// UseEnv is whether or not environment variables are used.
func UseEnv() bool { return appCfg.useEnv }
func (c *Cfg) UseEnv() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.useEnv
}

// SetUseEnv set's whether or not environment variables should be used with
// this cfg.
func SetUseEnv(b bool) { appCfg.SetUseEnv(b) }
func (c *Cfg) SetUseEnv(b bool) {
	c.RWMutex.Lock()
	c.useEnv = b
	c.RWMutex.Unlock()
}

// SetCfg goes through the initialized Settings and updates the updateable
// settings, if a new, valid value is found.  This applies to, in order: Env
// variables and config files. For any that are not found, or that are
// immutable, once set, the original initialization values are used.
//
// Updates to the application defaults will be applied as follows:
//    * if useCfg, the values found within the cfgFile will be applied.
//    * if useEnv, the values found in the env vars will be applied.
//
// Up through Flags, and with the exception of setting the cfg file, the order
// of precedence are:
//     command-line flags
//     environment variables
//     cfg file settings
//     application defaults
func SetCfg() error { return appCfg.SetCfg() }
func (c *Cfg) SetCfg() error {
	c.RWMutex.RLock()
	useCfg := c.useCfg
	useEnv := c.useEnv
	c.RWMutex.RUnlock()
	fname, err := c.GetStringE(CfgFile)
	if err != nil {
		return fmt.Errorf("SetCfg failed: %s", err.Error())
	}
	if useCfg {
		// Load the Cfg
		buff, err := getFileBytes(fname)
		if err != nil {
			return fmt.Errorf("reading %s failed: %s", fname, err.Error())
		}
		err = c.UpdateFromCfg(buff)
		if err != nil {
			return err
		}
	}
	if useEnv {
		err = c.UpdateFromEnv()
		if err != nil {
			return fmt.Errorf("setting cfg from env failed: %s", err.Error())
		}
	}
	return nil
}

// UpdateFromCfg updates the application's default values with the setting
// values found in the cfg. Only Cfg and Flag settings are updated.
func (c *Cfg) UpdateFromCfg(buff []byte) error {
	cfgSettings, err := c.processCfg(buff)
	if err != nil {
		return fmt.Errorf("UpdateFromCfg error processing cfg data: %s", err.Error())
	}
	// if nothing was returned and no error, nothing to do
	if cfgSettings == nil {
		return nil
	}
	// Go through settings and update setting values.
	for k, v := range cfgSettings.(map[string]interface{}) {
		// otherwise update the setting
		err = c.updateE(k, v)
		if err != nil {
			return fmt.Errorf("UpdateFromCfg error updating setting: %s", err.Error())
		}
	}
	return nil
}

// CfgsProcessed determines whether, or not, all of the cfg sources have been
// processed for a given Cfg.
func CfgProcessed() bool { return appCfg.CfgProcessed() }
func (c *Cfg) CfgProcessed() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	if c.useCfg && !c.cfgSet {
		return false
	}
	if c.useEnv && !c.envSet {
		return false
	}
	if c.useFlags && !c.argsFiltered {
		return false
	}
	// Either post registration cfg isn't being used, or everything is set.
	return true
}

// SetUsage sets flagSet.Usage
func SetUsage(f func()) { appCfg.SetUsage(f) }
func (c *Cfg) SetUsage(f func()) {
	c.RWMutex.Lock()
	c.flagSet.Usage = f
	c.RWMutex.Unlock()
}

// SetName set's the cfg's name.
func SetName(name string) { appCfg.SetName(name) }
func (c *Cfg) SetName(name string) {
	c.RWMutex.Lock()
	c.name = name
	c.RWMutex.Unlock()
}

// Name returns the cfg's name.
func Name() string { return appCfg.Name() }
func (c *Cfg) Name() string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.name
}

// IsCore returns whether the passed setting is a core setting.
func IsCoreE(name string) (bool, error) { return appCfg.IsCoreE(name) }
func (c *Cfg) IsCoreE(name string) (bool, error) {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	s, ok := c.settings[name]
	if !ok {
		return false, fmt.Errorf("IsCore: setting not found: %q", name)
	}
	return s.IsCore, nil
}

func IsCore(name string) bool { return appCfg.IsCore(name) }
func (c *Cfg) IsCore(name string) bool {
	b, _ := c.IsCoreE(name)
	return b
}

// IsCfg returns whether the passed setting is a cfg setting.
func IsCfgE(name string) (bool, error) { return appCfg.IsCfgE(name) }
func (c *Cfg) IsCfgE(name string) (bool, error) {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	s, ok := c.settings[name]
	if !ok {
		return false, fmt.Errorf("IsCfg: setting not found: %q", name)
	}
	return s.IsCfg, nil
}

func IsCfg(name string) bool { return appCfg.IsCfg(name) }
func (c *Cfg) IsCfg(name string) bool {
	b, _ := c.IsCfgE(name)
	return b
}

// IsEnv returns whether the passed setting is a env setting.
func IsEnvE(name string) (bool, error) { return appCfg.IsEnvE(name) }
func (c *Cfg) IsEnvE(name string) (bool, error) {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	s, ok := c.settings[name]
	if !ok {
		return false, fmt.Errorf("IsEnv: setting not found: %q", name)
	}
	return s.IsEnv, nil
}

func IsEnv(name string) bool { return appCfg.IsEnv(name) }
func (c *Cfg) IsEnv(name string) bool {
	b, _ := c.IsEnvE(name)
	return b
}

// IsFlag returns whether the passed setting is a flag setting.
func IsFlagE(name string) (bool, error) { return appCfg.IsFlagE(name) }
func (c *Cfg) IsFlagE(name string) (bool, error) {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	s, ok := c.settings[name]
	if !ok {
		return false, fmt.Errorf("IsFlag: setting not found: %q", name)
	}
	return s.IsFlag, nil
}

func IsFlag(name string) bool { return appCfg.IsFlag(name) }
func (c *Cfg) IsFlag(name string) bool {
	b, _ := c.IsFlagE(name)
	return b
}

// processCfg() is the entry point for reading the configuration file.
func (c *Cfg) processCfg(buff []byte) (cfg interface{}, err error) {
	// if it's not set to use a cfg file, nothing to do
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	if !c.useCfg {
		return nil, nil
	}
	setting, ok := c.settings[CfgFile]
	if !ok {
		// Wasn't configured, nothing to do. Not an error.
		return nil, nil
	}
	n := setting.Value.(string)
	if n == "" {
		// This isn't an error as config file is allowed to not exist
		// TODO:
		//	Possible add a CfgFileRequired flag
		return nil, nil
	}
	// This shouldn't happend, but lots of things happen that shouldn't.  It should
	// have been registered already. so if it doesn't exit, err.
	format, ok := c.settings[CfgFormat]
	if !ok {
		return nil, fmt.Errorf("processCfg error: format was not set")
	}
	if format.Value.(string) == "" {
		return nil, fmt.Errorf("processCfg error: format was not set")
	}
	format, _ = c.settings[CfgFormat]
	cfg, err = unmarshalCfgBytes(ParseFormat(format.Value.(string)), buff)
	if err != nil {
		return nil, fmt.Errorf("processCfg unmarshal error, %s: %s", n, err.Error())
	}
	return cfg, nil
}

// canUpdate checks to see if the passed setting key is updateable. If the key
// is not updateable, a false is returned along with an error.
func (c *Cfg) canUpdate(k string) (bool, error) {
	// See if the key exists, if it doesn't already exist, it can't be updated.
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	s, ok := c.settings[k]
	if !ok {
		return false, fmt.Errorf("cannot update %q: setting not found", k)
	}
	// See if there are any settings that prevent it from being overridden.  Core and
	// environment variables are never settable. Core must be set during registration.
	if s.IsCore {
		return false, fmt.Errorf("cannot update %q: core settings cannot be updated", k)
	}
	if s.IsFlag && c.argsFiltered {
		return false, fmt.Errorf("cannot update %q: flag settings cannot be updated after arg filtering", k)
	}
	// Everything else is updateable.
	return true, nil
}

// canOverride() checks to see if the setting can be overridden. Overrides only
// come from flags. If it can't be overridden, it must be set via application,
// environment variable, or cfg file.
func (c *Cfg) canOverride(k string) bool {
	// an empty key cannot Override
	if k == "" {
		return false
	}
	// See if the key exists, if it doesn't already exist, it can't be overridden
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	s, ok := c.settings[k]
	if !ok {
		return false
	}
	// See if there are any settings that prevent it from being overridden.
	// Core can never be overridden-must be a flag to override.
	if s.IsCore {
		return false
	}
	// flags can only be set prior to arg filtering, after which you must use
	// Override().
	if s.IsFlag && c.argsFiltered {
		return false
	}
	return true
}

// getFileBytes reads from the passed path and returns its contents as bytes,
// or an error.  The entire contents of the file are read at once.
func getFileBytes(p string) ([]byte, error) {
	cfg, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// formatFromFilename gets the format from the passed filename.  An error will
// be returned if either the format isn't supported or the extension doesn't
// exist.  If the passed string has multiple dots, the last dot is assumed to
// be the extension.
func formatFromFilename(s string) (Format, error) {
	if s == "" {
		return Unsupported, fmt.Errorf("no config filename")
	}
	parts := strings.Split(s, ".")
	format := ""
	// case 0 has already been evaluated
	switch len(parts) {
	case 1:
		return Unsupported, fmt.Errorf("unable to determine %s's config format: no extension", strings.TrimSpace(s))
	case 2:
		format = parts[1]
	default:
		// assume its the last part
		format = parts[len(parts)-1]
	}
	f := ParseFormat(format)
	if !f.isSupported() {
		return Unsupported, unsupportedFormatErr(format)
	}
	return f, nil
}

// unmarshalCfgBytes accepts bytes and unmarshals them using the correct
// format. Either the unmarshaled data or an error is returned.
//
// Supported formats:
//   json
//   toml
// TODO
//   add YAML support
//   add HCL support
func unmarshalCfgBytes(f Format, buff []byte) (interface{}, error) {
	var ret interface{}
	switch f {
	case JSON:
		err := json.Unmarshal(buff, &ret)
		if err != nil {
			return nil, err
		}
	case TOML:
		_, err := toml.Decode(string(buff), &ret)
		if err != nil {
			return nil, err
		}
	default:
		err := unsupportedFormatErr(f.String())
		return nil, err
	}
	return ret, nil
}
