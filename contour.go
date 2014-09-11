//
package contour

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// Environment Variable constants for common environment variables.
// Or you can supply your own values. Contour automatically downcases
// environment variables for consistency across formats.
var (
	EnvConfigFilename    string = "configfilename"
	EnvConfigFormat      string = "configformat"
	EnvLogConfigFilename string = "logconfigfilename"
	EnvLogging           string = "logging"
)

// settingAlias are aliases to settings, each setting is its own alias.
var settingAlias map[string]string = make(map[string]string)

// commandAlias are aliases to commands, each command is its own alias.
var commandAlias map[string]string = make(map[string]string)

// configFile holds the contents of the configuration file
var configFile map[string]interface{} = make(map[string]interface{})

// config holds the current application configuration
var AppConfig *Config = &Config{settings: map[string]*setting{}}

const (
	settingNotFoundErr = " setting was not found"
)

// Config is a group of settings and holds all of the application setting
// information. Even though contour automatically uses environment variables,
// unless its told to ignore them, it still needs to maintain state
// information about each setting so it knows how to handle attempst to update.
// TODO:
//	* support ignoring environment variables
//
type Config struct {

	// code is the shortcode for this configuration. It is mostly used to
	// prefix environment variables, when used.
	code     string
	settings map[string]*setting
}

// GetAppConfig returns the AppConfig to the caller. Any contour function
// called uses this config.
func GetAppConfig() *Config {
	return AppConfig
}

// NewConfig returns a *Config to the caller. Any config created by NewConfig()
// is independent of the AppConfig
func NewConfig() *Config {
	return &Config{settings: map[string]*setting{}}
}

// GetAppCode returns the app code for the config. If set, this is used as
// the prefix for environment variables and configuration setting names.
func (c *Config) GetAppCode() string {
	return c.code
}

// SetAppCode set's the appcode. This can only be done once. If it is already
// set, it will return an error.
func (c *Config) SetAppCode(s string) error {
	if c.code != "" {
		return errors.New("appCode is already set. AppCode is immutable. Once set, it cannot be altered")
	}

	c.code = s
	return nil
}

// setting holds the information for a configuration setting.
type setting struct {
	// Type is the datatype for the setting
	Type string

	// The current value of the setting
	Value interface{}

	// Code of the setting
	Code string

	// Immutable: Once the Value has been set, it cannot be changed. This
	// allows for registering a setting without a value, so it can be
	// updated later--becoming immutable in the process.
	Immutable bool

	// IsCore: whether or not this is considered a core setting. Core
	// settings if for things like application name, where you don't want
	// anything else overwriting that value, once set, and you want to be
	// able to overwrite any existing ENV value if contour hasn't already
	// set it. Once set, Immutable is also true.
	IsCore bool

	// IsEnv: whether or not the original source of this setting was its
	// environment variables, vs. flags or config, etc. This is tracked
	// because it has implications on override behavior.
	IsEnv bool

	// IsFlag:  whether or not this is a flag.
	IsFlag bool
}

// SetConfig goes through the initialized settings and updates the updateable
// settings if a new, valid value is found. This applies to, in order: Env
// variables and config files. For any that are not found, or that are
// immutable, once set, the original initialization values are used.
//
// The merged configuration settings are then  written to their respective
// environment variables. At this point, only args, or in application setting
// changes, can change the non-immutable settings.
func SetConfig() error {
	// Load any set environment variables into AppConfig. Core and already
	// set Write Once settings are not updated from env.
	loadEnvs()

	// Set all the environment variables. This is the application settings
	// merged with any already existing environment variable values.
	err := SetEnvs()
	if err != nil {
		return err
	}

	// Load the Config file.
	err = loadConfigFile()
	if err != nil {
		return err
	}

	//  Save the config file settings to their env variables, if allowed.
	return setEnvFromConfigFile()
}

// getConfigFormat gets the configured config filename and returns the format
// it is in, if it is a supported format; otherwise an error.
func getConfigFormat(s string) (string, error) {
	if s == "" {
		return "", errors.New("a config filename was expected, none received")
	}

	parts := strings.Split(s, ".")
	format := ""

	// case 0 has already been evaluated
	switch len(parts) {
	case 1:
		return "", errors.New("unable to determine config format, the configuration file, " + strings.TrimSpace(s) + ", doesn't have an extension")

	case 2:
		format = parts[1]

	default:
		// assume its the last part
		format = parts[len(parts)-1]
	}

	if !isSupportedFormat(format) {
		return "", errors.New(format + " is an unsupported format for configuration files")
	}

	return format, nil

}

