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

// configFile holds the contents of the configuration file
var configFile map[string]interface{} = make(map[string]interface{})

// config holds the current application configuration
var AppConfig *Config = &Config{Settings: map[string]*setting{}}

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

// NewConfig returns a *Config to the caller
func NewConfig() *Config {
	AppConfig = &Config{Settings: map[string]*setting{}}
	return AppConfig
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
func SetBoolFlag(k, v string, b bool) {
	// see if the key exists
	_, ok := AppConfig.Settings[k]
	if ok {
		// see if it can be overridden
		if AppConfig.Settings[k].IsIdempotent || AppConfig.Settings[k].SourceIsEnv {
			return
		}

		// override it
		AppConfig.Settings[k].Value = b
		AppConfig.Settings[k].IsFlag = true
		AppConfig.Settings[k].Code = v
		return
	}
	
	// otherwise add it
	bs := utils.BoolToString(b)
	os.Setenv(AppConfig.GetCode() + k, bs)
	AppConfig.Settings[k] = &setting{Value: b, Code: v, IsFlag: true}
}

// resetAppConfig resets the application's configuration struct to empty.
// This does not affect their respective environment variables
func resetAppConfig() {
	AppConfig = &Config{Settings: map[string]*setting{}}
}

			
/*
func Get(k string) interface{} {

}

func GetBool(k string) bool {

}

func GetInt(k string) int {

}

// GetInterface is a convenience wrapper function to Get
func GetInterface(k string) interface{} {
	return Get(k)
}

*/
