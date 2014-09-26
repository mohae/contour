package contour

import (
	"bytes"
	"fmt"
	_ "os"
	"sync"
)

// AppConfig returns the configs[app]. If it doesn't exist, one is initialized
// and returned.
//
// Contour has a set of functions that implicitly interact with configs[app].
// If the application is only going to use one configuration, this is what
// should be used as one can just interact with contour, instead of directly
// with the app config, which is also supported.
func AppConfig() *Cfg {
	return configs[0]
}

// Config returns the config for the passed key, if it exists, or an error.
func Config(k string) (*Cfg, error) {
	i, err := configIndex(k)
	if err != nil {
		return nil, err
	}

	return configs[i], nil
}

// NewConfig returns a *Cfg to the caller. This config is added to configs
// using the passed key value. If a config using the requested key already
// exists, an error is returned.
func NewConfig(k string) (c *Cfg, err error) {
	i, err := configIndex(k)
	if err == nil {
		return configs[i], fmt.Errorf("%q configuration already exists", k)
	}

	// Grow, if needed
	if cap(configs) == configCount {
		inc := growConfigIncrement	
		// 0 means double each time
		if inc == 0 {
			inc = configCount * 2
		}

		tempCfgs := configs
		configs = make([]*Cfg, configCount + 1, configCount + inc)
		copy(configs, tempCfgs[:configCount])

		tempNames := configNames
		configNames = make([]string, configCount + 1, configCount + inc)
		copy(configNames, tempNames[:configCount])
	}

	configNames[configCount] = k
	configs[configCount] = &Cfg{name: k, settings: make([]*setting, initSettingsCap)}
	configCount++
	return configs[configCount - 1], nil
}

