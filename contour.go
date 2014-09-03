//
package contour

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	utils "github.com/mohae/utilitybelt"
)

// Environment Variable constants for common environment variables.
// Or you can supply your own values. Contour automatically downcases
// environment variables for consistency across formats.
var (
	EnvConfigFilename string = "configfilename"
	EnvConfigFormat string = "configformat"
	EnvLogFilename string = "logfilename"
	EnvLogging string = "logging"
)

// settingAlias are aliases to settings, each setting is its own alias.
var settingAlias map[string]string = make(map[string]string)

// commandAlias are aliases to commands, each command is its own alias.
var commandAlias map[string]string = make(map[string]string)

// configFile holds the contents of the configuration file
var configFile map[string]interface{} = make(map[string]interface{})

// config holds the current application configuration
var AppConfig *Config = &Config{Settings: map[string]*setting{}}

const (
	SettingNotFoundErr = " setting was not found"
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
	code 	string
	Settings 	map[string]*setting
}

// GetAppConfig returns the AppConfig to the caller. Any contour function
// called uses this config.
func GetAppConfig() *Config {
	return AppConfig
}

// NewConfig returns a *Config to the caller. Any config created by NewConfig()
// is independent of the AppConfig
func NewConfig() *Config {
	return &Config{Settings: map[string]*setting{}}
}

func (c *Config) GetCode() string {
	return c.code
}

// SetAppCode set's the appcode. This can only be done once. If it is already
// set, it will return an error.
func (c *Config) SetCode(s string) error {
	if c.code != "" {
		return errors.New("appCode is already set. AppCode is idempotent. Once set, it cannot be altered")
	}

	c.code = s

	return nil
}

// setting holds the information for a configuration setting.
type setting struct {
	// Code of the setting
	Code string

	// Type is the datatype for the setting
	Type string

	// The current value of the setting
	Value interface{}

	// IsFlag:  whether or not this is a flag.
	IsFlag bool

	// IsIdempotent: whether or not this value can be overwritten once set.
	IsIdempotent bool

	// SourceIsEnv: whether or not the original source of this setting was
	// its environment variables, vs. flags or config, etc. This is tracked
	// because it has implications on override behavior.
	SourceIsEnv bool

	// IsCore: whether or not this is considered a core setting. Core 
	// settings if for things like application name, where you don't want
	// anything else overwriting that value, once set, and you want to be
	// able to overwrite any existing ENV value if contour hasn't already
	// set it. Once set, IsIdempotent is also true.
	IsCore bool
}

// getConfigFormat gets the configured config filename and returns the format
// it is in, if it is a supported format; otherwise an error.
func getConfigFormat() (string, error) {
	// If the format is already set, we don't override the setting.
	// A nil is returned because this is not an error.
	format := os.Getenv(EnvConfigFormat)
	if format != "" {
		// See if the set format is supported, if it is, 
		if isSupportedFormat(format) {
			return format, nil
		}
	}	
	
	// Format is either unknown, or not supported; try to resolve it.
	fname := os.Getenv(EnvConfigFilename)
	if fname == "" {
		return "", errors.New("unable to determine config format, filename not set")
	}
		
	parts := strings.Split(fname, ".")

	// case 0 has already been evaluated
	switch len(parts) {
	case 1: 
		return "", errors.New("unable to determine config format, the configuration file " + fname + " doesn't have an extension")
	case 2:
		format = parts[1]
	default:
		// assume its the last part
		format = parts[len(parts) - 1]
	}

	if isSupportedFormat(format) {
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
                return  false
        }

        return false
}

// SetFromEnv updates the configuration from the environment variable values.
// A setting is only updated if it IsUpdateable.
func SetFromEnv() {
	for k, _ := range AppConfig.Settings {
		// See if k exists as an env variable
		v := os.Getenv(k)
		if v == "" {
			continue
		}

		// Only core and idempotent are updateable
		if !CanSetUsingEnv(k) {
			continue
		}

		// Gotten this far, set it
		AppConfig.Settings[k].Value = v
		AppConfig.Settings[k].SourceIsEnv = true
	}
	
}

// SetFromConfigFile populates configFile from the configured config file.
// The config file entries are then processed, updating their associated
// settings. A setting is only updated if it IsUpdateable.
func SetFromConfigFile() error {
	// ConfigFile should be set and its format type should be known.
	err := LoadConfigFile()
	if err != nil {
		return err
	}

	// ProcessConfigFile, setting what's appropriate.
	for k, v := range configFile {
		// Find the key in the settings
		_, ok := AppConfig.Settings[k]
		if !ok {
			// skip settings that don't already exist
			continue
		}

		// Skip if IsIdempotent, IsCore, SourceIsEnv since they aren't
		// overridable by ConfigFile.
		if !IsUpdateable(k) {
			continue
		}

		// Update the setting with file's
		Update(k, v) 
	}

	return nil
}

