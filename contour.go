//
package contour

import (
	"bytes"
	"encoding/json"
	// "encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	// "os"
	// "strconv"
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

func (f Format) String() string {
	switch f {
	case JSON:
		return "json"
	case TOML:
		return "toml"
	case YAML:
		return "yaml"
	case XML:
		return "xml"
	default:
		return "unsupported"
	}
}

// isSupported checks to see if the passed string represents a supported config
// format.
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

// Contour variable names for the pre-configured core setting names that it
// comes with. These are public and are directly settable if you wish to use
// your own values, just set them before doing anything with Contour.
var (
	CfgFile   string = "cfgfile"
	CfgFormat string = "cfgformat"
)

func init() {
	initCfgs()
}

// initConfigs initializes the configs var. This can be called to reset it in
// testing too.
func initCfgs() {
	appCfg = NewCfg(app)
}

// formatFromFilename gets the format from the passed filename.  An error will
// be returned if either the format isn't supported or the extension doesn't
// exist.  If the passed string has multiple dots, the last dot is assumed to
// be the extension.
func formatFromFilename(s string) (Format, error) {
	if s == "" {
		return Unsupported, fmt.Errorf("no config filename")
	}
	parts := strings.Split(s, ".")
	format := ""
	// case 0 has already been evaluated
	switch len(parts) {
	case 1:
		return Unsupported, fmt.Errorf("unable to determine %s's config format: no extension", strings.TrimSpace(s))
	case 2:
		format = parts[1]
	default:
		// assume its the last part
		format = parts[len(parts)-1]
	}
	f := ParseFormat(format)
	if !f.isSupported() {
		return Unsupported, unsupportedFormatErr(format)
	}
	return f, nil
}

// getCfgFile() is the entry point for reading the configuration file.
func (c *Cfg) getFile() (cfg interface{}, err error) {
	// if it's not set to use a cfg file, nothing to do
	if !c.useCfgFile {
		return nil, nil
	}
	setting, ok := c.settings[CfgFile]
	if !ok {
		// Wasn't configured, nothing to do. Not an error.
		return nil, nil
	}
	n := setting.Value.(string)
	if n == "" {
		// This isn't an error as config file is allowed to not exist
		// TODO:
		//	Possible add a CfgFileRequired flag
		return nil, nil
	}
	// This shouldn't happend, but lots of things happen that shouldn't.  It should
	// have been registered already. so if it doesn't exit, err.
	format, ok := c.settings[CfgFormat]
	if !ok {
		return nil, fmt.Errorf("cfgvformat was not set")
	}
	if format.Value.(string) == "" {
		return nil, fmt.Errorf("cfg format was not set")
	}
	fBytes, err := readCfgFile(n)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", n, err)
	}
	format, _ = c.settings[CfgFormat]
	cfg, err = unmarshalFormatReader(ParseFormat(format.Value.(string)), bytes.NewReader(fBytes))
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling %s: %s", n, err)
	}
	return cfg, nil
}

// readCfgFile reads the configFile and returns the resulting slice. The entire
// contents of the file are read at once.
func readCfgFile(n string) ([]byte, error) {
	cfg, err := ioutil.ReadFile(n)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// unmarshalFormatReader accepts an io.Reader and unmarshals it using the
// correct format.
//
// Supported formats:
//   json
//   toml
// TODO
//   add YAML support
//   add HCL support
func unmarshalFormatReader(f Format, r io.Reader) (interface{}, error) {
	b := new(bytes.Buffer)
	b.ReadFrom(r)
	var ret interface{}
	switch f {
	case JSON:
		err := json.Unmarshal(b.Bytes(), &ret)
		if err != nil {
			return nil, err
		}
	case TOML:
		_, err := toml.Decode(b.String(), &ret)
		if err != nil {
			return nil, err
		}
	default:
		err := unsupportedFormatErr(f.String())
		return nil, err
	}
	return ret, nil
}

// canUpdate checks to see if the passed setting key is updateable. If the key
// is not updateable, a false is returned along with an error.
func canUpdate(k string) (bool, error) { return appCfg.canUpdate(k) }
func (c *Cfg) canUpdate(k string) (bool, error) {
	// See if the key exists, if it doesn't already exist, it can't be updated.
	s, ok := c.settings[k]
	if !ok {
		return false, fmt.Errorf("cannot update %q: not found", k)
	}
	// See if there are any settings that prevent it from being overridden.  Core and
	// environment variables are never settable. Core must be set during registration.
	if s.IsCore {
		return false, fmt.Errorf("cannot update %q: core settings cannot be updated", k)
	}
	// Everything else is updateable.
	return true, nil
}

// canOverride is a helper funciton for the global appCfg.
func canOverride(k string) bool {
	return appCfg.canOverride(k)
}

// canOverride() checks to see if the setting can be overridden. Overrides only
// come from flags. If it can't be overridden, it must be set via application,
// environment variable, or cfg file.
func (c *Cfg) canOverride(k string) bool {
	// See if the key exists, if it doesn't already exist, it can't be overridden
	_, ok := c.settings[k]
	if !ok {
		return false
	}
	// See if there are any settings that prevent it from being overridden.
	// Core can never be overridden-must be a flag to override.
	if c.settings[k].IsCore || !c.settings[k].IsFlag {
		return false
	}
	return true
}

// notFoundErr returns a standardized not found error.
func notFoundErr(k string) error {
	return fmt.Errorf("%s not found", k)
}

// settingNotFoundErr standardizes the settingNotFound err output.
func settingNotFoundErr(k string) error {
	return fmt.Errorf("setting not found: %s", k)
}

// unsupportedFormatErr standardizes the unsupportedFormatErr output.
func unsupportedFormatErr(k string) error {
	return fmt.Errorf("unsupported cfg format: %s", k)
}
