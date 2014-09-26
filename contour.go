//
package contour

import (
	"bytes"
	"encoding/json"
	_ "encoding/xml"
	_ "flag"
	"fmt"
	"io"
	"io/ioutil"
	_ "os"
	_ "strconv"
	"strings"

	_ "code.google.com/p/gcfg"
	"github.com/BurntSushi/toml"
)

const app = "app"

const (
	Unsupported Format = iota
	JSON
	TOML
	YAML
	XML
)

type Format int

func ParseFormatE(s string) (Format, error) {
	ls := strings.ToLower(s)
	switch ls {
	case "json":
		return JSON, nil
	case "toml":
		return TOML, nil
	case "yaml", "yml":
		return YAML, nil
	case "xml":
		return XML, nil
	}

	return Unsupported, unsupportedFormatErr(ls)
}

func ParseFormat(s string) Format {
	f, _ := ParseFormatE(s)
	return f
}

// configs allows for support of multiple configurations. The main application
// config is 'app'. Calling any of Contour's function versions of
// config.method() is the equivelant of calling config[app].method().
var configs []*Cfg
var configNames []string
var configCount int

// Contour ironment variable names for the pre-configured core setting names
// that it comes with. These are public and are directly settable if you wish
// to use your own values. Just set them before doing anything with Contour.
var (
	CfgFile   string = "cfgfile"
	CfgFormat string = "cfgformat"
)

// The following settings are used to gain better control over the configs
// slice growth, to minimize garbage generation. Practically speaking, only
// growthConfigIncrement makes any difference as the initConfigCap is used
// during init() and any changes to it will not affect run-time behavior.
// Instead, fork contour and customize that value for your application if you
// expect your application to benefit from a different initConfigCap value.
//
// initConfigCap sets the initial capacity for configurations
// growthConfigIncrement sets the amount to grow the config slice.
var initConfigCap int = 4
var growConfigIncrement int

// The following settings are used to gain better control over the settings
// slice growth, to minimize garbage generation.
// initSettingsCap sets the initial capacity for settings.
// growSettingsIncrement sets the amount to grow the settings slice.
// These can be specified per Cfg too.
var initSettingsCap int = 10
var growSettingsIncrement int = 10

func init() {
	initConfigs()
}

// SetGrowthConfigIncrement sets the amount by which to grow the capacity of
// the configs slice; 0 (zero) means double capacity each time, which is the
// default.
func SetGrowConfigIncrement(i int) {
	if i < 0 {
		i = 0
	}

	growConfigIncrement = i
}

// SetInitSettingsCap sets the initial capacity for the settings slice in new
// configurations.
func SetInitSettingsCap(i int) {
	if i < 2 {
		i = 2
	}
	initSettingsCap = i
}

// SetGrowSettingsIncrement sets the amount by which to grow the capacity of
// the settings slice; 0 (zero) means double capacity each time.
func SetGrowSettingsIncrement(i int) {
	if i < 0 {
		i = 0
	}

	growSettingsIncrement = i
}

// configFormat gets the configured config filename and returns the format
// it is in, if it is a supported format; otherwise an error.
func configFormat(s string) (Format, error) {
	if s == "" {
		return Unsupported, fmt.Errorf("a config filename was expected, none received")
	}

	parts := strings.Split(s, ".")
	format := ""

	// case 0 has already been evaluated
	switch len(parts) {
	case 1:
		return Unsupported, fmt.Errorf("unable to determine config format, the configuration file, " + strings.TrimSpace(s) + ", doesn't have an extension")

	case 2:
		format = parts[1]

	default:
		// assume its the last part
		format = parts[len(parts)-1]
	}

	f := ParseFormat(format)
	is := f.isSupported()
	if !is {
		err := unsupportedFormatErr(format)
		logger.Error(err)
		return Unsupported, err
	}

	return f, nil

}

// isSupportedFormat checks to see if the passed string represents a supported
// config format.
func (f Format) isSupported() bool {
	switch f {
	case YAML:
		return true
	case JSON:
		return true
	case TOML:
		return true
	case XML:
		return true
	}

	return false
}