// LoadConfigFile() is the entry point for reading the configuration file.
func LoadConfigFile() error {
	n := os.Getenv(EnvConfigFilename)
	if n == "" {
		// This isn't an error as config file is allowed to not exist
		// TODO:
		//	Possible add a ConfigFileRequired flag
		//	Add logging support
		return nil
	}

	fBytes, err := readConfigFile(n)
	if err != nil {
		return err
	}

	err = MarshalFormatReader(os.Getenv(EnvConfigFormat),bytes.NewReader(fBytes)) 
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

// MarshalFormatReader 
func MarshalFormatReader(t string, r io.Reader) error {
	b := new(bytes.Buffer)
	b.ReadFrom(r)

	switch t{
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

// SetIdempotentString sets the value of idempotent configuration settings
// that are strings.
func SetIdempotentString(k, v string) error {
	// see if the key already exists, if it does, it can't be set
	_, ok := AppConfig.Settings[k]
	if ok {
		return nil
	}

	// Set the environment variable for it
	err := os.Setenv(k, v)
	if err != nil {
		return err
	}

	AppConfig.Settings[k] = &setting{Type: "string",  Value: v, IsIdempotent: true}

	return nil
}

// SetIdemString is a convenience function that wraps SetIdempotentString()
func SetIdemString(k, v string) {
	SetIdempotentString(k, v)
}

// SetBoolFlag sets the value of a configuration setting that is also a flag.
func SetBoolFlag(k, f string, v bool) {
	// see if the key exists
	_, ok := AppConfig.Settings[k]
	if ok {
		// see if it can be set
		if AppConfig.Settings[k].IsIdempotent || AppConfig.Settings[k].SourceIsEnv {
			return
		}

		// replace current settings with new
		AppConfig.Settings[k].Value = v
		AppConfig.Settings[k].IsFlag = true
		AppConfig.Settings[k].Code = f
		return
	}
	
	// otherwise add it
	s := utils.BoolToString(v)
	os.Setenv(AppConfig.GetCode() + k, s)
	AppConfig.Settings[k] = &setting{Value: v, Code: f, IsFlag: true}
}

// SetStringFlag sets the value of a configuration setting that is also a flag.
func SetStringFlag(k, f string, v string) {
	// see if the key exists
	_, ok := AppConfig.Settings[k]
	if ok {
		// see if it can be set
		if AppConfig.Settings[k].IsIdempotent || AppConfig.Settings[k].SourceIsEnv || AppConfig.Settings[k].IsCore {
			return
		}

		// replace current settings with new
		AppConfig.Settings[k].Value = v
		AppConfig.Settings[k].IsFlag = true
		AppConfig.Settings[k].Code = f
		return
	}
	
	// otherwise add it
	os.Setenv(AppConfig.GetCode() + k, v)
	AppConfig.Settings[k] = &setting{Value: v, Code: f, IsFlag: true}
}

// Update updates the passed key with the passed value.
func Update(k string, v interface{}) error {
	if !IsUpdateable(k) {
		return nil
	}

	err := os.Setenv(AppConfig.GetCode() + k, v.(string))
	if err != nil {
		return err
	}

	set(k, v)

	return nil
}	


func set(k string, v interface{}) {
	// Cast according to type for this key
	switch AppConfig.Settings[k].Type {
	case "string":
		AppConfig.Settings[k].Value = v.(string)
			
	case "bool":
		AppConfig.Settings[k].Value = v.(bool)

	case "int", "int8", "int16", "int32", "int64":
		AppConfig.Settings[k].Value = v.(int)

	default:
		AppConfig.Settings[k].Value = v
	}
}

// IsUpdateable checks to see if the passed setting key is updateable.
func IsUpdateable(k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be
	// updated.
	_, ok := AppConfig.Settings[k]
	if !ok {
		return false
	}

	// See if there are any settings that prevent it from being updated.
	if AppConfig.Settings[k].IsIdempotent || AppConfig.Settings[k].SourceIsEnv || AppConfig.Settings[k].IsCore {
		return false
	}
	
	return true
}

// Override overrides the setting, if it is overrideable. This is used to
// override any environment variable that had pre-existing values.
func Override(k string, v interface{}) error {
	if !IsOverrideable(k) {
		return nil
	}
	
	err := os.Setenv(AppConfig.GetCode() + k, v.(string))
	if err != nil {
		return err
	}

	// Set the new value
	set(k, v)

	return nil
}

// IsOverrideable() checks to see if the setting can be overridden. Overrides 
// only come from args and flags. ConfigFile settings must be updated instead.
func IsOverrideable(k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be
	// overridden
	_, ok := AppConfig.Settings[k]
	if !ok {
		return false
	}

	// See if there are any settings that prevent it from being overridden.
	if AppConfig.Settings[k].IsIdempotent || AppConfig.Settings[k].IsCore {
		return false
	}
	
	return true	
}

// CanSetUsingEnv checks to see if the setting is settable using an env
// variable.
func CanSetUsingEnv(k string) bool {
	// If something is flagged as idempotent, it can't be set be changed.
	// Same with core settings.
	if AppConfig.Settings[k].IsIdempotent || AppConfig.Settings[k].IsCore {
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
		err := errors.New(alias + " is an alias of the command " + v + " cannot make it an alias of " + command)
		return err
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
	AppConfig = &Config{Settings: map[string]*setting{}}
}

func Get(k string) (interface{}, error) {
	_, ok := AppConfig.Settings[k]
	if !ok {
		return nil, notFoundErr(k)
	}

	return AppConfig.Settings[k].Value, nil
}

func GetBool(k string) (bool, error) {
	_, ok := AppConfig.Settings[k]
	if !ok {
		return false, notFoundErr(k)
	}
	
	return AppConfig.Settings[k].Value.(bool), nil
}

func GetInt(k string) (int, error) {
	_, ok := AppConfig.Settings[k]
	if !ok {
		return 0, notFoundErr(k)
	}

	return AppConfig.Settings[k].Value.(int), nil
}

func GetString(k string) (string, error) {
	_, ok := AppConfig.Settings[k]
	if !ok {
		return "", notFoundErr(k)
	}

	return AppConfig.Settings[k].Value.(string), nil
}

// GetInterface is a convenience wrapper function to Get
func GetInterface(k string) (interface{}, error) {
	return Get(k)
}

func notFoundErr(k string) error {
	return errors.New(k + " not found")
}