// isSupportedFormat checks to see if the passed string represents a supported
// config format.
func isSupportedFormat(s string) bool {
	switch s {
	case "json":
		return true
	case "toml":
		return true
	default:
		return false
	}

	return false
}

// loadEnvs updates the configuration from the environment variable values.
// A setting is only updated if it IsUpdateable.
func loadEnvs() {
	for k, _ := range AppConfig.settings {
		// See if k exists as an env variable
		v := os.Getenv(k)
		if v == "" {
			continue
		}

		// Core is not updateable
		if AppConfig.settings[k].IsCore {
			continue
		}

		// If its readonly, see if its set. If it isn't it can be.
		if AppConfig.settings[k].Immutable {
			if AppConfig.settings[k].Value != nil {
				continue
			}
		}

		// Gotten this far, set it
		AppConfig.settings[k].Value = v
		AppConfig.settings[k].IsEnv = true
	}

}

// loadConfigFile() is the entry point for reading the configuration file.
func loadConfigFile() error {
	setting, ok := AppConfig.settings[EnvConfigFilename]
	if !ok {
		// Wasn't configured, nothing to do. Not an error.
		return nil
	}

	n := setting.Value.(string)
	if n == "" {
		// This isn't an error as config file is allowed to not exist
		// TODO:
		//	Possible add a ConfigFileRequired flag
		//	Add logging support
		return nil
	}

	// This shouldn't happend, but lots of things happen that shouldn't.
	// It should have been registered already. so if it doesn't exit, err.
	if AppConfig.settings[EnvConfigFormat].Value == nil {
		return errors.New("Unable to load configuration value, its format type was not set")
	}

	fBytes, err := readConfigFile(n)
	if err != nil {
		return err
	}

	err = marshalFormatReader(AppConfig.settings[EnvConfigFormat].Value.(string), bytes.NewReader(fBytes))
	if err != nil {
		return err
	}

	return nil
}

