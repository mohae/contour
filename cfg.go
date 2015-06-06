package contour

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
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
	settings map[string]*setting
	// Whether configuration settings have been registered and set.
	useCfgFile bool
	cfgFileSet bool
	// tracks the vars that are exposed to cfg
	cfgVars map[string]struct{}
	// useEnv: whether this config writes to and reads from environment
	// variables. If false, Settings are stored only in Config.
	useEnv bool
	envSet bool
	// Whether flags have been registered and set.
	useFlags bool
	flagsSet bool
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
	return &Cfg{name: name, errOnMissingCfg: true, searchPath: true, flagSet: flag.NewFlagSet(name, flag.ContinueOnError), settings: map[string]*setting{}, cfgVars: map[string]struct{}{}, useFlags: true, shortFlags: map[string]string{}}
}

// Loadenv is a convenience function for the global appCfg.
func Loadenv() error {
	return appCfg.Loadenv()
}

// Loadenv, if cfg.useEnvs, checks the cfg's env vars and updates the settings
// if they are set.
func (c *Cfg) Loadenv() error {
	if !c.useEnv {
		return nil
	}
	var err error
	for k, v := range c.settings {
		if !v.IsEnv {
			continue
		}
		tmp := os.Getenv(fmt.Sprintf("%s_%s", c.name, k))
		if tmp != "" {
			switch v.Type {
			case "bool":
				b, _ := strconv.ParseBool(tmp)
				err = c.UpdateBoolE(k, b)
			case "int":
				i, err := strconv.Atoi(tmp)
				if err != nil {
					return fmt.Errorf("Loadenv error while parsing %s: %s", fmt.Sprintf("%s_%s", c.name, k), err)
				}
				err = c.UpdateIntE(k, i)
			case "int64":
				i, err := strconv.ParseInt(tmp, 10, 64)
				if err != nil {
					return fmt.Errorf("Loadenv error while parsing %s: %s", fmt.Sprintf("%s_%s", c.name, k), err)
				}
				err = c.UpdateInt64E(k, i)
			case "string":
				err = c.UpdateStringE(k, tmp)
			default:
				return fmt.Errorf("%s has an unsupported env variable type: %s", k, v.Type)
			}
			if err != nil {
				return fmt.Errorf("Loadenv error while setting %s: %s", fmt.Sprintf("%s_%s", c.name, k), err)
			}
		}
	}
	c.envSet = true
	return nil
}

// Code is a convenience function for the global appCfg.
func ErrOnMissingCfg() bool {
	return appCfg.ErrOnMissingCfg()
}

// ErrOnMissingCfg returns whether a missing config file should result in an
// error. This only applies when useCfg == true
func (c *Cfg) ErrOnMissingCfg() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.errOnMissingCfg
}

// SetErrOnMissingCfg is a convenience function for the global appCfg.
func SetErrOnMissingCfg(b bool) {
	appCfg.SetErrOnMissingCfg(b)
}

// SetErrOnMissingCfg returns whether a missing config file should result in an
// error. This only applies when useCfg == true
func (c *Cfg) SetErrOnMissingCfg(b bool) {
	c.RWMutex.Lock()
	c.errOnMissingCfg = b
	c.RWMutex.Unlock()
}

// SearchPath is a convenience function for the global appCfg.
func SearchPath() bool {
	return appCfg.SearchPath()
}

// SearchPath returns whether or not the Path environment variable should be
// searched when looking for the config file.
func (c *Cfg) SearchPath() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.searchPath
}

// SetSearchPath is a convenience function for the global appCfg.
func SetSearchPath(b bool) {
	appCfg.SetSearchPath(b)
}

// SetSearchPath set's whether or not the user's PATH env variable should be
// searched for the cfg file.
func (c *Cfg) SetSearchPath(b bool) {
	c.RWMutex.Lock()
	c.searchPath = b
	c.RWMutex.Unlock()
}

// UseCfgFile is a convenience function for the global appCfg.
func UseCfgFile() bool {
	return appCfg.UseCfgFile()
}

// UseCfgFile returns whether this cfg uses a CfgFile.
func (c *Cfg) UseCfgFile() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.useCfgFile
}

// SetUseCfgFile is a convenience function for the global appCfg.
func SetUseCfgFile(b bool) {
	appCfg.SetUseCfgFile(b)
}

// SetUseCfgFile set's whether a cfg file should be used with this cfg.
func (c *Cfg) SetUseCfgFile(b bool) {
	c.RWMutex.Lock()
	c.useCfgFile = b
	c.RWMutex.Unlock()
}

// UseEnv is a convenience function for the global appCfg.
func UseEnv() bool {
	return appCfg.useEnv
}

// UseEnv is whether or not environment variables are used.
func (c *Cfg) UseEnv() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.useEnv
}

// SetUseEnv is a convenience function for the global appCfg.
func SetUseEnv(b bool) {
	appCfg.SetUseEnv(b)
}

// SetUseEnv set's whether or not environment variables should be used with
// this cfg.
func (c *Cfg) SetUseEnv(b bool) {
	c.RWMutex.Lock()
	c.useEnv = b
	c.RWMutex.Unlock()
}

// SetCfg is a convenience function for the global appCfg.
func SetCfg() error {
	return appCfg.SetCfg()
}

// SetCfg goes through the initialized Settings and updates the updateable
// settings, if a new, valid value is found.  This applies to, in order: Env
// variables and config files. For any that are not found, or that are
// immutable, once set, the original initialization values are used.
//
// The merged configuration Settings are then  written to their respective
// environment variables. At this point, only args, or in application setting
// changes, can change the non-immutable Settings.
func (c *Cfg) SetCfg() error {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	// Load the Config file.
	err := c.setFromFile()
	if err != nil {
		return err
	}
	c.cfgFileSet = true
	return nil
}

func (c *Cfg) setFromFile() error {
	f, err := c.getFile()
	if err != nil {
		return err
	}
	// if nothing was returned and no error, nothing to do
	if f == nil {
		return nil
	}
	// Go through the file contents and update the Cfg
	for k, v := range f.(map[string]interface{}) {
		// Find the key in the settings
		_, ok := c.settings[k]
		if !ok {
			// skip settings that don't already exist
			continue
		}
		err := c.updateE(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// CfgFileProcessed sis a convenience function for the global appCfg.
func CfgFileProcessed() bool {
	return appCfg.CfgFileProcessed()
}

// CfgFileProcessed determines whether, or not, all of the configurations, for a
// given config, have been processed.
func (c *Cfg) CfgFileProcessed() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	if c.useCfgFile && !c.cfgFileSet {
		return false
	}
	if c.useEnv && !c.envSet {
		return false
	}
	if c.useFlags && !c.flagsSet {
		return false
	}
	// Either post registration configuration isn't being used, or
	// everything is set.
	return true
}

// SetUsage sets appCfg's usage func
func SetUsage(f func()) {
	appCfg.SetUsage(f)
}

// SetUsage sets flagSet.Usage

func (c *Cfg) SetUsage(f func()) {
	c.RWMutex.Lock()
	c.flagSet.Usage = f
	c.RWMutex.Unlock()
}
