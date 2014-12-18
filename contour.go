//
package contour

import (
	"bytes"
	"encoding/json"
	_ "encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	_ "os"
	_ "strconv"
	"strings"

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

// appCfg: contour's global config; contour config functinos operate on this.
var appCfg *Cfg

// Contour ironment variable names for the pre-configured core setting names
// that it comes with. These are public and are directly settable if you wish
// to use your own values. Just set them before doing anything with Contour.
var (
	CfgFilename string = "cfgfilename"
	CfgFormat   string = "cfgformat"
)

func init() {
	initCfgs()
}

// initConfigs initializes the configs var. This can be called to reset it in
// testing too.
func initCfgs() {
	appCfg = &Cfg{name: app, settings: map[string]*setting{}}
}

// formatFromFilename gets the format from the passed filename.
// If the returned format is not supported, or there isn't a file
// extension, an error is returned. If the passed string has
// multiple dots, the last dot is assumed to be the extension.
func formatFromFilename(s string) (Format, error) {
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
	if !f.isSupported() {
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
		return false
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
// getCfgFile() is the entry point for reading the configuration file.
func (c *Cfg) getFile() (cfg interface{}, err error) {
	setting, ok := c.settings[CfgFilename]
	if !ok {
		// Wasn't configured, nothing to do. Not an error.
		return nil, nil
	}

	n := setting.Value.(string)
	if n == "" {
		// This isn't an error as config file is allowed to not exist
		// TODO:
		//	Possible add a ConfigFileRequired flag
		//	Add logging support
		return nil, nil
	}

	fmt.Printf("%#v\n", c.settings)
	// This shouldn't happend, but lots of things happen that shouldn't.
	// It should have been registered already. so if it doesn't exit, err.
	format, ok := c.settings[CfgFormat]
	if !ok {
		return nil, fmt.Errorf("Unable to load the cfg file, the configuration format type was not set")
	}
	if format.Value.(string) == "" {
		return nil, fmt.Errorf("Unable to load the cfg file, the configuration format type was not set")
	}

	fBytes, err := readCfg(n)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	format, _ = c.settings[CfgFormat]
	cfg, err = marshalFormatReader(ParseFormat(format.Value.(string)), bytes.NewReader(fBytes))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return cfg, nil
}

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
func marshalFormatReader(f Format, r io.Reader) (interface{}, error) {
	b := new(bytes.Buffer)
	b.ReadFrom(r)

	var ret interface{}
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
		err := unsupportedFormatErr(f.String())
		logger.Error(err)
		return nil, err
	}
	return ret, nil
}

// canUpdate checks to see if the passed setting key is updateable.
// TODO the logic flow is wonky because it could be simplified but
// want hard check for core and not sure about conf/flag/env stuff yet.
// so the wierdness sits for now.
func (c *Cfg) canUpdate(k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be
	// updated.
	_, ok := c.settings[k]
	if !ok {
		return false
	}

	// See if there are any settings that prevent it from being overridden.
	// Core and ironment variables are never settable. Core must be set
	// during registration.
	if c.settings[k].IsCore {
		return false
	}

	// Only flags and conf types are updateable, otherwise they must be
	// registered or set.
	if c.settings[k].IsCfg || c.settings[k].IsFlag {
		return true
	}

	return false
}

func canUpdate(k string) bool {
	return appCfg.canUpdate(k)
}

// canOverride() checks to see if the setting can be overridden. Overrides
// only come from flags. If it can't be overridden, it must be set via
// application, environment variable, or cfg file.
func (c *Cfg) canOverride(k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be
	// overridden
	_, ok := c.settings[k]
	if !ok {
		return false
	}

	// See if there are any settings that prevent it from being overridden.
	// Core can never be overridden
	// Must be a flag to override.
	if c.settings[k].IsCore || !c.settings[k].IsFlag {
		return false
	}

	return true
}

func canOverride(k string) bool {
	return appCfg.canOverride(k)
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
	return fmt.Errorf("unsupported configuration file format: %s", k)
}