// readConfigFile reads the configFile
func readConfigFile(n string) ([]byte, error) {
	cfg, err := ioutil.ReadFile(n)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// marshalFormatReader
func marshalFormatReader(t string, r io.Reader) error {
	b := new(bytes.Buffer)
	b.ReadFrom(r)

	switch t {
	case "json":
		err := json.Unmarshal(b.Bytes(), &configFile)
		if err != nil {
			return err
		}

	case "toml":
		_, err := toml.Decode(b.String(), &configFile)
		if err != nil {
			return err
		}

	}
	return nil
}

// canUpdate checks to see if the passed setting key is updateable.
func canUpdate(k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be
	// updated.
	_, ok := AppConfig.settings[k]
	if !ok {
		return false
	}

	// See if there are any settings that prevent it from being overridden.
	// Core and Environment variables are never settable. Core must be set
	// during registration and Environment variables must be set using
	// the Override functions.
	if AppConfig.settings[k].IsCore || AppConfig.settings[k].IsEnv {
		return false
	}

	// Immutable variables are only settable if they are not set.
	// This does not apply to boolean as there is no way to determine if
	// the value is unset. So bool immutables are only writable when they
	// are registered.
	if (AppConfig.settings[k].Immutable && AppConfig.settings[k].Value != "") || (AppConfig.settings[k].Immutable && AppConfig.settings[k].Type == "bool") {
		return false
	}
	return true
}

// canOverride() checks to see if the setting can be overridden. Overrides
// only come from args and flags. ConfigFile settings must be set instead.
func canOverride(k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be
	// overridden
	_, ok := AppConfig.settings[k]
	if !ok {
		return false
	}

	// See if there are any settings that prevent it from being overridden.
	if AppConfig.settings[k].IsCore {
		return false
	}

	// Immutable variables are only settable if they are not set.
	// This does not apply to boolean as there is no way to determine if
	// the value is unset. So bool immutables are only writable when they
	// are registered.
	if (AppConfig.settings[k].Immutable && AppConfig.settings[k].Value != "") || (AppConfig.settings[k].Immutable && AppConfig.settings[k].Type == "bool") {
		return false
	}

	return true
}

// AddCommandAlias adds an alias for a command. The first time a command is
// added, it's added as an alias of itself too.
func AddCommandAlias(command, alias string) error {
	// see if an alias already exists
	v, ok := commandAlias[alias]
	if ok {
		return errors.New(alias + " is an alias of the command " + v + " cannot make it an alias of " + command)
	}

	// see if the command already has aliases
	v, ok = commandAlias[command]
	if !ok {
		// Add it as an alias of itself first
		commandAlias[command] = command
	}

	commandAlias[alias] = command
	return nil
}

// AddSettingAlias adds an alias for a setting. The first time a setting is
// added, it's added as an alias of itself too.
func AddSettingAlias(setting, alias string) error {
	// see if an alias already exists
	v, ok := settingAlias[alias]
	if ok {
		err := errors.New(alias + " is an alias of the setting " + v + " cannot make it an alias of " + setting)
		return err
	}

	// see if the setting already has aliases
	v, ok = settingAlias[setting]
	if !ok {
		// Add it as an alias of itself first
		settingAlias[setting] = setting
	}

	settingAlias[alias] = setting
	return nil
}

// Set
// resetAppConfig resets the application's configuration struct to empty.
// This does not affect their respective environment variables
func resetAppConfig() {
	AppConfig = &Config{settings: map[string]*setting{}}
}

// notFoundErr returns a standardized not found error.
func notFoundErr(k string) error {
	return errors.New(k + " not found")
}

// FilterArgs takes the passed args and filter's the flags out of them.
// The populated flags override their settings, according to the override
// rules. Successful overrides result in the relevant AppConfig setting
// being updated along with the environment variable.
//
// Any args left, after filtering, are returned to the caller.
// TODO: refactor this for greater abstraction (as long as it doesn't come
//	at the cost of reflection)
func FilterArgs(flagSet *flag.FlagSet, args []string) ([]string, error) {
	// Get the flag filters from the config variable information.
	boolFilterNames := GetBoolFilterNames()

	// Preallocate the worst case scenario.
	boolFilters := make([]*bool, len(boolFilterNames))
	bFilterNames := make([]string, len(boolFilterNames))
	var flags int

	for _, name := range boolFilterNames {
		if AppConfig.settings[name].IsFlag {
			boolFilters[flags] = flagSet.Bool(name, AppConfig.settings[name].Value.(bool), fmt.Sprintf("filter %s", name))
			bFilterNames[flags] = name
			flags++
		}
	}

	// Get the flag filters from the config variable information.
	intFilterNames := GetIntFilterNames()

	// Preallocate the worst case scenario.
	intFilters := make([]*int, len(intFilterNames))
	iFilterNames := make([]string, len(intFilterNames))
	flags = 0

	for _, name := range intFilterNames {
		if AppConfig.settings[name].IsFlag {
			intFilters[flags] = flagSet.Int(name, AppConfig.settings[name].Value.(int), fmt.Sprintf("filter %s", name))
			iFilterNames[flags] = name
			flags++
		}
	}
	// Get the flag filters from the config variable information.
	stringFilterNames := GetStringFilterNames()

	// Preallocate the worst case scenario.
	stringFilters := make([]*string, len(stringFilterNames))
	sFilterNames := make([]string, len(stringFilterNames))
	flags = 0

	for _, name := range stringFilterNames {
		if AppConfig.settings[name].IsFlag {
			fmt.Println(flags)
			stringFilters[flags] = flagSet.String(name, AppConfig.settings[name].Value.(string), fmt.Sprintf("filter %s", name))
			sFilterNames[flags] = name
			flags++
		}
	}

	// Parse args for flags
	err := flagSet.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse of command-line arguments failed: %s", err)
	}

	// Get the remaining args
	cmdArgs := flagSet.Args()

	// Process the captured values
	for i, v := range boolFilters {
		if v != AppConfig.settings[bFilterNames[i]].Value {
			Override(bFilterNames[i], v)
		}
	}

	for i, v := range intFilters {
		if v != AppConfig.settings[iFilterNames[i]].Value {
			Override(iFilterNames[i], v)
		}
	}

	for i, v := range stringFilters {
		if v != AppConfig.settings[sFilterNames[i]].Value {
			Override(sFilterNames[i], v)
		}
	}

	return cmdArgs, nil
}