func (f Format) String() string {
	switch f {
	case YAML:
		return "yaml"
	case JSON:
		return "json"
	case TOML:
		return "toml"
	case XML:
		return "xml"
	}

	return "unsupported"
}
// RegisterConfigFilename set's the configuration file's name. The name is
// parsed for a valid extension--one that is a supported format--and saves
// that value too. If it cannot be determined, the extension info is not set.
// These are considered core values and cannot be changed from command-line
// and configuration files. (IsCore == true).
func RegisterConfigFilename(k, v string) error {
	if v == "" {
		return fmt.Errorf("A config filename was expected, none received")
	}

	if k == "" {
		return fmt.Errorf("A key for the config filename setting was expected, none received")
	}

	configs[0].RegisterStringCore(k, v)

	// TODO redo this given new paradigm
	// Register it first. If a valid config format isn't found, an error
	// will be returned, so registering it afterwords would mean the
	// setting would not exist.
	configs[0].RegisterString(CfgFormat, "")
	format, err := configFormat(v)
	if err != nil {
		return err
	}

	configs[0].RegisterString(CfgFormat, format.String())

	return nil
}

// RegisterSetting checks to see if the entry already exists and adds the
// new setting if it does not.
func RegisterSetting(Type string, k string, short string, v interface{}, Usage string, Default string, IsCore, IsCfg, IsFlag bool) {
	configs[0].RegisterSetting(Type, k, short, v, Usage, Default, IsCore, IsCfg, IsFlag)
}

/*
// loads updates the configuration from the environment variable values.
// A setting is only updated if it IsUpdateable.
func loads() {
	if !appConfig.use {
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
		appConfig.settings[k].Is = true
	}

}
*/

// readCfg reads the configFile
func readCfg(n string) ([]byte, error) {
	cfg, err := ioutil.ReadFile(n)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return cfg, nil
}

// marshalFormatReader
func marshalFormatReader(f Format, r io.Reader) (map[string]interface{}, error) {
	b := new(bytes.Buffer)
	b.ReadFrom(r)

	ret := make(map[string]interface{})
	switch f {
	case JSON:
		err := json.Unmarshal(b.Bytes(), &ret)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

	case TOML:
		_, err := toml.Decode(b.String(), &ret)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	default:
		err := fmt.Errorf("unsupported configuration file format type: %s", )
		logger.Error(err)
		return nil, err
	}
	return ret, nil
}

func canUpdate(i int) bool {
	return configs[0].canUpdate(i)
}

func canOverride(i int) bool {
	return configs[0].canOverride(i)
}

/*
// AddCommandAlias adds an alias for a command. The first time a command is
// added, it's added as an alias of itself too.
func AddCommandAlias(command, alias string) error {
	// see if an alias already exists
	v, ok := commandAlias[alias]
	if ok {
		return fmt.Errorf(alias + " is an alias of the command " + v + " cannot make it an alias of " + command)
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
		err := fmt.Errorf(alias + " is an alias of the setting " + v + " cannot make it an alias of " + setting)
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
	configs = make([]*Cfg, 1)
	configs[0] = &Cfg{name: app, settings: make([]*setting, initSettingsCap)}
	configNames = make([]string, 1)
	configNames[0] = app
}

// Convenience functions for the main config
// Code returns the code for the config. If set, this is used as
// the prefix for environment variables and configuration setting names.
func Code() string {
	return configs[0].code
}

func UseEnv() bool {
	return configs[0].useEnv
}

// SetConfig goes through the initialized Settings and updates the updateable
// Settings if a new, valid value is found. This applies to, in order: Env
// variables and config files. For any that are not found, or that are
// immutable, once set, the original initialization values are used.
//
// The merged configuration Settings are then  written to their respective
// environment variables. At this point, only args, or in application setting
// changes, can change the non-immutable Settings.
func SetConfig() error {
	return configs[0].SetConfig()
}

// SetCode set's the code for this configuration. This can only be done once.
// If it is already set, it will return an error.
func SetCode(s string) error {
	if configs[0].code != "" {
		return fmt.Errorf("appCode is already set. AppCode is immutable. Once set, it cannot be altered")
	}

	configs[0].code = s
	return nil
}

// Config processed returns whether or not all of the config's settings have
// been processed.
func ConfigProcessed() bool {
	return configs[0].ConfigProcessed()
}

func configIndex(k string) (int, error)  {
	for i, v := range configNames {
		if v == k {
			return i, nil
		}
	}
	
	return -1, fmt.Errorf("%q config not found", k)
}

// notFoundErr returns a standardized not found error.
func notFoundErr(k string) error {
	return fmt.Errorf("%s not found", k)
}

// settingNotFoundErr adds the suffix ": setting " to k before calling
// notFoundErr
func settingNotFoundErr(k string) error {
	return notFoundErr(fmt.Sprintf("%s: setting", k))
}

func unsupportedFormatErr(k string) error {
	return fmt.Errorf("%s is not a supported configuration file format", k)
}
