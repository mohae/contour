package contour

import (
	"fmt"
	_ "os"
	"sync"

	"github.com/mohae/flag"
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

	lock    sync.RWMutex
	flagSet *flag.FlagSet
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
	useCfg bool
	cfgSet bool

	// useEnv: whether this config writes to and reads from environment
	// variables. If false, Settings are stored only in Config.
	useEnv bool
	envSet bool

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
	return &Cfg{name: name, settings: map[string]*setting{}}
}

// Code returns the code for the config. If set, this is used as
// the prefix for environment variables and configuration setting names.
func (c *Cfg) Code() string {
	return c.code
}

func (c *Cfg) UseEnv() bool {
	return c.useEnv
}

// SetCfg goes through the initialized Settings and updates the updateable
// Settings if a new, valid value is found. This applies to, in order: Env
// variables and config files. For any that are not found, or that are
// immutable, once set, the original initialization values are used.
//
// The merged configuration Settings are then  written to their respective
// environment variables. At this point, only args, or in application setting
// changes, can change the non-immutable Settings.
func (c *Cfg) SetCfg() error {
	// Load any set environment variables into appConfig. Core and already
	// set Write Once Settings are not updated from env.
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

	c.cfgSet = true
	return nil
}

func (c *Cfg) setFromFile() error {
	f, err := c.getFile()
	if err != nil {
		fmt.Println(err)
		logger.Error(err)
		return err
	}

	// Go through the file contents and update the Cfg
	for k, v := range f {
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

/*
// Set env accepts a key and value and sets a single environment variable from that
func (c *Cfg) Setenv(k string, v interface{}) error {
	// if we aren't using environment variables, do nothing.
	if !appConfig.UseEnv() {
		return nil
	}

	var tmp string
	var err error

	switch appConfig.Settings[k].Type {
	case "string":
		err = os.Setenv(k, *v.(*string))

	case "int":
		err = os.Setenv(k, string(*v.(*int)))

	case "bool":
		tmp = strconv.FormatBool(*v.(*bool))
		err = os.Setenv(k, tmp)

	default:
		err = fmt.Errorf("Unable to set env variable for %s: type is unsupported %s", k, appConfig.Settings[k].Type)
	}
	return err
}
*/

// SetCode set's the code for this configuration. This can only be done once.
// If it is already set, it will return an error.
func (c *Cfg) SetCode(s string) error {
	if c.code != "" {
		return fmt.Errorf("appCode is already set. AppCode is immutable. Once set, it cannot be altered")
	}

	c.code = s
	return nil
}

// CfgProcessed determines whether, or not, all of the configurations, for a
// given config, have been processed.
func (c *Cfg) CfgProcessed() bool {
	if c.useCfg && !c.cfgSet {
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

// Convenience functions for the main config
// Code returns the code for the config. If set, this is used as
// the prefix for environment variables and configuration setting names.
func Code() string {
	return appCfg.code
}

func UseEnv() bool {
	return appCfg.useEnv
}

// SetCfg goes through the initialized Settings and updates the updateable
// Settings if a new, valid value is found. This applies to, in order: Env
// variables and config files. For any that are not found, or that are
// immutable, once set, the original initialization values are used.
//
// The merged configuration Settings are then  written to their respective
// environment variables. At this point, only args, or in application setting
// changes, can change the non-immutable Settings.
func SetCfg() error {
	return appCfg.SetCfg()
}

/*
// Set env accepts a key and value and sets a single environment variable from that
func (c *Cfg) Setenv(k string, v interface{}) error {
	// if we aren't using environment variables, do nothing.
	if !appConfig.UseEnv() {
		return nil
	}

	var tmp string
	var err error

	switch appConfig.Settings[k].Type {
	case "string":
		err = os.Setenv(k, *v.(*string))

	case "int":
		err = os.Setenv(k, string(*v.(*int)))

	case "bool":
		tmp = strconv.FormatBool(*v.(*bool))
		err = os.Setenv(k, tmp)

	default:
		err = fmt.Errorf("Unable to set env variable for %s: type is unsupported %s", k, appConfig.Settings[k].Type)
	}
	return err
}
*/

// SetCode set's the code for this configuration. This can only be done once.
// If it is already set, it will return an error.
func SetCode(s string) error {
	if appCfg.code != "" {
		return fmt.Errorf("appCode is already set. AppCode is immutable. Once set, it cannot be altered")
	}

	appCfg.code = s
	return nil
}

// Config processed returns whether or not all of the config's settings have
// been processed.
func CfgProcessed() bool {
	return appCfg.CfgProcessed()
}

// SetFlagSetUsage sets flagSet.Usage
func (c *Cfg) SetFlagSetUsage(f func()) {
	c.flagSet.Usage = f
}

// SetFlagSetUsage sets flagSet.Usage
func SetFlagSetUsage(f func()) {
	//	appCfg.SetFlagSetUsage(f)
	appCfg.flagSet.Usage = f
}
