//
package contour

import (
	_"bytes"
	_"encoding/json"
	_"encoding/xml"
	"errors"
	_"flag"
	_"fmt"
	_"io"
	_"io/ioutil"
	_"os"
	"strings"
	_"strconv"

	_ "code.google.com/p/gcfg"
	_ "github.com/BurntSushi/toml"
)

const app = "app"

// configs allows for support of multiple configurations. The main application
// config is 'app'. Calling any of Contour's function versions of
// config.method() is the equivelant of calling config[app].method().
var configs map[string]*Cfg

// Contour Environment variable names for the pre-configured core setting names
// that it comes with. These are public and are directly settable if you wish
// to use your own values. Just set them before doing anything with Contour.
var (
	EnvCfgFile	string = "cfgfile"
	EnvCfgFormat	string = "cfgformat"
	EnvLogCfgFile	string = "logcfgfile"
	EnvLogging	string = "logging"
)

func init() {
	initConfigs()
}
// configFormat gets the configured config filename and returns the format
// it is in, if it is a supported format; otherwise an error.
func configFormat(s string) (string, error) {
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
	case "yaml", "yml":
		return true
	case "ini":
		return true
	case "xml":
		return true
	default:
		return false
	}

	return false
}

/*
// loadEnvs updates the configuration from the environment variable values.
// A setting is only updated if it IsUpdateable.
func loadEnvs() {
	if !appConfig.useEnv {
		return
	}

	for k, _ := range appConfig.settings {
		// See if k exists as an env variable
		v := os.Getenv(k)
		if v == "" {
			continue
		}

		// Core is not updateable
		if appConfig.settings[k].IsCore {
			continue
		}

		// If its readonly, see if its set. If it isn't it can be.
		if appConfig.settings[k].Immutable {
			if appConfig.settings[k].Value != nil {
				continue
			}
		}

		// Gotten this far, set it
		appConfig.settings[k].Value = v
		appConfig.settings[k].IsEnv = true
	}

}

// loadConfigFile() is the entry point for reading the configuration file.
func loadConfigFile() error {
	setting, ok := appConfig.settings[EnvConfigFilename]
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
	if appConfig.settings[EnvConfigFormat].Value == nil {
		return errors.New("Unable to load configuration value, its format type was not set")
	}

	fBytes, err := readConfigFile(n)
	if err != nil {
		return err
	}

	err = marshalFormatReader(appConfig.settings[EnvConfigFormat].Value.(string), bytes.NewReader(fBytes))
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
	_, ok := appConfig.settings[k]
	if !ok {
		return false
	}

	// See if there are any settings that prevent it from being overridden.
	// Core and Environment variables are never settable. Core must be set
	// during registration and Environment variables must be set using
	// the Override functions.
	if appConfig.settings[k].IsCore || appConfig.settings[k].IsEnv {
		return false
	}

	// Immutable variables are only settable if they are not set.
	// This does not apply to boolean as there is no way to determine if
	// the value is unset. So bool immutables are only writable when they
	// are registered.
	if (appConfig.settings[k].Immutable && appConfig.settings[k].Value != "") || (appConfig.settings[k].Immutable && appConfig.settings[k].Type == "bool") {
		return false
	}
	return true
}

// canOverride() checks to see if the setting can be overridden. Overrides
// only come from args and flags. ConfigFile settings must be set instead.
func canOverride(k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be
	// overridden
	_, ok := appConfig.settings[k]
	if !ok {
		return false
	}

	// See if there are any settings that prevent it from being overridden.
	if appConfig.settings[k].IsCore {
		return false
	}

	// Immutable variables are only settable if they are not set.
	// This does not apply to boolean as there is no way to determine if
	// the value is unset. So bool immutables are only writable when they
	// are registered.
	if (appConfig.settings[k].Immutable && appConfig.settings[k].Value != "") || (appConfig.settings[k].Immutable && appConfig.settings[k].Type == "bool") {
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


*/

// initConfigs initializes the configs var. This can be called to reset it in
// testing too.
func initConfigs() {
	configs = make(map[string]*Cfg)
}

// notFoundErr returns a standardized not found error.
func notFoundErr(k string) error {
	return errors.New(k + " not found")
}

// settingNotFoundErr adds the suffix ": setting " to k before calling
// notFoundErr
func settingNotFoundErr(k string) error {
	return notFoundErr(k + ": setting")
}