// Cfg is a group of Settings and holds all of the application setting
// information. Even though contour automatically uses environment variables,
// unless its told to ignore them, it still needs to maintain state
// information about each setting so it knows how to handle attempst to update.
// TODO:
//	* support ignoring environment variables
//
type Cfg struct {
	name string

	lock sync.RWMutex

	// code is the shortcode for this configuration. It is mostly used to
	// prefix environment variables, when used.
	code string

	// file is the name of the config file associated with this config.
	file string

	// encoding is what encoding scheme is used for this config file
	encoding string

	// Settings contains a slice of the configuration settings
	settings []*setting
	settingsCount int

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

// Code returns the code for the config. If set, this is used as
// the prefix for environment variables and configuration setting names.
func (c *Cfg) Code() string {
	return c.code
}

func (c *Cfg) UseEnv() bool {
	return c.useEnv
}

// SetConfig goes through the initialized Settings and updates the updateable
// Settings if a new, valid value is found. This applies to, in order: Env
// variables and config files. For any that are not found, or that are
// immutable, once set, the original initialization values are used.
//
// The merged configuration Settings are then  written to their respective
// environment variables. At this point, only args, or in application setting
// changes, can change the non-immutable Settings.
func (c *Cfg) SetConfig() error {
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

// setEnvFromConfigFile goes through all the settings in the configFile and
// checks to see if the setting is updateable; saving those that are to their
// environment variable.
func (c *Cfg) setCfg(cf map[string]interface{}) error {
	if !c.UseEnv() {
		return nil
	}
	
	for k, v := range cf {
		_, err := c.settingIndex(k)
		if err != nil {
			// skip settings that don't already exist
			continue
		}

		err = c.updateE(k, v)
		if err != nil {
			return err
		}

	}

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
		idx, err := c.settingIndex(k)
		if err != nil {
			// skip settings that don't already exist
			continue
		}

		err = c.idxUpdateE(idx, v)
		if err != nil {
			return err
		}

	}

	return nil
}

func (c *Cfg) appendSetting(s *setting, pos int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if cap(c.settings) == c.settingsCount {
		tmp := c.settings
		c.settings = make([]*setting, c.settingsCount + 1, c.settingsCount + growSettingsIncrement)
		copy(c.settings, tmp[:c.settingsCount])
	}
	c.settings[c.settingsCount] = s
	c.settingsCount++
}

func (c *Cfg) deleteSetting(s string) {
	c.lock.RLock()
	for i, setting := range c.settings {
		if setting.Name == s {
			// delete 
			c.lock.RUnlock()
			c.deleteSettingIndex(i)
			return
		}
	}
	c.lock.RUnlock()
	return
}

// deleteSettingIndex deletes the setting at pos by appending around it.
func (c *Cfg) deleteSettingIndex(idx int) {
	c.lock.Lock()
	c.settings = append(c.settings[:idx], c.settings[idx+1:]...)
	c.lock.Unlock()
}

func (c *Cfg) settingIndex(k string) (int, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for i, setting := range c.settings {
		if setting.Name == k {
			return i, nil
		}
	}

	return -1, notFoundErr(k)
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

// getCfgFile() is the entry point for reading the configuration file.
func (c *Cfg) getFile() (map[string]interface{}, error) {
	// Configuration file not existing is ok as running without config
	// file is supported.
	idx, err := c.settingIndex(CfgFile)
	if err != nil {
		return nil, nil
	}

	n := c.settings[idx].Value.(string)
	if n == "" {
		// This isn't an error as config file is allowed to not exist
		// TODO:
		//	Possible add a ConfigFileRequired flag
		//	Add logging support
		return nil, nil
	}

	// This shouldn't happend, but lots of things happen that shouldn't.
	// It should have been registered already. so if it doesn't exit, err.
	if c.settings[idx].Value == nil {
		return nil, fmt.Errorf("Unable to load the cfg file, the configuration format type was not set")
	}

	fBytes, err := readCfg(n)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	idx, err = c.settingIndex(CfgFormat)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	cfg := make(map[string]interface{})
	cfg, err = marshalFormatReader(ParseFormat(c.settings[idx].Value.(string)), bytes.NewReader(fBytes))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return cfg, nil
}

// SetCode set's the code for this configuration. This can only be done once.
// If it is already set, it will return an error.
func (c *Cfg) SetCode(s string) error {
	if c.code != "" {
		return fmt.Errorf("appCode is already set. AppCode is immutable. Once set, it cannot be altered")
	}

	c.code = s
	return nil
}

// ConfigProcessed determines whether, or not, all of the configurations, for a
// given config, have been processed.
func (c *Cfg) ConfigProcessed() bool {
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


// Config methods
// RegisterConfigFilename set's the configuration file's name. The name is
// parsed for a valid extension--one that is a supported format--and saves
// that value too. If it cannot be determined, the extension info is not set.
// These are considered core values and cannot be changed from command-line
// and configuration files. (IsCore == true).
func (c *Cfg) RegisterConfigFilename(k, v string) error {
	if v == "" {
		return fmt.Errorf("A config filename was expected, none received")
	}

	if k == "" {
		return fmt.Errorf("A key for the config filename setting was expected, none received")
	}

	c.RegisterStringCore(k, v)

	// Register it first. If a valid config format isn't found, an error
	// will be returned, so registering it afterwords would mean the
	// setting would not exist.
	c.RegisterString(CfgFormat, "")
	format, err := configFormat(v)
	if err != nil {
		return err
	}

	// Now we can update the format, since it wasn't set before, it can be
	// set now before it becomes read only.
	c.RegisterString(CfgFormat, format.String())

	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func (c *Cfg) RegisterSetting(Type, k, Short string, v interface{}, Usage string, Default string, IsCore, IsCfg, IsFlag bool) error  {

	// find the index for the setting. If its found, return an error
	_, err := c.settingIndex(k)
	if err == nil {
		err := fmt.Errorf("cannot register %q, it is already registered")
		logger.Error(err)
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	// Keep track of whether or not a config is being used. If a setting is
	// registered as a config setting, it is assumed a configuration source
	// is being used.
	if IsCfg {
		c.useCfg = true
	}

	// Keep track of whether or not flags are being used. If a setting is
	// registered as a flag setting, it is assumed that flags are being 
	// used.
	if IsFlag {
		c.useFlags = true
	}

	// Otherwise register it as a new setting.
	c.settings[c.settingsCount] = &setting{
		Type:      Type,
		Value:     v,
		Short:     Short,
		Usage:     Usage,
		Default:   Default,
		IsCore:    IsCore,
		IsCfg:     IsCfg,
		IsFlag:    IsFlag,
	}
	
	c.settingsCount++
	return nil
}


func (c *Cfg) updateE(k string, v interface{}) error {
	// setting must exist to update
	idx, err := c.settingIndex(k)
	if err != nil {
		logger.Error(err)
		return err
	}

	return c.idxUpdateE(idx, v)
}

// Only non-core settings are updateable.
// Flags must use Override* to update settings.
// save it to its environment variable.
func (c *Cfg) idxUpdateE(i int, v interface{}) error {
	if !c.canUpdate(i) {
		err := fmt.Errorf("config[%s]: %s is not updateable", c.name, c.settings[i].Name)
		logger.Warn(err)
		return err
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.settings[i].Value = v
	return nil
}

// canUpdate checks to see if the passed setting key is updateable.
// TODO the logic flow is wonky because it could be simplified but
// want hard check for core and not sure about conf/flag/env stuff yet.
// so the wierdness sits for now.
func (c *Cfg) canUpdate(i int) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	// See if there are any settings that prevent it from being overridden.
	// Core and ironment variables are never settable. Core must be set
	// during registration.
	if c.settings[i].IsCore {
		return false
	}

	// Only flags and conf types are updateable, otherwise they must be
	// registered or set.
	if c.settings[i].IsCfg || c.settings[i].IsFlag {
		return true
	}

	return false
}

func (c *Cfg) Override(k string, v interface{}) error {
	if v == nil {
		return nil
	}

	idx, err := c.settingIndex(k)
	if err != nil {
		logger.Error(err)
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	// If it can't be overriden, 
	if c.settings[idx].IsCore || !c.settings[idx].IsFlag {
		err := fmt.Errorf("%v: setting is not a flag. Only flags can be overridden", k)
		logger.Warn(err)
		return err
	}

/*
	// Write to environment variable
	err := c.Setenv(k, v)
	if err != nil {
		logger.Error(err)
		return err
	}
*/

	c.settings[idx].Value = v
	return nil
}

func Override(k string, v interface{}) error {
	return configs[0].Override(k, v)
}

// canOverride() checks to see if the setting can be overridden. Overrides
// only come from flags. If it can't be overridden, it must be set via
// application, environment variable, or cfg file.
func (c *Cfg) canOverride(i int) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	// See if there are any settings that prevent it from being overridden.
	// Core can never be overridden
	// Must be a flag to override.
	if c.settings[i].IsCore || !c.settings[i].IsFlag {
		return false
	}

	return true
}

// GetE returns the setting Value as an interface{}. If its not a valid
// setting, an error is returned.
func (c *Cfg) GetE(k string) (interface{}, error) {
	idx, err := c.settingIndex(k)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.settings[idx].Value, nil
}

// Convenience functions for configs[app]
// Get returns the setting Value as an interface{}.
// GetE returns the setting Value as an interface{}.
func GetE(k string) (interface{}, error) {
	return configs[0].GetE(k)
}

func (c *Cfg) Get(k string) interface{} {
	s, _ := c.GetE(k)
	return s
}

func Get(k string) interface{} {
	s, _ := configs[0].GetE(k)
	return s
}

// GetInterfaceE is a convenience wrapper function to Get
func (c *Cfg) GetInterfaceE(k string) (interface{}, error) {
	return c.GetE(k)
}

// GetInterfaceE is a convenience wrapper function to Get
func GetInterfaceE(k string) (interface{}, error) {
	return configs[0].GetE(k)
}

// GetInterfac returns the setting Value as an interface
func (c *Cfg) GetInterface(k string) interface{} {
	return c.Get(k)
}

// GetInterface is a convenience wrapper function to Get
func GetInterface(k string) interface{} {
	return configs[0].Get(k)
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
