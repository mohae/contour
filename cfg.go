package contour

import (
	"flag"
	"fmt"
	_ "os"
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
	searchPath      bool
	flagSet         *flag.FlagSet
	// code is the shortcode for this configuration. It is mostly used to
	// prefix environment variables, when used.
	code string
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
	// useEnv: whether this config writes to and reads from environment
	// variables. If false, Settings are stored only in Config.
	useEnv   bool
	envSet   bool
	envNames map[string]string // maps flag vars to environment names
	// Whether flags have been registered and set.
	useFlags bool
	flagsSet bool
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
	return &Cfg{name: name, errOnMissingCfg: true, searchPath: true, flagSet: flag.NewFlagSet(name, flag.ContinueOnError), settings: map[string]*setting{}}
}

// Code is a convenience functions for the appCfg global.
func Code() string {
	return appCfg.Code()
}

// Code returns the code for the config. If set, this is used as
// the prefix for environment variables and configuration setting names.
func (c *Cfg) Code() string {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.code
}

// SetCode is a convenience functions for the appCfg global.
func SetCode(s string) error {
	return appCfg.SetCode(s)
}

// SetCode set's the code for this configuration. This can only be done once.
// If it is already set, it will return an error.
func (c *Cfg) SetCode(s string) error {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	if c.code != "" {
		return fmt.Errorf("this configuration's code is already set and cannot be overridden")
	}
	c.code = s
	return nil
}

// Code is a convenience functions for the appCfg global.
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

// SetErrOnMissingCfg is a convenience functions for the appCfg global.
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

// SearchPath is a convenience functions for the appCfg global.
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

// SetSearchPath is a convenience functions for the appCfg global.
func SetSearchPath(b bool) {
	appCfg.SetSearchPath(b)
}

func (c *Cfg) SetSearchPath(b bool) {
	c.RWMutex.Lock()
	c.searchPath = b
	c.RWMutex.Unlock()
}

// UseCfgFile is a convenience functions for the appCfg global.
func UseCfgFile() bool {
	return appCfg.UseCfgFile()
}

// UseCfgFile returns whether this cfg uses a CfgFile.
func (c *Cfg) UseCfgFile() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.useCfgFile
}

// SetUseCfgFile is a convenience functions for the appCfg global.
func SetUseCfgFile(b bool) {
	appCfg.SetUseCfgFile(b)
}

// SetErrOnMissingCfg returns whether a missing config file should result in
// an error.
func (c *Cfg) SetUseCfgFile(b bool) {
	c.RWMutex.Lock()
	c.useCfgFile = b
	c.RWMutex.Unlock()
}

// UseEnv is a convenience functions for the appCfg global.
func UseEnv() bool {
	return appCfg.useEnv
}

func (c *Cfg) UseEnv() bool {
	c.RWMutex.RLock()
	defer c.RWMutex.RUnlock()
	return c.useEnv
}

// SetUseEnv is a convenience functions for the appCfg global.
func SetUseEnv(b bool) {
	appCfg.SetUseEnv(b)
}

func (c *Cfg) SetUseEnv(b bool) {
	c.RWMutex.Lock()
	c.useEnv = b
	c.RWMutex.Unlock()
}

// SetCfg is a convenience function for the global appCfg.
func SetCfgFile() error {
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
	// Load any set environment variables into appConfig. Core and already set
	// Write Once Settings are not updated from env.
	//	c.loadEnvs()

	// Set all the environment variables. This is the application Settings
	// merged with any already existing environment variable values.
	//	err := c.Setenvs()
	//	if err != nil {
	//		return err
	//	}

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

// SetUsage sets flagSet.Usage
func (c *Cfg) SetUsage(f func()) {
	c.RWMutex.Lock()
	c.flagSet.Usage = f
	c.RWMutex.Unlock()
}

// SetUsage sets appCfg's usage func
func SetUsage(f func()) {
	appCfg.SetUsage(f)
}
